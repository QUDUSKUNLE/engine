package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
)

// RoutesAdaptor registers all API routes
func RoutesAdaptor(public *echo.Group, handler *handlers.HTTPHandler) *echo.Group {
	// Register all domain-specific routes
	AuthRoutes(public, handler)
	DiagnosticRoutes(public, handler)
	ScheduleRoutes(public, handler)
	MedicalRecordRoutes(public, handler)
	AppointmentRoutes(public, handler)
	PaymentRoutes(public, handler)
	AvailabilityRoutes(public, handler)
	return public
}
