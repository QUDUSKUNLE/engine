package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/middlewares"
	"github.com/medicue/core/utils"
)

// routeConfig defines the structure for route configuration
type routeConfig struct {
	method      string
	path        string
	handler     echo.HandlerFunc
	factory     func() interface{}
	description string // Added for better logging and documentation
}


// registerRoutes registers the given routes with the specified group
func registerRoutes(group *echo.Group, routes []routeConfig) {
	for _, route := range routes {
		utils.Info("Registering route",
			utils.LogField{Key: "method", Value: route.method},
			utils.LogField{Key: "path", Value: "/v1" + route.path},
			utils.LogField{Key: "description", Value: route.description})

		switch route.method {
		case http.MethodPost:
			group.POST(
				route.path,
				route.handler,
				middlewares.BodyValidationInterceptorFor(route.factory),
			)
		case http.MethodGet:
			group.GET(
				route.path,
				route.handler,
				middlewares.BodyValidationInterceptorFor(route.factory),
			)
		case http.MethodPut:
			group.PUT(
				route.path,
				route.handler,
				middlewares.BodyValidationInterceptorFor(route.factory),
			)
		case http.MethodDelete:
			group.DELETE(
				route.path,
				route.handler,
				middlewares.BodyValidationInterceptorFor(route.factory),
			)
		default:
			utils.Error("Unsupported HTTP method",
				utils.LogField{Key: "method", Value: route.method},
				utils.LogField{Key: "path", Value: route.path})
		}
	}
}
