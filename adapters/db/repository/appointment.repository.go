package repository

import (
	"context"
	"fmt"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
	"github.com/jackc/pgx/v5"
)

// Ensure Repository implements AppointmentRepository
var _ ports.AppointmentRepository = (*Repository)(nil)

func (r *Repository) IsValidTestType(ctx context.Context, testType string) bool {
	// List of valid test types - this should match the domain validation
	validTypes := []string{
		"BLOOD_TEST", "URINE_TEST", "X_RAY", "MRI", "CT_SCAN", "ULTRASOUND", "ECG",
		"EEG", "BIOPSY", "SKIN_TEST", "IMMUNOLOGY_TEST", "HORMONE_TEST", "VIRAL_TEST",
		"BACTERIAL_TEST", "PARASITIC_TEST", "FUNGAL_TEST", "MOLECULAR_TEST", "TOXICOLOGY_TEST",
		"ECHO", "COVID_19_TEST", "BLOOD_SUGAR_TEST", "LIPID_PROFILE", "HEMOGLOBIN_TEST",
		"THYROID_TEST", "LIVER_FUNCTION_TEST", "KIDNEY_FUNCTION_TEST", "URIC_ACID_TEST",
		"VITAMIN_D_TEST", "VITAMIN_B12_TEST", "HEMOGRAM", "COMPLETE_BLOOD_COUNT",
		"BLOOD_GROUPING", "HEPATITIS_B_TEST", "HEPATITIS_C_TEST", "HIV_TEST",
		"MALARIA_TEST", "DENGUE_TEST", "TYPHOID_TEST", "COVID_19_ANTIBODY_TEST",
		"COVID_19_RAPID_ANTIGEN_TEST", "COVID_19_RT_PCR_TEST", "PREGNANCY_TEST",
		"ALLERGY_TEST", "GENETIC_TEST", "OTHER",
	}

	for _, valid := range validTypes {
		if valid == testType {
			return true
		}
	}
	return false
}

// AppointmentTxRepository represents a transaction-aware repository
type AppointmentTxRepository struct {
	*Repository
	tx pgx.Tx
}

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

func (repo *Repository) MarkReminderSent(ctx context.Context, ID string) error {
	return nil
}

// RescheduleAppointment reschedules an appointment to a new time
func (repo *Repository) RescheduleAppointment(ctx context.Context, params db.RescheduleAppointmentParams) (*db.Appointment, error) {
	return repo.database.RescheduleAppointment(ctx, params)
}

// BeginTx starts a new transaction
func (r *Repository) BeginTx(ctx context.Context) (ports.AppointmentTx, error) {
	if r.conn == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	qtx := r.database.WithTx(tx)
	return &AppointmentTxRepository{
		Repository: &Repository{database: qtx},
		tx:         tx,
	}, nil
}

// CreateSchedule creates a schedule within the same transaction
func (t *AppointmentTxRepository) CreateSchedule(ctx context.Context, params db.Create_Diagnostic_ScheduleParams) (*db.DiagnosticSchedule, error) {
	return t.database.Create_Diagnostic_Schedule(ctx, params)
}

// Commit commits the transaction
func (t *AppointmentTxRepository) Commit(ctx context.Context) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Commit(ctx)
}

// Rollback rolls back the transaction
func (t *AppointmentTxRepository) Rollback(ctx context.Context) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Rollback(ctx)
}
