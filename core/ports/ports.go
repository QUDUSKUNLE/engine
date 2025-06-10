package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

type ScheduleRepository interface {
	CreateDiagnosticSchedule(ctx context.Context, arg db.Create_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error)
	GetDiagnosticSchedule(ctx context.Context, req db.Get_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error)
	GetDiagnosticSchedules(ctx context.Context, req db.Get_Diagnostic_SchedulesParams) ([]*db.DiagnosticSchedule, error)
	UpdateDiagnosticSchedule(ctx context.Context, req db.Update_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error)
	DeleteDiagnosticSchedule(ctx context.Context, req db.Delete_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error)
	GetDiagnosticScheduleByCentre(ctx context.Context, req db.Get_Diagnsotic_Schedule_By_CentreParams) (*db.DiagnosticSchedule, error)
	GetDiagnosticSchedulesByCentre(ctx context.Context, req db.Get_Diagnsotic_Schedules_By_CentreParams) ([]*db.DiagnosticSchedule, error)
	UpdateDiagnosticScheduleByCentre(ctx context.Context, req db.Update_Diagnostic_Schedule_By_CentreParams) (*db.DiagnosticSchedule, error)
}

type UserRepository interface {
	GetUsers(ctx context.Context, params db.GetUsersParams) ([]*db.User, error)
	CreateUser(ctx context.Context, user db.CreateUserParams) (*db.CreateUserRow, error)
	GetUser(ctx context.Context, id string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	UpdateUser(ctx context.Context, user db.UpdateUserParams) (*db.UpdateUserRow, error)
}

type DiagnosticRepository interface {
	CreateDiagnosticCentre(ctx context.Context, arg db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error)
	GetDiagnosticCentre(ctx context.Context, id string) (*db.DiagnosticCentre, error)
	GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error)
	UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error)
	// Add other methods as needed
}
