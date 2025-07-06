package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medivue/adapters/db"
	"github.com/medivue/core/domain"
	"github.com/medivue/core/utils"
)

// PaymentError represents a payment service error with context
type PaymentError struct {
	Code           string         // Internal error code
	Message        string         // User-friendly error message
	Err            error          // Original error
	TransactionID  string         // Payment transaction ID if available
	PaymentID      string         // Payment ID if available
	ProviderError  error          // Original payment provider error
	HTTPStatusCode int            // HTTP status code for the error
	Retryable      bool           // Whether the error can be retried
	Metadata       map[string]any // Additional error context
}

func (e *PaymentError) Error() string {
	if e.TransactionID != "" {
		return fmt.Sprintf("%s: %s (Transaction: %s) - %v", e.Code, e.Message, e.TransactionID, e.Err)
	}
	return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
}

func (e *PaymentError) Unwrap() error {
	return e.Err
}

// Error codes for payment service
const (
	ErrCodeValidation           = "VALIDATION_ERROR"
	ErrCodePaymentVerification  = "VERIFICATION_ERROR"
	ErrCodePaymentNotFound      = "PAYMENT_NOT_FOUND"
	ErrCodeInvalidStatus        = "INVALID_STATUS"
	ErrCodeSystemError          = "SYSTEM_ERROR"
	ErrCodeProviderError        = "PROVIDER_ERROR"
	ErrCodeInsufficientFunds    = "INSUFFICIENT_FUNDS"
	ErrCodeDuplicateTransaction = "DUPLICATE_TRANSACTION"
	ErrCodeTimeout              = "TIMEOUT_ERROR"
	ErrCodeNetworkError         = "NETWORK_ERROR"
	ErrCodeInvalidAmount        = "INVALID_AMOUNT"
	ErrCodeAuthenticationError  = "AUTH_ERROR"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	RetryableErrors []string
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     2 * time.Second,
		Multiplier:      2.0,
		RetryableErrors: []string{
			ErrCodeTimeout,
			ErrCodeNetworkError,
			ErrCodePaymentVerification,
		},
	}
}

// isRetryableError checks if an error should trigger a retry
func (s *ServicesHandler) isRetryableError(err error) bool {
	var pErr *PaymentError
	if errors.As(err, &pErr) {
		if pErr.Retryable {
			return true
		}
		for _, code := range DefaultRetryConfig().RetryableErrors {
			if pErr.Code == code {
				return true
			}
		}
		return false
	}

	// For non-PaymentError types, retry on certain error patterns
	errStr := err.Error()
	return strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporary") ||
		strings.Contains(errStr, "deadline exceeded")
}

// withRetry executes the given function with exponential backoff retry logic
func (s *ServicesHandler) withRetry(ctx context.Context, operation string, f func() error) error {
	config := DefaultRetryConfig()
	var lastErr error
	interval := config.InitialInterval

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err := f()
		if err == nil {
			return nil
		}

		lastErr = err
		if attempt == config.MaxAttempts || !s.isRetryableError(err) {
			break
		}

		// Log retry attempt with structured logging
		utils.Warn("Operation failed, retrying...",
			utils.LogField{Key: "operation", Value: operation},
			utils.LogField{Key: "attempt", Value: attempt},
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "next_retry_in", Value: interval.String()})

		// Wait with exponential backoff
		select {
		case <-ctx.Done():
			return &PaymentError{
				Code:    ErrCodeTimeout,
				Message: "operation cancelled",
				Err:     ctx.Err(),
				Metadata: map[string]interface{}{
					"operation": operation,
					"attempts":  attempt,
				},
				Retryable: false,
			}
		case <-time.After(interval):
			interval = time.Duration(float64(interval) * config.Multiplier)
			if interval > config.MaxInterval {
				interval = config.MaxInterval
			}
		}
	}

	// If we got here, all retries failed
	var pErr *PaymentError
	if errors.As(lastErr, &pErr) {
		pErr.Retryable = false // Mark as non-retryable since we've exhausted retries
		return pErr
	}

	return &PaymentError{
		Code:    ErrCodeSystemError,
		Message: fmt.Sprintf("operation %s failed after %d attempts", operation, config.MaxAttempts),
		Err:     lastErr,
		Metadata: map[string]interface{}{
			"operation": operation,
			"attempts":  config.MaxAttempts,
		},
		Retryable: false,
	}
}

