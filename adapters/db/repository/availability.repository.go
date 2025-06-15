package repository

import (
	"context"

	"github.com/medicue/adapters/db"
)

func (r *Repository) CreateAvailability(ctx context.Context, req db.Create_AvailabilityParams) (*db.DiagnosticCentreAvailability, error) {
	return r.database.Create_Availability(ctx, req)
}

func (r *Repository) DeleteAvailability(ctx context.Context, param db.Delete_AvailabilityParams) error {

	return r.database.Delete_Availability(ctx, param)
}

func (r *Repository) GetAvailability(ctx context.Context, params db.Get_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error) {
	// TODO: Implement the actual availability retrieval logic
	return nil, nil
}

func (r *Repository) UpdateAvailability(ctx context.Context, params db.Update_AvailabilityParams) (*db.DiagnosticCentreAvailability, error) {
	return r.database.Update_Availability(ctx, params)
}
