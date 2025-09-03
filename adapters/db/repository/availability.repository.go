package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
)

// Ensure Repository implements AvailabilityRepository
var _ ports.AvailabilityRepository = (*Repository)(nil)

func (repo *Repository) CreateAvailability(
	ctx context.Context,
	req db.Create_AvailabilityParams,
) ([]*db.DiagnosticCentreAvailability, error) {
	return repo.database.Create_Availability(ctx, req)
}

func (repo *Repository) CreateSingleAvailability(
	ctx context.Context,
	params db.Create_Single_AvailabilityParams,
) (*db.DiagnosticCentreAvailability, error) {
	return repo.database.Create_Single_Availability(ctx, params)
}

func (repo *Repository) DeleteAvailability(
	ctx context.Context,
	param db.Delete_AvailabilityParams,
) error {
	return repo.database.Delete_Availability(ctx, param)
}

func (repo *Repository) GetAvailability(
	ctx context.Context,
	params db.Get_AvailabilityParams,
) ([]*db.DiagnosticCentreAvailability, error) {
	return repo.database.Get_Availability(ctx, params)
}

func (repo *Repository) GetDiagnosticAvailability(
	ctx context.Context,
	req string,
) ([]*db.DiagnosticCentreAvailability, error) {
	return repo.database.Get_Diagnostic_Availability(ctx, req)
}

func (repo *Repository) UpdateAvailability(
	ctx context.Context,
	params db.Update_AvailabilityParams,
) (*db.DiagnosticCentreAvailability, error) {
	return repo.database.Update_Availability(ctx, params)
}

func (repo *Repository) UpdateManyAvailability(
	ctx context.Context,
	arg db.Update_Many_AvailabilityParams,
) ([]*db.DiagnosticCentreAvailability, error) {
	return repo.database.Update_Many_Availability(ctx, arg)
}
