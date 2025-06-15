package services

import (
	"github.com/medicue/adapters/ex"
	"github.com/medicue/core/ports"
)

type ServicesHandler struct {
	UserRepo            ports.UserRepository
	ScheduleRepo        ports.ScheduleRepository
	DiagnosticRepo      ports.DiagnosticRepository
	RecordRepo          ports.RecordRepository
	AvailabilityRepo    ports.AvailabilityRepository
	PaymentRepo         ports.PaymentRepository
	AppointmentRepo     ports.AppointmentRepository
	FileRepo            ports.FileService
	notificationService ports.NotificationService
}

func ServicesAdapter(
	useRepo ports.UserRepository,
	scheduleRepo ports.ScheduleRepository,
	diagnosticCentreRepo ports.DiagnosticRepository,
	availabilityRepo ports.AvailabilityRepository,
	record ports.RecordRepository,
	paymentPort ports.PaymentRepository,
	appointmentPort ports.AppointmentRepository,
) *ServicesHandler {
	return &ServicesHandler{
		UserRepo:            useRepo,
		ScheduleRepo:        scheduleRepo,
		DiagnosticRepo:      diagnosticCentreRepo,
		AvailabilityRepo:    availabilityRepo,
		PaymentRepo:         paymentPort,
		AppointmentRepo:     appointmentPort,
		RecordRepo:          record,
		FileRepo:            &ex.LocalFileService{},
		notificationService: &ex.EmailAdapter{},
	}
}
