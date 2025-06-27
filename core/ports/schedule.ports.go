package ports

import (
	"context"

	"github.com/medivue/adapters/db"
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
