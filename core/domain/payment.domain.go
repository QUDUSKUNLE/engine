package domain

import (
	"time"

	"github.com/diagnoxix/adapters/db"
)

const (
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodWallet   PaymentMethod = "wallet"
)

const (
	PaymentProviderPaystack    PaymentProvider = "PAYSTACK"
	PaymentProviderFlutterwave PaymentProvider = "FLUTTERWAVE"
	PaymentProviderStripe      PaymentProvider = "STRIPE"
	PaymentProviderMonnify     PaymentProvider = "MONNIFY"
)

type (
	// PaymentProvider represents the payment provider
	PaymentProvider string

	// PaymentMethod represents the payment method
	PaymentMethod string
	// CreatePaymentDTO represents the request body for creating a payment
	CreatePaymentDTO struct {
		AppointmentID     string           `json:"appointment_id" validate:"required,uuid"`
		Amount            float64          `json:"amount" validate:"required,gt=0"`
		Currency          string           `json:"currency" validate:"required,len=3"`
		PaymentMethod     db.PaymentMethod `json:"payment_method" validate:"required,oneof=card transfer cash wallet"`
		PaymentProvider   string           `json:"payment_provider" validate:"required,oneof=PAYSTACK FLUTTERWAVE STRIPE MONNIFY"`
		PaymentMetadata   interface{}      `json:"payment_metadata,omitempty"`
		ProviderReference string           `json:"provider_reference,omitempty"`
	}
	// GetPaymentDTO represents the request parameters for getting a payment
	GetPaymentDTO struct {
		PaymentID string `param:"payment_id" validate:"required,uuid"`
	}
	// ListPaymentsDTO represents the query parameters for listing payments
	ListPaymentsDTO struct {
		DiagnosticCentreID string           `query:"diagnostic_centre_id" validate:"omitempty,uuid"`
		PatientID          string           `query:"patient_id" validate:"omitempty,uuid"`
		Status             db.PaymentStatus `query:"status" validate:"omitempty,oneof=pending success failed refunded cancelled"`
		FromDate           time.Time        `query:"from_date" validate:"omitempty"`
		ToDate             time.Time        `query:"to_date" validate:"omitempty,gtefield=FromDate"`
		PaginationQueryDTO
	}
	// RefundPaymentDTO represents the request body for refunding a payment
	RefundPaymentDTO struct {
		PaymentID    string  `param:"payment_id" validate:"required,uuid"`
		RefundAmount float64 `json:"refund_amount" validate:"required,gt=0"`
		RefundReason string  `json:"refund_reason" validate:"required,max=500"`
		RefundedBy   string  `json:"refunded_by"`
	}
	// PaymentWebhookDTO represents the request body for payment webhook
	PaymentWebhookDTO struct {
		PaymentID     string           `json:"payment_id" validate:"required"`
		Status        db.PaymentStatus `json:"status" validate:"required,oneof=pending success failed refunded cancelled"`
		TransactionID string           `json:"transaction_id" validate:"required"`
		Metadata      interface{}      `json:"metadata,omitempty"`
	}
	VerifyPaymentDTO struct {
		Reference string `param:"reference" validate:"required"`
	}
)
