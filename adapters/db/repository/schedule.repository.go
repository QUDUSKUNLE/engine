package repository

import (
	"context"

	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
)

// Ensure Repository implements DiagnosticCentreRepository
var _ ports.ScheduleRepository = (*Repository)(nil)

// CreateDiagnosticSchedule creates a new diagnostic schedule
func (repo *Repository) CreateDiagnosticSchedule(ctx context.Context, arg db.Create_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error) {
	return repo.database.Create_Diagnostic_Schedule(ctx, arg)
}

// GetDiagnosticSchedule retrieves a diagnostic schedule by its parameters
func (repo *Repository) GetDiagnosticSchedule(ctx context.Context, req db.Get_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnostic_Schedule(ctx, req)
}

// GetDiagnosticSchedules retrieves multiple diagnostic schedules by parameters
func (repo *Repository) GetDiagnosticSchedules(ctx context.Context, req db.Get_Diagnostic_SchedulesParams) ([]*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnostic_Schedules(ctx, req)
}

// UpdateDiagnosticSchedule updates a diagnostic schedule
func (repo *Repository) UpdateDiagnosticSchedule(ctx context.Context, req db.Update_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error) {
	return repo.database.Update_Diagnostic_Schedule(ctx, req)
}

// DeleteDiagnosticSchedule deletes a diagnostic schedule
func (repo *Repository) DeleteDiagnosticSchedule(ctx context.Context, req db.Delete_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error) {
	return repo.database.Delete_Diagnostic_Schedule(ctx, req)
}

// GetDiagnosticScheduleByCentre retrieves a diagnostic schedule by centre parameters
func (repo *Repository) GetDiagnosticScheduleByCentre(ctx context.Context, req db.Get_Diagnsotic_Schedule_By_CentreParams) (*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnsotic_Schedule_By_Centre(ctx, req)
}

// GetDiagnosticSchedulesByCentre retrieves multiple diagnostic schedules by centre parameters
func (repo *Repository) GetDiagnosticSchedulesByCentre(ctx context.Context, req db.Get_Diagnsotic_Schedules_By_CentreParams) ([]*db.DiagnosticSchedule, error) {
	return repo.database.Get_Diagnsotic_Schedules_By_Centre(ctx, req)
}

// UpdateDiagnosticScheduleByCentre updates a diagnostic schedule by centre
func (repo *Repository) UpdateDiagnosticScheduleByCentre(ctx context.Context, req db.Update_Diagnostic_Schedule_By_CentreParams) (*db.DiagnosticSchedule, error) {
	return repo.database.Update_Diagnostic_Schedule_By_Centre(ctx, req)
}
