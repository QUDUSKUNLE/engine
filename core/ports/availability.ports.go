package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

// AvailabilityRepository defines the interface for availability data operations
type AvailabilityRepository interface {
	CreateAvailability(ctx context.Context, availability db.Create_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error)
	UpdateAvailability(ctx context.Context, update db.Update_AvailabilityParams) (*db.DiagnosticCentreAvailability, error)
	GetAvailability(ctx context.Context, params db.Get_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error)
	DeleteAvailability(ctx context.Context, req db.Delete_AvailabilityParams) error
	UpdateManyAvailability(ctx context.Context, arg db.Update_Many_AvailabilityParams) ([]*db.DiagnosticCentreAvailability, error)
}
