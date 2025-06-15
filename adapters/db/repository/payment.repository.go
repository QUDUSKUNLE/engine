package repository

import (
	"context"

	"github.com/medicue/adapters/db"
)

// CreatePayment creates a new payment record
func (r *Repository) CreatePayment(ctx context.Context, payment db.Create_PaymentParams) (*db.Payment, error) {
	return r.database.Create_Payment(ctx, payment)
}

// GetPayment retrieves a payment by ID
func (r *Repository) GetPayment(ctx context.Context, id string) (*db.Payment, error) {
	return r.database.Get_Payment(ctx, id)
}

// ListPayments retrieves a list of payments with filtering
func (r *Repository) ListPayments(ctx context.Context, params db.List_PaymentsParams) ([]*db.Payment, error) {
	return r.database.List_Payments(ctx, params)
}

// RefundPayment processes a payment refund
func (r *Repository) RefundPayment(ctx context.Context, params db.Refund_PaymentParams) (*db.Payment, error) {
	return r.database.Refund_Payment(ctx, params)
}

// UpdatePaymentStatus updates the payment status and related information
func (r *Repository) UpdatePaymentStatus(ctx context.Context, params db.Update_Payment_StatusParams) (*db.Payment, error) {
	return r.database.Update_Payment_Status(ctx, params)
}