// validatePaymentStatus checks if a status transition is valid
func (s *ServicesHandler) validatePaymentStatus(current, new db.PaymentStatus) error {
	// Define valid status transitions
	validTransitions := map[db.PaymentStatus][]db.PaymentStatus{
		db.PaymentStatusPending:   {db.PaymentStatusSuccess, db.PaymentStatusFailed},
		db.PaymentStatusSuccess:   {db.PaymentStatusRefunded},
		db.PaymentStatusFailed:    {db.PaymentStatusPending}, // Allow retry
		db.PaymentStatusRefunded:  {},                        // Terminal state
		db.PaymentStatusCancelled: {db.PaymentStatusPending}, // Allow retry after cancellation
	}

	if transitions, exists := validTransitions[current]; exists {
		for _, valid := range transitions {
			if valid == new {
				return nil
			}
		}
	}

	return &PaymentError{
		Code:    ErrCodeInvalidStatus,
		Message: fmt.Sprintf("invalid status transition from %s to %s", current, new),
		Err:     fmt.Errorf("invalid status transition from %s to %s", current, new),
		Metadata: map[string]interface{}{
			"current_status": current,
			"new_status":     new,
		},
	}
}

// validateAmount validates payment amount constraints
func (s *ServicesHandler) ValidateAmount(amount float64, currency string) error {
	// Define currency-specific constraints
	minAmounts := map[string]float64{
		"NGN": 100.00, // Minimum amount for NGN
		"USD": 0.50,   // Minimum amount for USD
		"EUR": 0.50,   // Minimum amount for EUR
	}

	maxAmounts := map[string]float64{
		"NGN": 10000000.00, // Maximum amount for NGN
		"USD": 100000.00,   // Maximum amount for USD
		"EUR": 100000.00,   // Maximum amount for EUR
	}

	// Check if currency is supported
	if _, ok := minAmounts[currency]; !ok {
		return &PaymentError{
			Code:    ErrCodeValidation,
			Message: fmt.Sprintf("unsupported currency: %s", currency),
			Err:     fmt.Errorf("unsupported currency"),
			Metadata: map[string]interface{}{
				"currency": currency,
			},
		}
	}

	// Validate amount range
	if amount < minAmounts[currency] {
		return &PaymentError{
			Code:    ErrCodeInvalidAmount,
			Message: fmt.Sprintf("amount below minimum allowed for %s (%.2f)", currency, minAmounts[currency]),
			Err:     fmt.Errorf("amount below minimum"),
			Metadata: map[string]interface{}{
				"currency":    currency,
				"amount":      amount,
				"min_allowed": minAmounts[currency],
			},
		}
	}

	if amount > maxAmounts[currency] {
		return &PaymentError{
			Code:    ErrCodeInvalidAmount,
			Message: fmt.Sprintf("amount above maximum allowed for %s (%.2f)", currency, maxAmounts[currency]),
			Err:     fmt.Errorf("amount above maximum"),
			Metadata: map[string]interface{}{
				"currency":    currency,
				"amount":      amount,
				"max_allowed": maxAmounts[currency],
			},
		}
	}

	return nil
}

