package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

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
