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
	}

	for _, r := range routes {
		switch r.method {
		case http.MethodPost:
			public.POST(
				r.path,
				r.handler,
				middlewares.BodyValidationInterceptorFor(r.factory),
			)
		case http.MethodGet:
			public.GET(
				r.path,
				r.handler,
				middlewares.BodyValidationInterceptorFor(r.factory),
			)
		case http.MethodPut:
			public.PUT(
				r.path,
				r.handler,
				middlewares.BodyValidationInterceptorFor(r.factory),
			)
			// Add more HTTP methods as needed
		}
	}
	return public
}