// validateAndParseMetadata validates payment metadata and ensures it has required fields
func (s *ServicesHandler) validateAndParseMetadata(rawMetadata interface{}) ([]byte, error) {
	if rawMetadata == nil {
		return nil, nil
	}

	// Convert to JSON
	metadata, err := json.Marshal(rawMetadata)
	if err != nil {
		return nil, &PaymentError{
			Code:    ErrCodeValidation,
			Message: "failed to marshal payment metadata",
			Err:     err,
			Metadata: map[string]interface{}{
				"raw_metadata": fmt.Sprintf("%+v", rawMetadata),
			},
		}
	}

	// Validate structure and required fields
	var parsed map[string]interface{}
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, &PaymentError{
			Code:    ErrCodeValidation,
			Message: "invalid metadata format",
			Err:     err,
			Metadata: map[string]interface{}{
				"metadata": string(metadata),
			},
		}
	}

	// Check for required fields
	requiredFields := []string{
		"provider_reference",
		"transaction_time",
		"payment_method",
		"currency",
	}

	missingFields := []string{}
	for _, field := range requiredFields {
		if _, exists := parsed[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return nil, &PaymentError{
			Code:    ErrCodeValidation,
			Message: "missing required metadata fields",
			Err:     fmt.Errorf("missing fields: %v", missingFields),
			Metadata: map[string]interface{}{
				"missing_fields":  missingFields,
				"provided_fields": parsed,
			},
		}
	}

	// Add validation timestamp and version
	parsed["validation_timestamp"] = time.Now().UTC().Format(time.RFC3339)
	parsed["metadata_version"] = "1.0"

	// Re-marshal with validation info
	metadata, err = json.Marshal(parsed)
	if err != nil {
		return nil, &PaymentError{
			Code:    ErrCodeSystemError,
			Message: "failed to marshal validated metadata",
			Err:     err,
			Metadata: map[string]interface{}{
				"parsed_data": parsed,
			},
		}
	}

	return metadata, nil
}

// isPaymentVerificationError determines if an error is a verification error that should trigger a retry
func (s *ServicesHandler) IsPaymentVerificationError(err error) bool {
	var pErr *PaymentError
	if errors.As(err, &pErr) {
		return pErr.Code == ErrCodePaymentVerification
	}
	// For non-PaymentError types, retry on certain error patterns
	return err != nil && (strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "temporary"))
}

// CreatePayment creates a new payment for an appointment
func (s *ServicesHandler) CreatePayment(ctx echo.Context) error {
	dto := ctx.Get(utils.ValidatedBodyDTO).(*domain.CreatePaymentDTO)
	// Validate appointment exists and is in valid state
	appointment, err := s.AppointmentRepo.GetAppointment(ctx.Request().Context(), dto.AppointmentID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return fmt.Errorf("appointment not found: %v", utils.ErrNotFound)
		}
		return fmt.Errorf("failed to get appointment: %v", err)
	}

	if appointment.Status != db.AppointmentStatusPending && appointment.Status != db.AppointmentStatusConfirmed {
		return fmt.Errorf("appointment is not in valid state for payment: %v", utils.ErrBadRequest)
	}

	// Create payment
	var pgAmount pgtype.Numeric
	if err := pgAmount.Scan(dto.Amount); err != nil {
		return fmt.Errorf("invalid amount format: %v", utils.ErrBadRequest)
	}

	var pgMetadata []byte
	if dto.PaymentMetadata != nil {
		metadata, err := json.Marshal(dto.PaymentMetadata)
		if err != nil {
			return fmt.Errorf("invalid payment metadata format: %v", utils.ErrBadRequest)
		}
		pgMetadata = metadata
	}

	appointmentUUID, err := uuid.Parse(dto.AppointmentID)
	if err != nil {
		return fmt.Errorf("invalid appointment ID: %v", utils.ErrBadRequest)
	}

	payment, err := s.PaymentRepo.CreatePayment(ctx.Request().Context(), db.Create_PaymentParams{
		AppointmentID:      appointmentUUID.String(),
		PatientID:          appointment.PatientID,
		DiagnosticCentreID: appointment.DiagnosticCentreID,
		Amount:             pgAmount,
		Currency:           dto.Currency,
		PaymentMethod:      db.PaymentMethod(dto.PaymentMethod),
		PaymentMetadata:    pgMetadata,
	})
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return utils.ResponseMessage(http.StatusCreated, payment, ctx)
}

