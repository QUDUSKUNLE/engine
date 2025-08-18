package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// AppointmentRoutes registers all appointment-related routes
func AppointmentRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	appointmentGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/appointments",
			handler:     handler.CreateAppointment,
			factory:     func() interface{} { return &domain.CreateAppointmentDTO{} },
			description: "Create a new appointment",
		},
		{
			method:  http.MethodGet,
			path:    "/appointments/:appointment_id",
			handler: handler.GetAppointment,
			factory: func() interface{} {
				return &domain.GetAppointmentDTO{}
			},
			description: "Get appointment by appointment_id",
		},
		{
			method:  http.MethodGet,
			path:    "/appointments",
			handler: handler.ListAppointments,
			factory: func() interface{} {
				return &domain.ListAppointmentsDTO{}
			},
			description: "List appointments with filters",
		},
		{
			method:  http.MethodPost,
			path:    "/appointments/:appointment_id/cancel",
			handler: handler.CancelAppointment,
			factory: func() interface{} {
				return &domain.CancelAppointmentDTO{}
			},
			description: "Cancel an appointment",
		},
		{
			method:      http.MethodPost,
			path:        "/appointments/:appointment_id/reschedule",
			handler:     handler.RescheduleAppointment,
			factory:     func() interface{} { return &domain.RescheduleAppointmentDTO{} },
			description: "Reschedule an appointment",
		},
		{
			method:      http.MethodPost,
			path:        "/appointments/confirm_appointment",
			handler:     handler.ConfirmAppointment,
			factory:     func() interface{} { return &domain.ConfirmAppointmentDTO{} },
			description: "Confirm an appointment",
		},
	}

	registerRoutes(group, appointmentGroup)
}
