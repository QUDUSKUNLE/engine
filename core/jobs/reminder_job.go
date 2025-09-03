package jobs

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/adapters/ex/templates/emails"
	"github.com/jackc/pgx/v5/pgtype"
	// "github.com/diagnoxix/core/utils"
)

func (j *ReminderJob) Start() {
	// Run every hour to check for appointments that need reminders
	j.scheduler.Every(1).Hour().Do(j.sendAppointmentReminders)
	j.scheduler.StartAsync()
}

func (j *ReminderJob) Stop() {
	j.scheduler.Stop()
}

func (j *ReminderJob) sendAppointmentReminders() {
	// ctx := context.Background()
	// now := time.Now()

	// // Get appointments in the next 24 hours that haven't had reminders sent
	// appointments, err := j.appointmentRepo.GetUpcomingAppointments(ctx, db.GetUpcomingAppointmentsParams{
	// 	StartTime: now,
	// 	EndTime:   now.Add(24 * time.Hour),
	// 	Statuses:  []db.AppointmentStatus{db.AppointmentStatusConfirmed},
	// })
	// if err != nil {
	// 	utils.Error("Failed to get upcoming appointments for reminders",
	// 		utils.LogField{Key: "error", Value: err.Error()})
	// 	return
	// }

	// for _, appointment := range appointments {
	// 	// Skip if reminder already sent (check a flag or separate table)
	// 	if appointment.ReminderSent {
	// 		continue
	// 	}

	// 	// Send reminder email
	// 	if err := j.sendReminderEmail(appointment); err != nil {
	// 		utils.Error("Failed to send reminder email",
	// 			utils.LogField{Key: "error", Value: err.Error()},
	// 			utils.LogField{Key: "appointment_id", Value: appointment.ID})
	// 		continue
	// 	}

	// 	// Mark reminder as sent
	// 	if err := j.appointmentRepo.MarkReminderSent(ctx, appointment.ID); err != nil {
	// 		utils.Error("Failed to mark reminder as sent",
	// 			utils.LogField{Key: "error", Value: err.Error()},
	// 			utils.LogField{Key: "appointment_id", Value: appointment.ID})
	// 	}
	// }
}

func (j *ReminderJob) SendReminderEmail(appointment *db.Appointment) error {
	// Get patient details by email
	patient, err := j.userRepo.GetUserByEmail(
		context.Background(),
		pgtype.Text{String: appointment.PatientID, Valid: true},
	)
	if err != nil {
		return err
	}

	// Get centre details
	centre, err := j.diagnosticRepo.GetDiagnosticCentre(
		context.Background(),
		appointment.DiagnosticCentreID,
	)
	if err != nil {
		return err
	}

	data := emails.AppointmentEmailData{
		EmailData: emails.EmailData{
			AppName: "Diagnoxix",
		},
		PatientName:     patient.Fullname.String,
		AppointmentID:   appointment.ID,
		AppointmentDate: appointment.AppointmentDate.Time,
		TimeSlot:        appointment.TimeSlot,
		CentreName:      centre.DiagnosticCentreName,
		// Status:          appointment.Status,
	}

	return j.notificationSvc.SendEmail(patient.Email.String, "Appointment Reminder", emails.TemplateAppointmentReminder, data)
}
