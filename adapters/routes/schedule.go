package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/core/domain"
)

// ScheduleRoutes registers all schedule-related routes
func ScheduleRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	scheduleGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/diagnostic_schedules",
			handler:     handler.CreateSchedule,
			factory:     func() interface{} { return &domain.CreateScheduleDTO{} },
			description: "Create schedule",
		},
		{
			method:      http.MethodGet,
			path:        "/diagnostic_schedules/:schedule_id",
			handler:     handler.GetSchedule,
			factory:     func() interface{} { return &domain.GetDiagnosticScheduleParamDTO{} },
			description: "Get schedule details",
		},
		{
			method:      http.MethodPut,
			path:        "/diagnostic_schedules/:schedule_id",
			handler:     handler.UpdateSchedule,
			factory:     func() interface{} { return &domain.UpdateScheduleDTO{} },
			description: "Update a schedule",
		},
		{
			method:      http.MethodGet,
			path:        "/diagnostic_schedules",
			handler:     handler.GetSchedules,
			factory:     func() interface{} { return &domain.GetDiagnosticSchedulesQueryDTO{} },
			description: "Get all schedules",
		},
	}

	registerRoutes(group, scheduleGroup)
}
