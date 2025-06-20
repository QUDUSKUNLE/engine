package repository

import (
	"context"

	"github.com/medicue/adapters/db"
)

func (r *Repository) CreateAvailability(ctx context.Context, req db.Create_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error) {
	return r.database.Create_Availability(ctx, req)
}

func (r *Repository) CreateSingleAvailability(ctx context.Context, params db.Create_Single_AvailabilityParams) (*db.DiagnosticCentreAvailability, error) {
	return r.database.Create_Single_Availability(ctx, params)
}

func (r *Repository) DeleteAvailability(ctx context.Context, param db.Delete_AvailabilityParams) error {
	return r.database.Delete_Availability(ctx, param)
}

func (r *Repository) GetAvailability(ctx context.Context, params db.Get_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error) {
	return r.database.Get_Availability(ctx, params)
}

func (r *Repository) GetDiagnosticAvailability(ctx context.Context, req string) ([]*db.DiagnosticCentreAvailability, error) {
	return r.database.Get_Diagnostic_Availability(ctx, req)
}

func (r *Repository) UpdateAvailability(ctx context.Context, params db.Update_AvailabilityParams) (*db.DiagnosticCentreAvailability, error) {
	return r.database.Update_Availability(ctx, params)
}

func (r *Repository) UpdateManyAvailability(ctx context.Context, arg db.Update_Many_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error) {
	return r.database.Update_Many_Availability(ctx, arg)
}
