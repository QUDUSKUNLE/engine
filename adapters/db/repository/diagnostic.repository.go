package repository

import (
	"context"

	"github.com/medicue/adapters/db"
)

// DiagnosticCentreRepository abstracts diagnostic centre DB operations
//go:generate mockgen -destination=../../mocks/mock_diagnostic_centre_repository.go -package=mocks . DiagnosticCentreRepository

// DiagnosticCentreRepository defines the interface for diagnostic centre DB operations
// This allows for easier testing and swapping implementations
// You can use this interface in your service layer

type DiagnosticCentreRepository interface {
	CreateDiagnosticCentre(ctx context.Context, params db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error)
	UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentre(ctx context.Context, diaagnostic_centre_id string) (*db.DiagnosticCentre, error)
	GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentreByOwner(ctx context.Context, params db.Get_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error)
	ListDiagnosticCentresByOwner(ctx context.Context, params db.List_Diagnostic_Centres_ByOwnerParams) ([]*db.DiagnosticCentre, error)
	RetrieveDiagnosticCentres(ctx context.Context, params db.Retrieve_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error)
	SearchDiagnosticCentres(ctx context.Context, params db.Search_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error)
	SearchDiagnosticCentresByDoctor(ctx context.Context, params db.Search_Diagnostic_Centres_ByDoctorParams) ([]*db.DiagnosticCentre, error)
}

// Ensure Repository implements DiagnosticCentreRepository
var _ DiagnosticCentreRepository = (*Repository)(nil)

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