// GetPayment retrieves a payment by ID
func (s *ServicesHandler) GetPayment(ctx echo.Context) error {
	dto := ctx.Get(utils.ValidatedQueryParamDTO).(*domain.GetPaymentDTO)
	payment, err := s.PaymentRepo.GetPayment(ctx.Request().Context(), dto.PaymentID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return fmt.Errorf("payment not found: %v", utils.ErrNotFound)
		}
		return fmt.Errorf("failed to get payment: %v", err)
	}
	return utils.ResponseMessage(http.StatusOK, payment, ctx)
}

// ListPayments lists payments with filtering
func (s *ServicesHandler) ListPayments(ctx echo.Context) error {

	dto := ctx.Get(utils.ValidatedQueryParamDTO).(*domain.ListPaymentsDTO)
	var params db.List_PaymentsParams
	offset := (dto.Page - 1) * dto.PageSize
	params.Limit = int32(dto.PageSize)
	params.Offset = int32(offset)

	// Handle diagnostic centre ID
	if dto.DiagnosticCentreID != "" {
		centre, err := uuid.Parse(dto.DiagnosticCentreID)
		if err != nil {
			return fmt.Errorf("invalid diagnostic centre ID: %v", utils.ErrBadRequest)
		}
		params.Column1 = centre.String()
	}

	// Handle patient ID
	if dto.PatientID != "" {
		patient, err := uuid.Parse(dto.PatientID)
		if err != nil {
			return fmt.Errorf("invalid patient ID: %v", utils.ErrBadRequest)
		}
		params.Column2 = patient.String()
	}

	// Handle status
	if dto.Status != "" {
		params.Column3 = dto.Status
	}

	// Handle date range
	if !dto.FromDate.IsZero() {
		params.Column4.Time = dto.FromDate
		params.Column4.Valid = true
	}

	if !dto.ToDate.IsZero() {
		params.Column5.Time = dto.ToDate
		params.Column5.Valid = true
	}

	payments, err := s.PaymentRepo.ListPayments(ctx.Request().Context(), params)
	if err != nil {
		return fmt.Errorf("failed to list payments: %v", err)
	}

	return utils.ResponseMessage(http.StatusOK, payments, ctx)
}

