package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/adapters/ex/templates"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

// CreateAppointment creates a new appointment and associated schedule
func (service *ServicesHandler) CreateAppointment(context echo.Context) error {

	// Authenticate and authorize user - owner or manager only
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Parse appointment creation request
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.CreateAppointmentDTO)

	// Start transaction
	tx, err := service.AppointmentRepo.BeginTx(context.Request().Context())
	if err != nil {
		utils.Error("Failed to start transaction",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	defer tx.Rollback(context.Request().Context())

	// Normalize test type format
	testType := strings.ToUpper(strings.ReplaceAll(string(dto.TestType), " ", "_"))
	if !service.AppointmentRepo.IsValidTestType(context.Request().Context(), testType) {
		utils.Error("Failed to validate test_type",
			utils.LogField{Key: "error", Value: "Invalid test type"})
		return utils.ErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf("invalid test type: %s", testType),
			context,
		)
	}
	// Use the validated test type for schedule creation
	scheduleParams := db.Create_Diagnostic_ScheduleParams{
		UserID:             currentUser.UserID.String(),
		DiagnosticCentreID: dto.DiagnosticCentreID.String(),
		ScheduleTime:       toTimestamptz(dto.AppointmentDate),
		Doctor:             string(dto.PreferredDoctor),
		TestType:           testType,
		Notes:              toText(dto.Notes),
		AcceptanceStatus:   db.ScheduleAcceptanceStatusACCEPTED, // Auto-accept the schedule since validation is done
	}

	schedule, err := tx.CreateSchedule(context.Request().Context(), scheduleParams)
	if err != nil {
		utils.Error("Failed to create schedule",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Create appointment with confirmed status since schedule is auto-accepted
	appointmentParams := db.CreateAppointmentParams{
		PatientID:          currentUser.UserID.String(),
		ScheduleID:         schedule.ID,
		DiagnosticCentreID: dto.DiagnosticCentreID.String(),
		AppointmentDate:    toTimestamptz(dto.AppointmentDate),
		TimeSlot: fmt.Sprintf("%s-%s",
			dto.AppointmentDate.Format("15:04"),
			dto.AppointmentDate.Add(30*time.Minute).Format("15:04")),
		Status: db.AppointmentStatusPending,
		Notes:  toText(dto.Notes),
	}

	appointment, err := tx.CreateAppointment(context.Request().Context(), appointmentParams)
	if err != nil {
		utils.Error("Failed to create appointment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "patient_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Create payment record
	var pgAmount pgtype.Numeric
	// Format amount with exactly 2 decimal places and apply proper rounding
	amount := fmt.Sprintf("%.2f", dto.Amount)
	if err := pgAmount.Scan(amount); err != nil {
		utils.Error("Failed to convert amount to numeric",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "amount", Value: amount})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid amount format"), context)
	}

	// Generate unique reference for this payment
	paymentReference := fmt.Sprintf("MED-%s-%s", appointment.ID, time.Now().Format("20060102150405"))
	// Initialize Paystack transaction
	metadata := map[string]interface{}{
		"appointment_id": appointment.ID,
		"patient_id":     currentUser.UserID.String(),
		"centre_id":      dto.DiagnosticCentreID.String(),
		"test_type":      testType,
	}

	// Get User Email
	user, err := service.UserRepo.GetUser(context.Request().Context(), currentUser.UserID.String())

	if err != nil {
		utils.Error("Failed to get user email",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "appointment_id", Value: appointment.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	paystackResponse, err := service.paymentService.InitializeTransaction(
		user.Email.String,
		dto.Amount,
		paymentReference,
		metadata,
	)

	if err != nil {
		utils.Error("Failed to initialize payment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "appointment_id", Value: appointment.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("Payment initialized successfully",
		utils.LogField{Key: "paystack_response", Value: paystackResponse},
		utils.LogField{Key: "appointment_id", Value: appointment.ID},
	)

	metadataBytes, err := utils.MarshalJSONField(paystackResponse, context)
	if err != nil {
		utils.Error("Failed to marshal payment metadata",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	_, err = service.PaymentRepo.CreatePayment(context.Request().Context(), db.Create_PaymentParams{
		AppointmentID:      appointment.ID,
		PatientID:          currentUser.UserID.String(),
		DiagnosticCentreID: dto.DiagnosticCentreID.String(),
		Amount:             pgAmount,
		Currency:           "NGN",
		PaymentMethod:      db.PaymentMethodCard,
		PaymentProvider:    db.PaymentProviderPAYSTACK,
		PaymentMetadata:    metadataBytes,
		ProviderMetadata:   metadataBytes,
		ProviderReference:  pgtype.Text{String: paymentReference, Valid: true},
		TransactionID:      pgtype.Text{String: paystackResponse.Data.Reference, Valid: true},
	})

	if err != nil {
		utils.Error("Failed to create payment",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "appointment_id", Value: appointment.ID},
		)
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Commit transaction
	if err := tx.Commit(context.Request().Context()); err != nil {
		utils.Error("Failed to commit transaction",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, errors.New("failed to commit transaction"), context)
	}

	return utils.ResponseMessage(http.StatusCreated, map[string]interface{}{
		"message":     "Appointment created successfully",
		"appointment": appointment,
		"reference":   paystackResponse.Data,
	}, context)
}

// GetAppointment retrieves an appointment by ID
func (service *ServicesHandler) GetAppointment(context echo.Context) error {
	// Authentication check
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
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
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
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
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
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
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
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
	patient, err := service.UserRepo.GetUser(
		context.Background(),
		appointment.PatientID,
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

	// Get schedule details for test type
	schedule, err := service.ScheduleRepo.GetDiagnosticScheduleByCentre(
		context.Background(),
		db.Get_Diagnsotic_Schedule_By_CentreParams{
			ID:                 appointment.ScheduleID,
			DiagnosticCentreID: appointment.DiagnosticCentreID,
		},
	)
	if err != nil {
		utils.Error("Failed to get schedule details for confirmation email",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	data := templates.AppointmentEmailData{
		EmailData: templates.EmailData{
			AppName: "Medicue",
		},
		PatientName:     patient.Fullname.String,
		AppointmentID:   appointment.ID,
		AppointmentDate: appointment.AppointmentDate.Time,
		TimeSlot:        appointment.TimeSlot,
		CentreName:      centre.DiagnosticCentreName,
		TestType:        string(schedule.TestType),
		Notes:           appointment.Notes.String,
	}

	// Get email template
	emailBody, err := templates.GetAppointmentConfirmationTemplate(data)
	if err != nil {
		utils.Error("Failed to generate confirmation email template",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Send email
	if err := service.notificationService.SendEmail(patient.Email.String, "Appointment Confirmation", emailBody); err != nil {
		utils.Error("Failed to send confirmation email",
			utils.LogField{Key: "error", Value: err.Error()})
	}
}

func (service *ServicesHandler) sendAppointmentCancellationEmail(appointment *db.Appointment) {
	// Get patient details by email
	patient, err := service.UserRepo.GetUser(context.Background(), appointment.PatientID)
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
	patient, err := service.UserRepo.GetUser(context.Background(), appointment.PatientID)
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

// sendAppointmentRequestEmail sends an email to the patient about their appointment request
// func (service *ServicesHandler) sendAppointmentRequestEmail(appointment *db.Appointment) {
// 	// Get patient details
// 	patient, err := service.UserRepo.GetUser(context.Background(), appointment.PatientID)
// 	if err != nil {
// 		utils.Error("Failed to get patient details for request email",
// 			utils.LogField{Key: "error", Value: err.Error()})
// 		return
// 	}

// 	// Get centre details
// 	centre, err := service.DiagnosticRepo.GetDiagnosticCentre(context.Background(), appointment.DiagnosticCentreID)
// 	if err != nil {
// 		utils.Error("Failed to get centre details for request email",
// 			utils.LogField{Key: "error", Value: err.Error()})
// 		return
// 	}

// 	data := templates.AppointmentEmailData{
// 		EmailData: templates.EmailData{
// 			AppName: "Medicue",
// 			Title:   "Medicue - Your Diagnostic Appointment",
// 		},
// 		PatientName:     patient.Fullname.String,
// 		AppointmentID:   appointment.ID,
// 		AppointmentDate: appointment.AppointmentDate.Time,
// 		TimeSlot:        appointment.TimeSlot,
// 		CentreName:      centre.DiagnosticCentreName,
// 		Status:          string(appointment.Status),
// 		Notes:           appointment.Notes.String,
// 	}

// 	body, err := templates.GetAppointmentRequestTemplate(data)
// 	if err != nil {
// 		utils.Error("Failed to generate request email template",
// 			utils.LogField{Key: "error", Value: err.Error()})
// 		return
// 	}

// 	if err := service.notificationService.SendEmail(patient.Email.String, "Appointment Request Received", body); err != nil {
// 		utils.Error("Failed to send request email",
// 			utils.LogField{Key: "error", Value: err.Error()})
// 	}
// }

// notifyDiagnosticCentreOfNewAppointment notifies the diagnostic centre about a new appointment
func (service *ServicesHandler) notifyDiagnosticCentreOfNewAppointment(appointment *db.Appointment, centre *db.DiagnosticCentre) {
	// Get schedule details to get test type and doctor preference
	schedule, err := service.ScheduleRepo.GetDiagnosticScheduleByCentre(context.Background(), db.Get_Diagnsotic_Schedule_By_CentreParams{
		ID:                 appointment.ScheduleID,
		DiagnosticCentreID: appointment.DiagnosticCentreID,
	})
	if err != nil {
		utils.Error("Failed to get schedule details for centre notification",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Get patient details
	patient, err := service.UserRepo.GetUserByEmail(
		context.Background(),
		pgtype.Text{String: appointment.PatientID, Valid: true},
	)
	if err != nil {
		utils.Error("Failed to get patient details for centre notification",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Extract contact information
	var contact domain.Contact
	if err := utils.UnmarshalJSONField(centre.Contact, &contact, nil); err != nil {
		utils.Error("Failed to unmarshal centre contact",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	data := templates.AppointmentEmailData{
		EmailData: templates.EmailData{
			AppName: "Medicue",
		},
		StaffName:       contact.Email, // Use contact email as staff name
		PatientName:     patient.Fullname.String,
		AppointmentID:   appointment.ID,
		AppointmentDate: appointment.AppointmentDate.Time,
		TimeSlot:        appointment.TimeSlot,
		CentreName:      centre.DiagnosticCentreName,
		TestType:        string(schedule.TestType),
		SpecialNotes:    schedule.Notes.String,
		RequiredAction:  "New appointment confirmation received",
	}

	body, err := templates.GetStaffNotificationTemplate(data)
	if err != nil {
		utils.Error("Failed to generate staff notification template",
			utils.LogField{Key: "error", Value: err.Error()})
		return
	}

	// Send email to diagnostic centre's primary email
	if err := service.notificationService.SendEmail(contact.Email, "New Appointment Confirmation", body); err != nil {
		utils.Error("Failed to send centre notification email",
			utils.LogField{Key: "error", Value: err.Error()})
	}

	// If configured, also send SMS notifications to all phone numbers
	for _, phone := range contact.Phone {
		message := fmt.Sprintf(
			"New appointment received for %s on %s at %s. Patient: %s. Test: %s",
			centre.DiagnosticCentreName,
			appointment.AppointmentDate.Time.Format("Jan 2, 2006"),
			appointment.TimeSlot,
			patient.Fullname.String,
			schedule.TestType,
		)
		if err := service.notificationService.SendSMS(phone, message); err != nil {
			utils.Error("Failed to send SMS notification",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "phone", Value: phone})
		}
	}
}
