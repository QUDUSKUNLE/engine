package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
)

// Ensure Repository implements DiagnosticCentreRepository
var _ ports.ScheduleRepository = (*Repository)(nil)

func (repo *Repository) CreateDiagnosticSchedule(
	ctx context.Context,
	arg db.Create_Diagnostic_ScheduleParams,
) (*db.DiagnosticSchedule, error) {
	return repo.database.Create_Diagnostic_Schedule(ctx, arg)
}

func (repo *Repository) GetDiagnosticSchedule(
	ctx context.Context,
	req db.Get_Diagnostic_ScheduleParams,
) (*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnostic_Schedule(ctx, req)
}

func (repo *Repository) GetDiagnosticSchedules(
	ctx context.Context,
	req db.Get_Diagnostic_SchedulesParams,
) ([]*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnostic_Schedules(ctx, req)
}

func (repo *Repository) UpdateDiagnosticSchedule(
	ctx context.Context,
	req db.Update_Diagnostic_ScheduleParams,
) (*db.DiagnosticSchedule, error) {
	return repo.database.Update_Diagnostic_Schedule(ctx, req)
}

func (repo *Repository) DeleteDiagnosticSchedule(
	ctx context.Context,
	req db.Delete_Diagnostic_ScheduleParams,
) (*db.DiagnosticSchedule, error) {
	return repo.database.Delete_Diagnostic_Schedule(ctx, req)
}

func (repo *Repository) GetDiagnosticScheduleByCentre(
	ctx context.Context,
	req db.Get_Diagnsotic_Schedule_By_CentreParams,
) (*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnsotic_Schedule_By_Centre(ctx, req)
}

func (repo *Repository) GetDiagnosticSchedulesByCentre(
	ctx context.Context,
	req db.Get_Diagnsotic_Schedules_By_CentreParams,
) ([]*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnsotic_Schedules_By_Centre(ctx, req)
}

func (repo *Repository) UpdateDiagnosticScheduleByCentre(
	ctx context.Context,
	req db.Update_Diagnostic_Schedule_By_CentreParams,
) (*db.DiagnosticSchedule, error) {
	return repo.database.Update_Diagnostic_Schedule_By_Centre(ctx, req)
}
