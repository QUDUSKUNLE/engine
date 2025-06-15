package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

// Helper functions for type conversions
func toTimestamptz(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	ts.Time = t
	ts.Valid = !t.IsZero()
	return ts
}

func toText(s string) pgtype.Text {
	var text pgtype.Text
	text.String = s
	text.Valid = s != ""
	return text
}

func toNumeric(n float64) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(n)
	return num
}

func toUUID(id string) pgtype.UUID {
	var uid pgtype.UUID
	if parsed, err := uuid.Parse(id); err == nil {
		uid.Bytes = parsed
		uid.Valid = true
	}
	return uid
}

// CreateAppointment creates a new appointment
func (service *ServicesHandler) CreateAppointment(context echo.Context) error {
	// Authentication check for registered users
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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

	return utils.ResponseMessage(http.StatusCreated, appointment, context)
}

// GetAppointment retrieves an appointment by ID
func (service *ServicesHandler) GetAppointment(context echo.Context) error {
	// Authentication check
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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
		Limit:             int32(dto.PageSize),
		Offset:            int32((dto.Page - 1) * dto.PageSize),
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
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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
	if err := service.AppointmentRepo.CancelAppointment(context.Request().Context(), dto.AppointmentID); err != nil {
		utils.Error("Failed to cancel appointment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "appointment_id", Value: dto.AppointmentID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]string{"message": "Appointment cancelled successfully"}, context)
}

// RescheduleAppointment reschedules an appointment to a new time
func (service *ServicesHandler) RescheduleAppointment(context echo.Context) error {
	// Authentication check
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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

	return utils.ResponseMessage(http.StatusOK, rescheduledAppointment, context)
}
