package services

import (
	"github.com/medicue/adapters/ex"
	"github.com/medicue/core/ports"
)

type ServicesHandler struct {
	UserRepo       ports.UserRepository
	ScheduleRepo   ports.ScheduleRepository
	DiagnosticRepo ports.DiagnosticRepository
	RecordRepo     ports.RecordRepository
	FileRepo       ports.FileService
	emailService   ports.EmailService
}

func ServicesAdapter(useRepo ports.UserRepository, scheduleRepo ports.ScheduleRepository, diagnosticCentreRepo ports.DiagnosticRepository, record ports.RecordRepository) *ServicesHandler {
	return &ServicesHandler{
		UserRepo:       useRepo,
		ScheduleRepo:   scheduleRepo,
		DiagnosticRepo: diagnosticCentreRepo,
		RecordRepo:     record,
		FileRepo:       &ex.LocalFileService{},
		emailService:   &ex.EmailAdapter{},
	}
}
