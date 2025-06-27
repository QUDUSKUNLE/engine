package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/medivue/adapters/handlers"
)

// RoutesAdaptor registers all API routes
func RoutesAdaptor(public *echo.Group, handler *handlers.HTTPHandler) *echo.Group {
	// Health check route (no auth required)
	public.GET("/health", handler.HealthCheck)

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
