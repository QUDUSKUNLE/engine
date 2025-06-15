package repository

import (
	"context"

	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
)

// Ensure Repository implements AppointmentRepository
var _ ports.AppointmentRepository = (*Repository)(nil)

// CreateAppointment creates a new appointment
func (repo *Repository) CreateAppointment(ctx context.Context, params db.CreateAppointmentParams) (*db.Appointment, error) {
	return repo.database.CreateAppointment(ctx, params)
}

// GetAppointment retrieves an appointment by ID
func (repo *Repository) GetAppointment(ctx context.Context, id string) (*db.Appointment, error) {
	return repo.database.GetAppointment(ctx, id)
}

// ListAppointments lists appointments based on the provided parameters
func (repo *Repository) ListAppointments(ctx context.Context, params db.GetCentreAppointmentsParams) ([]*db.Appointment, error) {
	return repo.database.GetCentreAppointments(ctx, params)
}

// UpdateAppointment updates an appointment's details
func (repo *Repository) UpdateAppointment(ctx context.Context, params db.UpdateAppointmentPaymentParams) (*db.Appointment, error) {
	return repo.database.UpdateAppointmentPayment(ctx, params)
}

// CancelAppointment cancels an existing appointment
func (repo *Repository) CancelAppointment(ctx context.Context, id string) error {
	return repo.database.DeleteAppointment(ctx, id)
}

// RescheduleAppointment reschedules an appointment to a new time
func (repo *Repository) RescheduleAppointment(ctx context.Context, params db.RescheduleAppointmentParams) (*db.Appointment, error) {
	return repo.database.RescheduleAppointment(ctx, params)
}
