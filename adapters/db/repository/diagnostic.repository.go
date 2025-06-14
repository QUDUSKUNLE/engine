package repository

import (
	"context"

	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
)

// Ensure Repository implements DiagnosticCentreRepository
var _ ports.DiagnosticRepository = (*Repository)(nil)

func (repo *Repository) CreateDiagnosticCentre(ctx context.Context, params db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error) {
	return repo.database.Create_Diagnostic_Centre(ctx, params)
}

func (repo *Repository) UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	return repo.database.Update_Diagnostic_Centre_ByOwner(ctx, params)
}

func (repo *Repository) GetDiagnosticCentre(ctx context.Context, diaagnostic_centre_id string) (*db.DiagnosticCentre, error) {
	return repo.database.Get_Diagnostic_Centre(ctx, diaagnostic_centre_id)
}

func (repo *Repository) GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.DiagnosticCentre, error) {
	return repo.database.Get_Diagnostic_Centre_ByManager(ctx, params)
}

func (repo *Repository) GetDiagnosticCentreByOwner(ctx context.Context, params db.Get_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	return repo.database.Get_Diagnostic_Centre_ByOwner(ctx, params)
}

func (repo *Repository) GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error) {
	return repo.database.Get_Nearest_Diagnostic_Centres(ctx, params)
}

func (repo *Repository) ListDiagnosticCentresByOwner(ctx context.Context, params db.List_Diagnostic_Centres_ByOwnerParams) ([]*db.DiagnosticCentre, error) {
	return repo.database.List_Diagnostic_Centres_ByOwner(ctx, params)
}

func (repo *Repository) RetrieveDiagnosticCentres(ctx context.Context, params db.Retrieve_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error) {
	return repo.database.Retrieve_Diagnostic_Centres(ctx, params)
}

func (repo *Repository) SearchDiagnosticCentres(ctx context.Context, params db.Search_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error) {
	return repo.database.Search_Diagnostic_Centres(ctx, params)
}

func (repo *Repository) SearchDiagnosticCentresByDoctor(ctx context.Context, params db.Search_Diagnostic_Centres_ByDoctorParams) ([]*db.DiagnosticCentre, error) {
	return repo.database.Search_Diagnostic_Centres_ByDoctor(ctx, params)
}

func (repo *Repository) DeleteDiagnosticCentreByOwner(ctx context.Context, params db.Delete_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	return repo.database.Delete_Diagnostic_Centre_ByOwner(ctx, params)
}
