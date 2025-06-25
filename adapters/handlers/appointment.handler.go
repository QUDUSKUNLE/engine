package handlers

import (
	"github.com/labstack/echo/v4"
)

// CreateAppointment creates a new appointment for a diagnostic centre
// @Summary Create appointment
// @Description Create a new appointment for a diagnostic test
// @Tags Appointments
// @Accept json
// @Produce json
// @Param appointment body domain.CreateAppointmentDTO true "Appointment details"
// @Success 201 {object} db.Appointment
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 404 {object} utils.ErrorResponse "Diagnostic centre or schedule not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/appointments [post]
func (h *HTTPHandler) CreateAppointment(c echo.Context) error {
	// Validation happens in middleware
	return h.service.CreateAppointment(c)
}

// GetAppointment retrieves an appointment by ID
// @Summary Get appointment
// @Description Get an appointment by its ID
// @Tags Appointments
// @Accept json
// @Produce json
// @Param appointment_id path string true "Appointment ID"
// @Success 200 {object} db.Appointment
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Not authorized to view this appointment"
// @Failure 404 {object} utils.ErrorResponse "Appointment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/appointments/{appointment_id} [get]
func (h *HTTPHandler) GetAppointment(c echo.Context) error {
	// Validation happens in middleware
	return h.service.GetAppointment(c)
}

// ListAppointments lists appointments based on filters
// @Summary List appointments
// @Description List appointments with optional filters
// @Tags Appointments
// @Accept json
// @Produce json
// @Param diagnostic_centre_id query string false "Filter by diagnostic centre ID"
// @Param status query string false "Filter by status (pending, confirmed, in_progress, completed, cancelled, rescheduled)"
// @Param from_date query string false "Start date (RFC3339)"
// @Param to_date query string false "End date (RFC3339)"
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Success 200 {array} db.Appointment
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/appointments [get]
func (h *HTTPHandler) ListAppointments(c echo.Context) error {
	// Validation happens in middleware
	return h.service.ListAppointments(c)
}

// UpdateAppointment updates an existing appointment
// @Summary Update appointment
// @Description Update an existing appointment's details
// @Tags Appointments
// @Accept json
// @Produce json
// @Param appointment_id path string true "Appointment ID"
// @Param appointment body domain.UpdateAppointmentDTO true "Updated appointment details"
// @Success 200 {object} db.Appointment
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Not authorized to update this appointment"
// @Failure 404 {object} utils.ErrorResponse "Appointment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/appointments/{appointment_id} [put]
func (h *HTTPHandler) UpdateAppointment(c echo.Context) error {
	// Validation happens in middleware
	return nil
}

// CancelAppointment cancels an existing appointment
// @Summary Cancel appointment
// @Description Cancel an existing appointment
// @Tags Appointments
// @Accept json
// @Produce json
// @Param appointment_id path string true "Appointment ID"
// @Param cancellation body domain.CancelAppointmentDTO true "Cancellation details"
// @Success 200 {object} map[string]string "Cancellation success message"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data or appointment state"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Not authorized to cancel this appointment"
// @Failure 404 {object} utils.ErrorResponse "Appointment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/appointments/{appointment_id}/cancel [post]
func (h *HTTPHandler) CancelAppointment(c echo.Context) error {
	// Validation happens in middleware
	return h.service.CancelAppointment(c)
}

// RescheduleAppointment reschedules an existing appointment
// @Summary Reschedule appointment
// @Description Reschedule an existing appointment to a new time
// @Tags Appointments
// @Accept json
// @Produce json
// @Param appointment_id path string true "Appointment ID"
// @Param reschedule body domain.RescheduleAppointmentDTO true "Rescheduling details"
// @Success 200 {object} db.Appointment
// @Failure 400 {object} utils.ErrorResponse "Invalid input data or appointment state"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Not authorized to reschedule this appointment"
// @Failure 404 {object} utils.ErrorResponse "Appointment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/appointments/{appointment_id}/reschedule [post]
func (h *HTTPHandler) RescheduleAppointment(c echo.Context) error {
	// Validation happens in middleware
	return h.service.RescheduleAppointment(c)
}

func (h *HTTPHandler) ConfirmAppointment(c echo.Context) error {
	return h.service.ConfirmAppointment(c)
}
