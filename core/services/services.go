package services

import (
	"github.com/medicue/core/ports"
)

type ServicesHandler struct {
	UserRepo             ports.UserRepository
	ScheduleRepo         ports.ScheduleRepository
	DiagnosticCentreRepo ports.DiagnosticRepository
}

func ServicesAdapter(userRepo ports.UserRepository, scheduleRepo ports.ScheduleRepository, diagnosticCentreRepo ports.DiagnosticRepository) *ServicesHandler {
	return &ServicesHandler{
		UserRepo:             userRepo,
		ScheduleRepo:         scheduleRepo,
		DiagnosticCentreRepo: diagnosticCentreRepo,
	}
}
