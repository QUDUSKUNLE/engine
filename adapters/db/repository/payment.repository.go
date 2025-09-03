package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Ensure Repository implements PaymentRepository
var _ ports.PaymentRepository = (*Repository)(nil)

func (repo *Repository) CreatePayment(
	ctx context.Context,
	payment db.Create_PaymentParams,
) (*db.Payment, error) {
	return repo.database.Create_Payment(ctx, payment)
}

func (repo *Repository) GetPayment(
	ctx context.Context,
	id string,
) (*db.Payment, error) {
	return repo.database.Get_Payment(ctx, id)
}

func (repo *Repository) ListPayments(
	ctx context.Context,
	params db.List_PaymentsParams,
) ([]*db.Payment, error) {
	return repo.database.List_Payments(ctx, params)
}

func (repo *Repository) RefundPayment(
	ctx context.Context,
	params db.Refund_PaymentParams,
) (*db.Payment, error) {
	return repo.database.Refund_Payment(ctx, params)
}

func (repo *Repository) UpdatePaymentStatus(
	ctx context.Context,
	params db.Update_Payment_StatusParams,
) (*db.Payment, error) {
	return repo.database.Update_Payment_Status(ctx, params)
}

func (repo *Repository) GetPaymentByReference(
	ctx context.Context,
	reference string,
) (*db.Payment, error) {
	return repo.database.GetPaymentByReference(ctx, pgtype.Text{
		String: reference,
		Valid:  true,
	})
}

func (repo *Repository) BeginWith(
	ctx context.Context,
) (ports.DBTX, error) {
	return repo.conn.BeginTx(ctx, pgx.TxOptions{})
}
