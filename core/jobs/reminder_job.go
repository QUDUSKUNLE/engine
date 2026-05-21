package jobs

import (
	"context"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/adapters/external/templates/emails"
	"github.com/diagnoxix/core/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (job *ReminderJob) Start() {
	// Run every hour to check for appointments that need reminders
	job.scheduler.Every(1).Hour().Do(job.sendAppointmentReminders)
	job.scheduler.StartAsync()

	utils.Info("Reminder job manager started")
}

func (job *ReminderJob) Stop() {
	job.scheduler.Stop()
}

func (job *ReminderJob) sendAppointmentReminders() {
	ctx := context.Background()
	now := time.Now()

	utils.Info("Reminder job for the next 1 hour is about to run...")

	// Get appointments in the next 24 hours that haven't had reminders sent
	appointments, err := job.appointmentRepo.GetUpcomingAppointments(ctx, db.GetUpComingAppointmentsParams{
		AppointmentDate:   pgtype.Timestamptz{Time: now.Add(24 * time.Hour), Valid: true },
		Column1:  []string{string(db.AppointmentStatusConfirmed)},
	})
	if err != nil {
		utils.Error("Failed to get upcoming appointments for reminders",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	for _, appointment := range appointments {
		// Skip if reminder already sent (check a flag or separate table)
		if appointment.ReminderSent.Valid {
			continue
		}           
		// Send reminder email
		if err := job.SendReminderEmail(appointment); err != nil {
			utils.Error("Failed to send reminder email",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "appointment_id", Value: appointment.ID})
			continue
		}

		// Mark reminder as sent
		if err := job.appointmentRepo.MarkReminderSent(ctx, appointment.ID); err != nil {
			utils.Error("Failed to mark reminder as sent",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "appointment_id", Value: appointment.ID})
		}
	}

	utils.Info("Finished reminder job for the next 1 hour...")
}

func (job *ReminderJob) SendReminderEmail(param *db.GetUpComingAppointmentsRow) error {

	
	data := emails.AppointmentEmailData{
		EmailData: emails.EmailData{
			AppName: "Diagnoxix",
		},
		PatientName:     param.FullName.String,
		AppointmentID:   param.ID,
		AppointmentDate: param.AppointmentDate.Time,
		TimeSlot:        param.TimeSlot,
		CentreName:      param.DiagnosticCentreName.String,
		Status:          string(param.Status),
	}

	return job.notificationSvc.SendEmail(param.Email.String, "Appointment Reminder", emails.TemplateAppointmentReminder, data)
}
