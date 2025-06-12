package services

import (
	"github.com/medicue/adapters/ex"
	"github.com/medicue/core/ports"
)

type ServicesHandler struct {
	userRepo       ports.UserRepository
	scheduleRepo   ports.ScheduleRepository
	diagnosticRepo ports.DiagnosticRepository
	recordRepo     ports.RecordRepository
	fileRepo       ports.FileService
}
func ServicesAdapter(useRepo ports.UserRepository, scheduleRepo ports.ScheduleRepository, diagnosticCentreRepo ports.DiagnosticRepository, record ports.RecordRepository) *ServicesHandler {
	return &ServicesHandler{
		userRepo:       useRepo,
		scheduleRepo:   scheduleRepo,
		diagnosticRepo: diagnosticCentreRepo,
		recordRepo:     record,
		fileRepo:       &ex.LocalFileService{},
	}
}
