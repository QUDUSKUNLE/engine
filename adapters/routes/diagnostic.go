package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// DiagnosticRoutes registers all diagnostic centre-related routes
func DiagnosticRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	diagnosticGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/diagnostic_centres",
			handler:     handler.CreateDiagnostic,
			factory:     func() interface{} { return &domain.CreateDiagnosticDTO{} },
			description: "Create a diagnostic centre",
		},
		{
			method:  http.MethodPost,
			path:    "/diagnostic_centres/manager",
			handler: handler.CreateDiagnosticCentreManager,
			factory: func() interface{} {
				return &domain.DiagnosticCentreManagerRegisterDTO{}
			},
			description: "Create a diagnostic centre manager",
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_centres",
			handler: handler.SearchDiagnosticCentre,
			factory: func() interface{} {
				return &domain.SearchDiagnosticCentreQueryDTO{}
			},
			description: "List diagnostic centres",
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_centres/:diagnostic_centre_id",
			handler: handler.GetDiagnosticCentre,
			factory: func() interface{} {
				return &domain.GetDiagnosticParamDTO{}
			},
			description: "Get diagnostic centre by diagnostic_centre_id",
		},
		{
			method:      http.MethodPut,
			path:        "/diagnostic_centres/:diagnostic_centre_id",
			handler:     handler.UpdateDiagnosticCentre,
			factory:     func() interface{} { return &domain.UpdateDiagnosticBodyDTO{} },
			description: "Update diagnostic centre",
		},
		{
			method:      http.MethodDelete,
			path:        "/diagnostic_centres/:diagnostic_centre_id",
			handler:     handler.DeleteDiagnosticCentre,
			description: "Delete diagnostic centre",
		},
		{
			method:      http.MethodGet,
			path:        "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules",
			handler:     handler.GetDiagnosticCentreSchedules,
			factory:     func() interface{} { return &domain.GetDiagnosticSchedulesByCentreParamDTO{} },
			description: "Get Diagnostic centre schedules",
		},
		{
			method:      http.MethodPost,
			path:        "/diagnostic_centres_owner/kyc",
			handler:     handler.SubmitKYC,
			factory:     func() interface{} { return &domain.UpgradeDTO{} },
			description: "Submit diagnostic owner kyc",
		},
	}

	registerRoutes(group, diagnosticGroup)
}
