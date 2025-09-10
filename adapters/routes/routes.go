package routes

import (
	"github.com/diagnoxix/adapters/handlers"
	"github.com/labstack/echo/v4"
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
	AIRoutes(public, handler)
	CacheRoutes(public, handler)
	
	// WebSocket routes (separate group for ws prefix)
	wsGroup := public.Group("/ws")
	WebSocketRoutes(wsGroup, handler)
	
	return public
}
