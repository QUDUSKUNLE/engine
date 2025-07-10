package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medivue/adapters/handlers"
	"github.com/medivue/core/domain"
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
