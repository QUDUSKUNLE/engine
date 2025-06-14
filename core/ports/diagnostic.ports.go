package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

// DiagnosticRepository defines the interface for diagnostic centre data access
type DiagnosticRepository interface {
	// Existing methods
	CreateDiagnosticCentre(ctx context.Context, arg db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentre(ctx context.Context, id string) (*db.DiagnosticCentre, error)
	UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.DiagnosticCentre, error)
	GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error)

	// New methods
	DeleteDiagnosticCentreByOwner(ctx context.Context, params db.Delete_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentreByOwner(ctx context.Context, params db.Get_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	ListDiagnosticCentresByOwner(ctx context.Context, params db.List_Diagnostic_Centres_ByOwnerParams) ([]*db.DiagnosticCentre, error)
	RetrieveDiagnosticCentres(ctx context.Context, params db.Retrieve_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error)
	SearchDiagnosticCentres(ctx context.Context, params db.Search_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error)
}
