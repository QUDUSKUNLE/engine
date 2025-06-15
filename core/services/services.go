package services

import (
	"github.com/medicue/adapters/ex"
	"github.com/medicue/core/ports"
)

type ServicesHandler struct {
	UserRepo         ports.UserRepository
	ScheduleRepo     ports.ScheduleRepository
	DiagnosticRepo   ports.DiagnosticRepository
	RecordRepo       ports.RecordRepository
	AvailabilityRepo ports.AvailabilityRepository
	FileRepo         ports.FileService
	emailService     ports.EmailService
}

func ServicesAdapter(useRepo ports.UserRepository, scheduleRepo ports.ScheduleRepository, diagnosticCentreRepo ports.DiagnosticRepository, availabilityRepo ports.AvailabilityRepository, record ports.RecordRepository) *ServicesHandler {
	return &ServicesHandler{
		UserRepo:         useRepo,
		ScheduleRepo:     scheduleRepo,
		DiagnosticRepo:   diagnosticCentreRepo,
		AvailabilityRepo: availabilityRepo,
		RecordRepo:       record,
		FileRepo:         &ex.LocalFileService{},
		emailService:     &ex.EmailAdapter{},
	}
}
