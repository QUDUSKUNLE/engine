package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/adapters/ex/templates"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

// CreateAppointment creates a new appointment
func (service *ServicesHandler) CreateAppointment(context echo.Context) error {
	// Authentication check for registered users
	currentUser, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get validated DTO
	dto := context.Get(utils.ValidatedBodyDTO).(*domain.CreateAppointmentDTO)

	// Verify diagnostic centre exists
	_, err = service.DiagnosticRepo.GetDiagnosticCentre(context.Request().Context(), dto.DiagnosticCentreID)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found"), context)
	}

	// Verify schedule exists and is valid
	schedule, err := service.ScheduleRepo.GetDiagnosticScheduleByCentre(context.Request().Context(), db.Get_Diagnsotic_Schedule_By_CentreParams{
		ID:                 dto.ScheduleID,
		DiagnosticCentreID: dto.DiagnosticCentreID,
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("schedule not found"), context)
	}

	if schedule.AcceptanceStatus != db.ScheduleAcceptanceStatusACCEPTED {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("schedule is not accepted"), context)
	}

	// Create appointment
	params := db.CreateAppointmentParams{
		PatientID:          currentUser.UserID.String(),
		ScheduleID:         dto.ScheduleID,
		DiagnosticCentreID: dto.DiagnosticCentreID,
		AppointmentDate:    toTimestamptz(dto.AppointmentDate),
		TimeSlot:           dto.TimeSlot,
		Status:             db.AppointmentStatusPending,
		Notes:              toText(dto.Notes),
	}

	appointment, err := service.AppointmentRepo.CreateAppointment(context.Request().Context(), params)
	if err != nil {
		utils.Error("Failed to create appointment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "patient_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Send confirmation email asynchronously
	go service.sendAppointmentConfirmationEmail(appointment)

	return utils.ResponseMessage(http.StatusCreated, appointment, context)
}

// GetAppointment retrieves an appointment by ID
func (service *ServicesHandler) GetAppointment(context echo.Context) error {
	// Authentication check
	currentUser, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get validated DTO
	dto := context.Get(utils.ValidatedBodyDTO).(*domain.GetAppointmentDTO)

	// Get appointment
	appointment, err := service.AppointmentRepo.GetAppointment(context.Request().Context(), dto.AppointmentID)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("appointment not found"), context)
	}

	// Verify ownership
	if appointment.PatientID != currentUser.UserID.String() {
		return utils.ErrorResponse(http.StatusForbidden, errors.New("not authorized to view this appointment"), context)
	}

	return utils.ResponseMessage(http.StatusOK, appointment, context)
}