// RefundPayment processes a payment refund
func (s *ServicesHandler) RPayment(ctx echo.Context) error {
	// Get payment to validate it
	dto := ctx.Get(utils.ValidatedBodyDTO).(*domain.RefundPaymentDTO)
	payment, err := s.PaymentRepo.GetPayment(ctx.Request().Context(), dto.PaymentID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return fmt.Errorf("payment not found: %v", utils.ErrNotFound)
		}
		return fmt.Errorf("failed to get payment: %v", err)
	}

	// Validate payment can be refunded
	if payment.PaymentStatus != db.PaymentStatusSuccess {
		return fmt.Errorf("payment must be successful to be refunded: %v", utils.ErrBadRequest)
	}
	if payment.RefundAmount.Valid {
		return fmt.Errorf("payment has already been refunded: %v", utils.ErrBadRequest)
	}

	var paymentAmount float64
	err = payment.Amount.Scan(&paymentAmount)
	if err != nil {
		return fmt.Errorf("failed to parse payment amount: %w", err)
	}

	if dto.RefundAmount > paymentAmount {
		return fmt.Errorf("refund amount cannot be greater than payment amount: %v", utils.ErrBadRequest)
	}

	var refundAmount pgtype.Numeric
	if err := refundAmount.Scan(dto.RefundAmount); err != nil {
		return fmt.Errorf("invalid refund amount format: %v", utils.ErrBadRequest)
	}

	refundedByUUID, err := uuid.Parse(dto.RefundedBy)
	if err != nil {
		return fmt.Errorf("invalid refunded by ID: %v", utils.ErrBadRequest)
	}

	// Process refund
	refundedPayment, err := s.PaymentRepo.RefundPayment(ctx.Request().Context(), db.Refund_PaymentParams{
		ID:           dto.PaymentID,
		RefundAmount: refundAmount,
		RefundReason: pgtype.Text{String: dto.RefundReason, Valid: true},
		RefundedBy:   pgtype.UUID{Bytes: refundedByUUID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to process refund: %v", err)
	}

	return utils.ResponseMessage(http.StatusAccepted, refundedPayment, ctx)
}

// HandlePaymentWebhook processes payment webhooks from payment provider
func (s *ServicesHandler) HandlePaymentWebhook(ctx echo.Context) error {
	dto := ctx.Get(utils.ValidatedBodyDTO).(*domain.PaymentWebhookDTO)
	var pgMetadata []byte
	if dto.Metadata != nil {
		metadata, err := json.Marshal(dto.Metadata)
		if err != nil {
			return fmt.Errorf("invalid metadata format: %v", utils.ErrBadRequest)
		}
		pgMetadata = metadata
	}

	// Update payment status
	payment, err := s.PaymentRepo.UpdatePaymentStatus(ctx.Request().Context(), db.Update_Payment_StatusParams{
		ID:              dto.PaymentID,
		PaymentStatus:   db.PaymentStatus(dto.Status),
		TransactionID:   pgtype.Text{String: dto.TransactionID, Valid: true},
		PaymentMetadata: pgMetadata,
	})
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return fmt.Errorf("payment not found: %v", utils.ErrNotFound)
		}
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	return utils.ResponseMessage(http.StatusAccepted, payment, ctx)
}

// VerifyPayment verifies a payment using Paystack and updates the payment status
func (s *ServicesHandler) VerifyPayment(c echo.Context, reference string) error {
	// Start transaction
	tx, err := s.PaymentRepo.BeginWith(c.Request().Context())
	if err != nil {
		return &PaymentError{
			Code:    ErrCodeSystemError,
			Message: "failed to start transaction",
			Err:     err,
		}
	}
	defer tx.Rollback(c.Request().Context())

	// Verify and update payment with retries
	payment, err := s.verifyAndUpdatePaymentWithRetry(c.Request().Context(), reference)
	if err != nil {
		utils.Error("Payment verification failed",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "reference", Value: reference})
		return err
	}

	// Commit transaction
	if err := tx.Commit(c.Request().Context()); err != nil {
		return &PaymentError{
			Code:    ErrCodeSystemError,
			Message: "failed to commit transaction",
			Err:     err,
		}
	}

	utils.Info("Payment verified successfully",
		utils.LogField{Key: "payment_id", Value: payment.ID},
		utils.LogField{Key: "status", Value: payment.PaymentStatus})

	return utils.ResponseMessage(http.StatusOK, payment, c)
}

