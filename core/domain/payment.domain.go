package domain

import (
	"time"

	"github.com/medicue/adapters/db"
)

// // PaymentStatus represents the status of a payment
// type PaymentStatus string

// const (
// 	PaymentStatusPending   PaymentStatus = "pending"
// 	PaymentStatusSuccess   PaymentStatus = "success"
// 	PaymentStatusFailed    PaymentStatus = "failed"
// 	PaymentStatusRefunded  PaymentStatus = "refunded"
// 	PaymentStatusCancelled PaymentStatus = "cancelled"
// )

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodWallet   PaymentMethod = "wallet"
)

// CreatePaymentDTO represents the request body for creating a payment
type CreatePaymentDTO struct {
	AppointmentID   string           `json:"appointment_id" validate:"required,uuid"`
	Amount          float64          `json:"amount" validate:"required,gt=0"`
	Currency        string           `json:"currency" validate:"required,len=3"`
	PaymentMethod   db.PaymentMethod `json:"payment_method" validate:"required,oneof=card transfer cash wallet"`
	PaymentMetadata interface{}      `json:"payment_metadata,omitempty"`
}

// GetPaymentDTO represents the request parameters for getting a payment
type GetPaymentDTO struct {
	PaymentID string `param:"payment_id" validate:"required,uuid"`
}

// ListPaymentsDTO represents the query parameters for listing payments
type ListPaymentsDTO struct {
	DiagnosticCentreID string           `query:"diagnostic_centre_id" validate:"omitempty,uuid"`
	PatientID          string           `query:"patient_id" validate:"omitempty,uuid"`
	Status             db.PaymentStatus `query:"status" validate:"omitempty,oneof=pending success failed refunded cancelled"`
	FromDate           time.Time        `query:"from_date" validate:"omitempty"`
	ToDate             time.Time        `query:"to_date" validate:"omitempty,gtefield=FromDate"`
	Page               int              `query:"page" validate:"min=1"`
	PageSize           int              `query:"page_size" validate:"min=1,max=100"`
}

// RefundPaymentDTO represents the request body for refunding a payment
type RefundPaymentDTO struct {
	PaymentID    string  `param:"payment_id" validate:"required,uuid"`
	RefundAmount float64 `json:"refund_amount" validate:"required,gt=0"`
	RefundReason string  `json:"refund_reason" validate:"required,max=500"`
	RefundedBy   string  `json:"refunded_by"`
}

// PaymentWebhookDTO represents the request body for payment webhook
type PaymentWebhookDTO struct {
	PaymentID     string           `json:"payment_id" validate:"required"`
	Status        db.PaymentStatus `json:"status" validate:"required,oneof=pending success failed refunded cancelled"`
	TransactionID string           `json:"transaction_id" validate:"required"`
	Metadata      interface{}      `json:"metadata,omitempty"`
}
