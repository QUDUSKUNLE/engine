package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/adapters/middlewares"
	"github.com/medicue/core/domain"
)

type routeConfig struct {
	method  string
	path    string
	handler echo.HandlerFunc
	factory func() interface{}
}

func RoutesAdaptor(public *echo.Group, handler *handlers.HTTPHandler) *echo.Group {
	routes := []routeConfig{
		// Auth Routes
		{
			method:  http.MethodPost,
			path:    "/register",
			handler: handler.Register,
			factory: func() interface{} { return &domain.UserRegisterDTO{} },
		},
		{
			method:  http.MethodPost,
			path:    "/login",
			handler: handler.SignIn,
			factory: func() interface{} { return &domain.UserSignInDTO{} },
		},
		// Diagnostic centre routes
		{
			method:  http.MethodGet,
			path:    "/diagnostic_centres/:diagnostic_centre_id",
			handler: handler.GetDiagnosticCentre,
			factory: func() interface{} { return &domain.GetDiagnosticParamDTO{} },
		},
		{
			method:  http.MethodPost,
			path:    "/diagnostic_centre_manager",
			handler: handler.CreateDiagnosticCentreManager,
			factory: func() interface{} { return &domain.DiagnosticCentreManagerRegisterDTO{} },
		},
		{
			method:  http.MethodPost,
			path:    "/diagnostic_centres",
			handler: handler.CreateDiagnostic,
			factory: func() interface{} { return &domain.CreateDiagnosticDTO{} },
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_centres",
			handler: handler.SearchDiagnosticCentre,
			factory: func() interface{} { return &domain.SearchDiagnosticCentreQueryDTO{} },
		},
		{
			method:  http.MethodPut,
			path:    "/diagnostic_centres/:diagnostic_centre_id",
			handler: handler.UpdateDiagnosticCentre,
			factory: func() interface{} { return &domain.UpdateDiagnosticBodyDTO{} },
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules",
			handler: handler.GetDiagnosticSchedulesByCentre,
			factory: func() interface{} {
				return &domain.GetDiagnosticSchedulesByCentreParamDTO{}

			},
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules/:schedule_id",
			handler: handler.GetDiagnosticScheduleByCentre,
			factory: func() interface{} {
				return &domain.GetDiagnosticScheduleByCentreParamDTO{}

			},
		},
		{
			method:  http.MethodPut,
			path:    "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules/:schedule_id",
			handler: handler.UpdateDiagnosticScheduleByCentre,
			factory: func() interface{} { return &domain.UpdateDiagnosticScheduleByCentreDTO{} },
		},
		// Schedule Routes
		{
			method:  http.MethodPost,
			path:    "/diagnostic_schedules",
			handler: handler.CreateSchedule,
			factory: func() interface{} { return &domain.CreateScheduleDTO{} },
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_schedules/:schedule_id",
			handler: handler.GetSchedule,
			factory: func() interface{} { return &domain.GetDiagnosticScheduleParamDTO{} },
		},
		{
			method:  http.MethodPut,
			path:    "/diagnostic_schedules/:schedule_id",
			handler: handler.UpdateSchedule,
			factory: func() interface{} { return &domain.UpdateScheduleDTO{} },
		},
		{
			method:  http.MethodGet,
			path:    "/diagnostic_schedules",
			handler: handler.GetSchedules,
			factory: func() interface{} { return &domain.GetDiagnosticSchedulesQueryDTO{} },
		},
	}

	for _, r := range routes {
		switch r.method {
		case http.MethodPost:
			if r.factory != nil {
				public.POST(
					r.path,
					r.handler,
					middlewares.BodyValidationInterceptorFor(r.factory),
				)
			} else {
				public.POST(r.path, r.handler)
			}
		case http.MethodGet:
			if r.factory != nil {
				public.GET(
					r.path,
					r.handler,
					middlewares.BodyValidationInterceptorFor(r.factory),
				)
			} else {
				public.GET(r.path, r.handler)
			}
		case http.MethodPut:
			if r.factory != nil {
				public.PUT(
					r.path,
					r.handler,
					middlewares.BodyValidationInterceptorFor(r.factory),
				)
			} else {
				public.PUT(r.path, r.handler)
			}
			// Add more HTTP methods as needed
		}
	}
	return public
}