// verifyAndUpdatePaymentWithRetry handles payment verification and status updates with retry logic
func (s *ServicesHandler) verifyAndUpdatePaymentWithRetry(ctx context.Context, reference string) (*db.Payment, error) {
	var verificationResponse *domain.PaystackVerificationResponse
	var payment *db.Payment

	// Verify payment with provider using retry
	err := s.withRetry(ctx, "verify_payment", func() error {
		var err error
		paystackResponse, err := s.paymentService.VerifyTransaction(reference)
		if err != nil {
			return &PaymentError{
				Code:    ErrCodePaymentVerification,
				Message: "payment verification failed",
				Err:     err,
				Metadata: map[string]interface{}{
					"reference": reference,
				},
				Retryable: true,
			}
		}
		// Convert PaystackResponse.Data to domain.PaystackVerificationData
		verificationResponse = &domain.PaystackVerificationResponse{
			Status:  paystackResponse.Status,
			Message: paystackResponse.Message,
			Data: domain.PaystackVerificationData{
				Status:    paystackResponse.Data.Status,
				Reference: paystackResponse.Data.Reference,
				Amount:    paystackResponse.Data.Amount,
				Channel:   paystackResponse.Data.Channel,
				Currency:  paystackResponse.Data.Currency,
				PaidAt:    ParseTime(paystackResponse.Data.PaidAt),
				CreatedAt: ParseTime(paystackResponse.Data.CreatedAt),
				Customer: domain.PaystackCustomer{
					ID:        paystackResponse.Data.Customer.ID,
					FirstName: paystackResponse.Data.Customer.FirstName,
					LastName:  paystackResponse.Data.Customer.LastName,
					Email:     paystackResponse.Data.Customer.Email,
					Phone:     paystackResponse.Data.Customer.Phone,
				},
				Metadata: paystackResponse.Data.Metadata,
			},
		}
		if !verificationResponse.Status || verificationResponse.Data.Status != "success" {
			return &PaymentError{
				Code:    ErrCodePaymentVerification,
				Message: fmt.Sprintf("payment verification unsuccessful: %s", verificationResponse.Message),
				Err:     fmt.Errorf("invalid status: %s", verificationResponse.Data.Status),
				Metadata: map[string]interface{}{
					"reference":       reference,
					"response_status": verificationResponse.Status,
					"payment_status":  verificationResponse.Data.Status,
				},
				Retryable: verificationResponse.Data.Status == "pending",
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Get payment by reference using retry
	err = s.withRetry(ctx, "get_payment", func() error {
		var err error
		payment, err = s.PaymentRepo.GetPaymentByReference(ctx, reference)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				return &PaymentError{
					Code:    ErrCodePaymentNotFound,
					Message: "payment not found",
					Err:     err,
					Metadata: map[string]interface{}{
						"reference": reference,
					},
					Retryable: false,
				}
			}
			return &PaymentError{
				Code:    ErrCodeSystemError,
				Message: "failed to fetch payment",
				Err:     err,
				Metadata: map[string]interface{}{
					"reference": reference,
				},
				Retryable: true,
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Validate status transition
	if err := s.validatePaymentStatus(payment.PaymentStatus, db.PaymentStatusSuccess); err != nil {
		return nil, err
	}

	// Process and validate metadata
	metadataToUpdate := map[string]interface{}{
		"provider_reference": verificationResponse.Data.Reference,
		"transaction_time":   verificationResponse.Data.PaidAt.Format(time.RFC3339),
		"payment_method":     verificationResponse.Data.Channel,
		"currency":           verificationResponse.Data.Currency,
		"customer":           verificationResponse.Data.Customer,
		"original_metadata":  verificationResponse.Data.Metadata,
	}

	metadata, err := s.validateAndParseMetadata(metadataToUpdate)
	if err != nil {
		return nil, &PaymentError{
			Code:    ErrCodeValidation,
			Message: "failed to validate payment metadata",
			Err:     err,
			Metadata: map[string]interface{}{
				"reference": reference,
				"metadata":  metadataToUpdate,
			},
		}
	}

	// Update payment status with validated data
	updatedPayment, err := s.PaymentRepo.UpdatePaymentStatus(ctx, db.Update_Payment_StatusParams{
		ID:              payment.ID,
		PaymentStatus:   db.PaymentStatusSuccess,
		TransactionID:   pgtype.Text{String: verificationResponse.Data.Reference, Valid: true},
		PaymentMetadata: metadata,
	})
	if err != nil {
		return nil, &PaymentError{
			Code:    ErrCodeSystemError,
			Message: "failed to update payment status",
			Err:     err,
			Metadata: map[string]interface{}{
				"payment_id": payment.ID,
				"reference":  reference,
			},
			Retryable: true,
		}
	}

	return updatedPayment, nil
}
