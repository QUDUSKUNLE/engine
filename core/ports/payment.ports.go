package ports

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/adapters/ex/paystack"
)

// DBTX represents a database transaction interface
type DBTX interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	CreatePayment(
		ctx context.Context,
		payment db.Create_PaymentParams,
	) (*db.Payment, error)

	GetPayment(
		ctx context.Context,
		id string,
	) (*db.Payment, error)

	GetPaymentByReference(
		ctx context.Context,
		reference string,
	) (*db.Payment, error)

	ListPayments(
		ctx context.Context,
		params db.List_PaymentsParams,
	) ([]*db.Payment, error)

	RefundPayment(
		ctx context.Context,
		params db.Refund_PaymentParams,
	) (*db.Payment, error)

	UpdatePaymentStatus(
		ctx context.Context,
		params db.Update_Payment_StatusParams,
	) (*db.Payment, error)

	BeginWith(ctx context.Context) (DBTX, error)
}

// PaymentProviderService defines the interface for payment provider operations
type PaymentProviderService interface {
	InitializeTransaction(
		email,
		reference string,
		amount float64,
		metadata map[string]interface{},
	) (*paystack.PaystackTransactionResponse, error)
	VerifyTransaction(reference string) (*paystack.PaystackVerificationResponse, error)
}
