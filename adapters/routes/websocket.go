package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// WebSocketRoutes registers all WebSocket-related routes
func WebSocketRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	wsGroup := []routeConfig{
		{
			method:      http.MethodGet,
			path:        "/notifications",
			handler:     handler.WebSocketHandler,
			factory:     func() interface{} { return &domain.CapabilitiesDTO{} },
			description: "Establish WebSocket connection for real-time notifications",
		},
		{
			method:      http.MethodGet,
			path:        "/stats",
			handler:     handler.GetWebSocketStatsHandler,
			factory:     func() interface{} { return &domain.CapabilitiesDTO{} },
			description: "Get WebSocket connection statistics",
		},
		{
			method:      http.MethodPost,
			path:        "/test-notification",
			handler:     handler.SendTestNotificationHandler,
			factory:     func() interface{} { return &handlers.TestNotificationRequest{} },
			description: "Send test notification via WebSocket",
		},
	}

	registerRoutes(group, wsGroup)
}
