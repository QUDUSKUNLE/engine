package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/medivue/adapters/db"
	"github.com/medivue/core/ports"
)

// Ensure Repository implements DiagnosticCentreRepository
var _ ports.DiagnosticRepository = (*Repository)(nil)

// AppointmentTxRepository represents a transaction-aware repository
type DiagnosticTxRepository struct {
	*Repository
	tx pgx.Tx
}

func (repo *Repository) CreateDiagnosticCentre(ctx context.Context, params db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error) {
	return repo.database.Create_Diagnostic_Centre(ctx, params)
}

func (repo *Repository) UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	return repo.database.Update_Diagnostic_Centre_ByOwner(ctx, params)
}

func (repo *Repository) GetDiagnosticCentre(ctx context.Context, diaagnostic_centre_id string) (*db.Get_Diagnostic_CentreRow, error) {
	return repo.database.Get_Diagnostic_Centre(ctx, diaagnostic_centre_id)
}

func (repo *Repository) GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.Get_Diagnostic_Centre_ByManagerRow, error) {
	return repo.database.Get_Diagnostic_Centre_ByManager(ctx, params)
}

func (repo *Repository) GetDiagnosticCentreByOwner(ctx context.Context, params db.Get_Diagnostic_Centre_ByOwnerParams) (*db.Get_Diagnostic_Centre_ByOwnerRow, error) {
	return repo.database.Get_Diagnostic_Centre_ByOwner(ctx, params)
}

func (repo *Repository) GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error) {
	return repo.database.Get_Nearest_Diagnostic_Centres(ctx, params)
}

func (repo *Repository) ListDiagnosticCentresByOwner(ctx context.Context, params db.List_Diagnostic_Centres_ByOwnerParams) ([]*db.List_Diagnostic_Centres_ByOwnerRow, error) {
	return repo.database.List_Diagnostic_Centres_ByOwner(ctx, params)
}

func (repo *Repository) RetrieveDiagnosticCentres(ctx context.Context, params db.Retrieve_Diagnostic_CentresParams) ([]*db.Retrieve_Diagnostic_CentresRow, error) {
	return repo.database.Retrieve_Diagnostic_Centres(ctx, params)
}

func (repo *Repository) SearchDiagnosticCentres(ctx context.Context, params db.Search_Diagnostic_CentresParams) ([]*db.Search_Diagnostic_CentresRow, error) {
	return repo.database.Search_Diagnostic_Centres(ctx, params)
}

func (repo *Repository) SearchDiagnosticCentresByDoctor(ctx context.Context, params db.Search_Diagnostic_Centres_ByDoctorParams) ([]*db.Search_Diagnostic_Centres_ByDoctorRow, error) {
	return repo.database.Search_Diagnostic_Centres_ByDoctor(ctx, params)
}

func (repo *Repository) DeleteDiagnosticCentreByOwner(ctx context.Context, params db.Delete_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	return repo.database.Delete_Diagnostic_Centre_ByOwner(ctx, params)
}

// BeginTx starts a new transaction
func (r *Repository) BeginDiagnostic(ctx context.Context) (ports.DiagnosticTx, error) {
	if r.conn == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	qtx := r.database.WithTx(tx)
	return &DiagnosticTxRepository{
		Repository: &Repository{database: qtx},
		tx:         tx,
	}, nil
}

// Commit commits the transaction
func (t *DiagnosticTxRepository) Commit(ctx context.Context) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Commit(ctx)
}

// Rollback rolls back the transaction
func (t *DiagnosticTxRepository) Rollback(ctx context.Context) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Rollback(ctx)
}
