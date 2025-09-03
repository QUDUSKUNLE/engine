package ports

import (
	"context"

	"github.com/diagnoxix/adapters/db"
)

// DiagnosticRepository defines the interface for diagnostic centre data access
type DiagnosticRepository interface {
	BeginDiagnostic(
		ctx context.Context,
	) (DiagnosticTx, error)

	// Existing methods
	CreateDiagnosticCentre(
		ctx context.Context,
		arg db.Create_Diagnostic_CentreParams,
	) (*db.DiagnosticCentre, error)

	GetDiagnosticCentre(
		ctx context.Context,
		id string,
	) (*db.Get_Diagnostic_CentreRow, error)

	UpdateDiagnosticCentreByOwner(
		ctx context.Context,
		params db.Update_Diagnostic_Centre_ByOwnerParams,
	) (*db.DiagnosticCentre, error)

	GetDiagnosticCentreByManager(
		ctx context.Context,
		params db.Get_Diagnostic_Centre_ByManagerParams,
	) (*db.Get_Diagnostic_Centre_ByManagerRow, error)

	GetNearestDiagnosticCentres(
		ctx context.Context,
		params db.Get_Nearest_Diagnostic_CentresParams,
	) ([]*db.Get_Nearest_Diagnostic_CentresRow, error)

	// New methods
	DeleteDiagnosticCentreByOwner(
		ctx context.Context,
		params db.Delete_Diagnostic_Centre_ByOwnerParams,
	) (*db.DiagnosticCentre, error)

	GetDiagnosticCentreByOwner(
		ctx context.Context,
		params db.Get_Diagnostic_Centre_ByOwnerParams,
	) (*db.Get_Diagnostic_Centre_ByOwnerRow, error)

	ListDiagnosticCentresByOwner(
		ctx context.Context,
		params db.List_Diagnostic_Centres_ByOwnerParams,
	) ([]*db.List_Diagnostic_Centres_ByOwnerRow, error)

	RetrieveDiagnosticCentres(
		ctx context.Context,
		params db.Retrieve_Diagnostic_CentresParams,
	) ([]*db.Retrieve_Diagnostic_CentresRow, error)

	SearchDiagnosticCentres(
		ctx context.Context,
		params db.Search_Diagnostic_CentresParams,
	) ([]*db.Search_Diagnostic_CentresRow, error)

	AssignAdmin(
		ctx context.Context,
		params db.AssignAdminParams,
	) (*db.DiagnosticCentre, error)

	UnAssignAdmin(
		ctx context.Context,
		arg db.UnassignAdminParams,
	) (*db.DiagnosticCentre, error)
}

type TestPriceRepository interface {
	CreateTestPrice(
		ctx context.Context,
		test_price db.Create_Test_PriceParams,
	) ([]*db.DiagnosticCentreTestPrice, error)
}

type DiagnosticTx interface {
	DBTX
	CreateDiagnosticCentre(
		ctx context.Context,
		arg db.Create_Diagnostic_CentreParams,
	) (*db.DiagnosticCentre, error)
	CreateTestPrice(
		ctx context.Context,
		test_price db.Create_Test_PriceParams,
	) ([]*db.DiagnosticCentreTestPrice, error)
}
