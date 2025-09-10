package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// CacheRoutes registers all cache-related routes
func CacheRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	cacheGroup := []routeConfig{
		{
			method:      http.MethodGet,
			path:        "/cache/stats",
			handler:     handler.GetCacheStatsHandler,
			factory:     func() interface{} { return &domain.CapabilitiesDTO{} },
			description: "Get cache performance statistics",
		},
		{
			method:      http.MethodPost,
			path:        "/cache/clear",
			handler:     handler.ClearCacheHandler,
			factory:     func() interface{} { return &domain.CapabilitiesDTO{} },
			description: "Clear all cached AI responses",
		},
		{
			method:      http.MethodPost,
			path:        "/cache/invalidate",
			handler:     handler.InvalidateCacheKeyHandler,
			factory:     func() interface{} { return &handlers.CacheInvalidationRequest{} },
			description: "Invalidate specific cache entry",
		},
		{
			method:      http.MethodGet,
			path:        "/cache/health",
			handler:     handler.GetCacheHealthHandler,
			factory:     func() interface{} { return &domain.CapabilitiesDTO{} },
			description: "Get cache health status",
		},
	}

	registerRoutes(group, cacheGroup)
}
