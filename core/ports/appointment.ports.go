package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

// AppointmentRepository defines the interface for appointment data operations
type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, appointment db.CreateAppointmentParams) (*db.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*db.Appointment, error)
	ListAppointments(ctx context.Context, params db.GetCentreAppointmentsParams) ([]*db.Appointment, error)
	UpdateAppointment(ctx context.Context, params db.UpdateAppointmentPaymentParams) (*db.Appointment, error)
	CancelAppointment(ctx context.Context, id string) error
	RescheduleAppointment(ctx context.Context, params db.RescheduleAppointmentParams) (*db.Appointment, error)
	MarkReminderSent(ctx context.Context, id string) error
}