// ListAppointments lists appointments based on filters
func (service *ServicesHandler) ListAppointments(context echo.Context) error {
	// Authentication check
	currentUser, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get validated DTO
	dto := context.Get(utils.ValidatedBodyDTO).(*domain.ListAppointmentsDTO)

	// If no date range specified, default to next 30 days
	if dto.FromDate.IsZero() {
		dto.FromDate = time.Now()
	}
	if dto.ToDate.IsZero() {
		dto.ToDate = dto.FromDate.AddDate(0, 1, 0) // 1 month from FromDate
	}

	// Force patient ID to current user unless they are a centre manager
	if currentUser.UserType != db.UserEnumDIAGNOSTICCENTREMANAGER &&
		currentUser.UserType != db.UserEnumDIAGNOSTICCENTREOWNER {
		dto.PatientID = currentUser.UserID.String()
	}

	// Build status array
	var statuses []db.AppointmentStatus
	if dto.Status != "" {
		statuses = append(statuses, db.AppointmentStatus(dto.Status))
	}

	// Get appointments
	params := db.GetCentreAppointmentsParams{
		DiagnosticCentreID: dto.DiagnosticCentreID,
		Column2:            statuses, // Status array
		AppointmentDate:    toTimestamptz(dto.FromDate),
		AppointmentDate_2:  toTimestamptz(dto.ToDate),
		Limit:              int32(dto.PageSize),
		Offset:             int32((dto.Page - 1) * dto.PageSize),
	}

	appointments, err := service.AppointmentRepo.ListAppointments(context.Request().Context(), params)
	if err != nil {
		utils.Error("Failed to list appointments",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, appointments, context)
}

// CancelAppointment cancels an existing appointment
func (service *ServicesHandler) CancelAppointment(context echo.Context) error {
	// Authentication check
	currentUser, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get validated DTO
	dto := context.Get(utils.ValidatedBodyDTO).(*domain.CancelAppointmentDTO)

	// Verify appointment exists and belongs to user
	appointment, err := service.AppointmentRepo.GetAppointment(context.Request().Context(), dto.AppointmentID)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("appointment not found"), context)
	}

	if appointment.PatientID != currentUser.UserID.String() {
		return utils.ErrorResponse(http.StatusForbidden, errors.New("not authorized to cancel this appointment"), context)
	}

	// Verify appointment can be cancelled
	if appointment.Status != db.AppointmentStatusPending && appointment.Status != db.AppointmentStatusConfirmed {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("appointment cannot be cancelled in its current state"), context)
	}

	// Cancel appointment
	_ = db.CancelAppointmentParams{
		ID:                 dto.AppointmentID,
		CancellationReason: toText(dto.Reason),
		CancelledBy:        toUUID(currentUser.UserID.String()),
		CancellationFee:    toNumeric(0), // Fee could be configured based on business rules
	}

	err = service.AppointmentRepo.CancelAppointment(context.Request().Context(), dto.AppointmentID)
	if err != nil {
		return err
	}

	cancelledAppointment, err := service.AppointmentRepo.GetAppointment(context.Request().Context(), dto.AppointmentID)
	if err != nil {
		utils.Error("Failed to cancel appointment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "appointment_id", Value: dto.AppointmentID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Send cancellation email asynchronously
	go service.sendAppointmentCancellationEmail(cancelledAppointment)

	return utils.ResponseMessage(http.StatusOK, map[string]string{"message": "Appointment cancelled successfully"}, context)
}

// RescheduleAppointment reschedules an appointment to a new time
func (service *ServicesHandler) RescheduleAppointment(context echo.Context) error {
	// Authentication check
	currentUser, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get validated DTO
	dto := context.Get(utils.ValidatedBodyDTO).(*domain.RescheduleAppointmentDTO)

	// Verify appointment exists and belongs to user
	appointment, err := service.AppointmentRepo.GetAppointment(context.Request().Context(), dto.AppointmentID)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("appointment not found"), context)
	}

	if appointment.PatientID != currentUser.UserID.String() {
		return utils.ErrorResponse(http.StatusForbidden, errors.New("not authorized to reschedule this appointment"), context)
	}

	// Verify appointment can be rescheduled
	if appointment.Status != db.AppointmentStatusPending && appointment.Status != db.AppointmentStatusConfirmed {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("appointment cannot be rescheduled in its current state"), context)
	}

	// Verify new schedule exists and is valid
	newSchedule, err := service.ScheduleRepo.GetDiagnosticScheduleByCentre(context.Request().Context(), db.Get_Diagnsotic_Schedule_By_CentreParams{
		ID:                 dto.NewScheduleID,
		DiagnosticCentreID: appointment.DiagnosticCentreID,
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("new schedule not found"), context)
	}

	if newSchedule.AcceptanceStatus != db.ScheduleAcceptanceStatusACCEPTED {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("new schedule is not accepted"), context)
	}

	// Reschedule appointment
	params := db.RescheduleAppointmentParams{
		ID:                 dto.AppointmentID,
		ReschedulingReason: toText(dto.RescheduleReason),
		RescheduledBy:      toUUID(currentUser.UserID.String()),
		ReschedulingFee:    toNumeric(0), // Fee could be configured based on business rules
		ScheduleID:         dto.NewScheduleID,
		AppointmentDate:    toTimestamptz(dto.NewDate),
		TimeSlot:           dto.NewTimeSlot,
	}

	rescheduledAppointment, err := service.AppointmentRepo.RescheduleAppointment(context.Request().Context(), params)
	if err != nil {
		utils.Error("Failed to reschedule appointment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "appointment_id", Value: dto.AppointmentID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Send reschedule email asynchronously
	go service.sendAppointmentRescheduleEmail(rescheduledAppointment)

	return utils.ResponseMessage(http.StatusOK, rescheduledAppointment, context)
}

// Helper functions to send appointment emails
func (service *ServicesHandler) sendAppointmentConfirmationEmail(appointment *db.Appointment) {
	// Get patient details by email
	patient, err := service.UserRepo.GetUserByEmail(
		context.Background(),
		pgtype.Text{String: appointment.PatientID, Valid: true},
	)
	if err != nil {
		utils.Error("Failed to get patient details for confirmation email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Get centre details
	centre, err := service.DiagnosticRepo.GetDiagnosticCentre(
		context.Background(),
		appointment.DiagnosticCentreID,
	)
	if err != nil {
		utils.Error("Failed to get centre details for confirmation email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	data := templates.AppointmentEmailData{
		EmailData: templates.EmailData{
			AppName: "Medicue",
			// AppURL:  os.Getenv("APP_URL"),
		},
		PatientName:     patient.Fullname.String,
		AppointmentID:   appointment.ID,
		AppointmentDate: appointment.AppointmentDate.Time,
		TimeSlot:        appointment.TimeSlot,
		CentreName:      centre.DiagnosticCentreName,
		// Status:          appointment.Status,
		Notes: appointment.Notes.String,
	}

	body, err := templates.GetAppointmentConfirmationTemplate(data)
	if err != nil {
		utils.Error("Failed to generate confirmation email template",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	if err := service.notificationService.SendEmail(patient.Email.String, "Appointment Confirmation", body); err != nil {
		utils.Error("Failed to send confirmation email",
			utils.LogField{Key: "error", Value: err.Error()})
	}
}

func (service *ServicesHandler) sendAppointmentCancellationEmail(appointment *db.Appointment) {
	// Get patient details by email
	patient, err := service.UserRepo.GetUserByEmail(
		context.Background(),
		pgtype.Text{String: appointment.PatientID, Valid: true},
	)
	if err != nil {
		utils.Error("Failed to get patient details for cancellation email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Get centre details
	centre, err := service.DiagnosticRepo.GetDiagnosticCentre(
		context.Background(),
		appointment.DiagnosticCentreID,
	)
	if err != nil {
		utils.Error("Failed to get centre details for cancellation email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	data := templates.AppointmentEmailData{
		EmailData: templates.EmailData{
			AppName: "Medicue",
			// AppURL:  os.Getenv("APP_URL"),
		},
		PatientName:     patient.Fullname.String,
		AppointmentID:   appointment.ID,
		AppointmentDate: appointment.AppointmentDate.Time,
		TimeSlot:        appointment.TimeSlot,
		CentreName:      centre.DiagnosticCentreName,
		// Status:          appointment.Status,
	}

	body, err := templates.GetAppointmentCancellationTemplate(data)
	if err != nil {
		utils.Error("Failed to generate cancellation email template",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	if err := service.notificationService.SendEmail(patient.Email.String, "Appointment Cancelled", body); err != nil {
		utils.Error("Failed to send cancellation email",
			utils.LogField{Key: "error", Value: err.Error()})
	}
}

func (service *ServicesHandler) sendAppointmentRescheduleEmail(appointment *db.Appointment) {
	// Get patient details by email
	patient, err := service.UserRepo.GetUserByEmail(
		context.Background(),
		pgtype.Text{String: appointment.PatientID, Valid: true},
	)
	if err != nil {
		utils.Error("Failed to get patient details for reschedule email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Get centre details
	centre, err := service.DiagnosticRepo.GetDiagnosticCentre(
		context.Background(),
		appointment.DiagnosticCentreID,
	)
	if err != nil {
		utils.Error("Failed to get centre details for reschedule email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	data := templates.AppointmentEmailData{
		EmailData: templates.EmailData{
			AppName: "Medicue",
			// AppURL:  os.Getenv("APP_URL"),
		},
		PatientName:     patient.Fullname.String,
		AppointmentID:   appointment.ID,
		AppointmentDate: appointment.AppointmentDate.Time,
		TimeSlot:        appointment.TimeSlot,
		CentreName:      centre.DiagnosticCentreName,
		// Status:          appointment.Status,
	}

	body, err := templates.GetAppointmentRescheduleTemplate(data)
	if err != nil {
		utils.Error("Failed to generate reschedule email template",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	if err := service.notificationService.SendEmail(patient.Email.String, "Appointment Rescheduled", body); err != nil {
		utils.Error("Failed to send reschedule email",
			utils.LogField{Key: "error", Value: err.Error()})
	}
}
