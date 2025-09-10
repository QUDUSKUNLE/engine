package handlers

import (
	"net/http"

	"github.com/diagnoxix/core/utils"
	"github.com/labstack/echo/v4"
)

// GetCacheStatsHandler returns cache performance statistics
// @Summary Get cache statistics
// @Description Returns cache performance metrics including hit rate, response times, and memory usage
// @Tags Cache
// @Produce json
// @Success 200 {object} map[string]interface{} "Cache statistics"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/cache/stats [get]
func (h *HTTPHandler) GetCacheStatsHandler(c echo.Context) error {
	stats, err := h.service.AI.GetCacheStats(c.Request().Context())
	if err != nil {
		utils.Error("Failed to get cache stats", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve cache statistics",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// ClearCacheHandler clears all cached AI responses
// @Summary Clear AI cache
// @Description Removes all cached AI responses to force fresh API calls
// @Tags Cache
// @Produce json
// @Success 200 {object} map[string]interface{} "Cache cleared successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/cache/clear [post]
func (h *HTTPHandler) ClearCacheHandler(c echo.Context) error {
	err := h.service.AI.ClearCache(c.Request().Context())
	if err != nil {
		utils.Error("Failed to clear cache", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to clear cache",
		})
	}

	utils.Info("AI cache cleared by user", utils.LogField{Key: "user_id", Value: "admin"}) // You can extract user ID from JWT

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cache cleared successfully",
	})
}

// InvalidateCacheKeyHandler removes a specific cache entry
// @Summary Invalidate specific cache entry
// @Description Removes a specific cached AI response by operation and input hash
// @Tags Cache
// @Accept json
// @Produce json
// @Param invalidation body CacheInvalidationRequest true "Cache invalidation request"
// @Success 200 {object} map[string]interface{} "Cache entry invalidated"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/cache/invalidate [post]
func (h *HTTPHandler) InvalidateCacheKeyHandler(c echo.Context) error {
	var req CacheInvalidationRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind cache invalidation request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Cache invalidation validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	err := h.service.AI.InvalidateCacheKey(c.Request().Context(), req.Operation, req.Input)
	if err != nil {
		utils.Error("Failed to invalidate cache key", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to invalidate cache entry",
		})
	}

	utils.Info("Cache key invalidated", 
		utils.LogField{Key: "operation", Value: req.Operation},
		utils.LogField{Key: "user_id", Value: "admin"}) // Extract from JWT

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cache entry invalidated successfully",
	})
}

// GetCacheHealthHandler returns cache health status
// @Summary Get cache health
// @Description Returns the health status of the caching system
// @Tags Cache
// @Produce json
// @Success 200 {object} map[string]interface{} "Cache health status"
// @Router /v1/cache/health [get]
func (h *HTTPHandler) GetCacheHealthHandler(c echo.Context) error {
	// Try to get cache stats to verify cache is working
	_, err := h.service.AI.GetCacheStats(c.Request().Context())
	
	status := "healthy"
	message := "Cache is operational"
	
	if err != nil {
		status = "degraded"
		message = "Cache is not available, using fallback"
		utils.Warn("Cache health check failed", utils.LogField{Key: "error", Value: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  status,
		"message": message,
		"cache_available": err == nil,
	})
}

// CacheInvalidationRequest represents a cache invalidation request
type CacheInvalidationRequest struct {
	Operation string      `json:"operation" validate:"required,oneof=lab_interpretation symptom_analysis report_summary image_analysis anomaly_detection lab_package_analysis automated_report"`
	Input     interface{} `json:"input" validate:"required"`
}
