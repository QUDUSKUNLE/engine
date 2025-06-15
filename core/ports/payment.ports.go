package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment db.Create_PaymentParams) (*db.Payment, error)
	GetPayment(ctx context.Context, id string) (*db.Payment, error)
	ListPayments(ctx context.Context, params db.List_PaymentsParams) ([]*db.Payment, error)
	RefundPayment(ctx context.Context, params db.Refund_PaymentParams) (*db.Payment, error)
	UpdatePaymentStatus(ctx context.Context, params db.Update_Payment_StatusParams) (*db.Payment, error)
}
