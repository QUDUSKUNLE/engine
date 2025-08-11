package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// DiagnosticRoutes registers all diagnostic centre-related routes
func AIRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	aiGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/interpret",
			handler:     handler.InterpretLabHandler,
			factory:     func() interface{} { return &domain.LabTest{} },
			description: "Interpret test result",
		},
	}

	registerRoutes(group, aiGroup)
}
