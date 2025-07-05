package ports

import (
	"context"

	"github.com/medivue/adapters/db"
)

// AppointmentRepository defines the interface for appointment data operations
type AppointmentRepository interface {
	BeginTx(ctx context.Context) (AppointmentTx, error)
	CreateAppointment(ctx context.Context, appointment db.CreateAppointmentParams) (*db.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*db.Appointment, error)
	ListAppointments(ctx context.Context, params db.GetCentreAppointmentsParams) ([]*db.Appointment, error)
	UpdateAppointment(ctx context.Context, params db.UpdateAppointmentPaymentParams) (*db.Appointment, error)
	CancelAppointment(ctx context.Context, id string) error
	RescheduleAppointment(ctx context.Context, params db.RescheduleAppointmentParams) (*db.Appointment, error)
	MarkReminderSent(ctx context.Context, id string) error
	IsValidTestType(ctx context.Context, testType string) bool
}

// AppointmentTx represents a transaction for appointment operations
type AppointmentTx interface {
	CreateAppointment(ctx context.Context, appointment db.CreateAppointmentParams) (*db.Appointment, error)
	CreateSchedule(ctx context.Context, schedule db.Create_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error)
	CreatePayment(ctx context.Context, payment db.Create_PaymentParams) (*db.Payment, error)
	UpdateAppointment(ctx context.Context, update db.UpdateAppointmentPaymentParams) (*db.Appointment, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
