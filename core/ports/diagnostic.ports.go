package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

type DiagnosticRepository interface {
	CreateDiagnosticCentre(ctx context.Context, arg db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentre(ctx context.Context, id string) (*db.DiagnosticCentre, error)
	GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error)
	UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.DiagnosticCentre, error)
	// Add other methods as needed
}
