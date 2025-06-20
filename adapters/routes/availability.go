package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/core/domain"
)

// AvailabilityRoutes registers all availability-related routes
func AvailabilityRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	availabilityGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/availability",
			handler:     handler.CreateAvailability,
			factory:     func() interface{} { return &domain.CreateAvailabilityDTO{} },
			description: "Create availability slots",
		},
		{
			method:      http.MethodGet,
			path:        "/availability/:diagnostic_centre_id",
			handler:     handler.GetAvailability,
			description: "Get availability slots",
		},
		{
			method:      http.MethodPut,
			path:        "/availability/:diagnostic_centre_id/:day_of_week",
			handler:     handler.UpdateAvailability,
			factory:     func() interface{} { return &domain.UpdateAvailabilityDTO{} },
			description: "Update availability slot",
		},
		{
			method:      http.MethodPut,
			path:        "/availability/:diagnostic_centre_id",
			handler:     handler.UpdateManyAvailability,
			factory:     func() interface{} { return &domain.UpdateManyAvailabilityDTO{} },
			description: "Update multiple availability slots",
		},
		{
			method:      http.MethodDelete,
			path:        "/availability/:diagnostic_centre_id/:day_of_week",
			handler:     handler.DeleteAvailability,
			description: "Delete availability slot",
		},
	}

	registerRoutes(group, availabilityGroup)
}
