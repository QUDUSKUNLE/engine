package jobs

import (
	"time"

	"github.com/diagnoxix/core/ports"
	"github.com/go-co-op/gocron"
)

type ReminderJob struct {
	appointmentRepo ports.AppointmentRepository
	diagnosticRepo  ports.DiagnosticRepository
	notificationSvc ports.NotificationService
	userRepo        ports.UserRepository
	scheduler       *gocron.Scheduler
}

func NewReminderJob(
	appointmentRepo ports.AppointmentRepository,
	notificationSvc ports.NotificationService,
	userRepo ports.UserRepository,
	diagnosticRepo ports.DiagnosticRepository,
) *ReminderJob {
	return &ReminderJob{
		appointmentRepo: appointmentRepo,
		notificationSvc: notificationSvc,
		userRepo:        userRepo,
		diagnosticRepo:  diagnosticRepo,
		scheduler:       gocron.NewScheduler(time.UTC),
	}
}
