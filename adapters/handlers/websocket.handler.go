package handlers

import (
	"net/http"

	"github.com/diagnoxix/core/utils"
	"github.com/labstack/echo/v4"
)

// WebSocketHandler handles WebSocket connection
// @Summary Establish WebSocket connection for real-time notifications
// @Description Upgrades HTTP connection to WebSocket for real-time notification delivery
// @Tags WebSocket
// @Param user_id query string true "User ID for the WebSocket connection"
// @Success 101 {string} string "Switching Protocols - WebSocket connection established"
// @Failure 400 {object} map[string]interface{} "Bad request - missing user_id"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ws/notifications [get]
func (h *HTTPHandler) WebSocketHandler(c echo.Context) error {
	return h.service.WebSocketManager.HandleWebSocket(c)
}

// GetWebSocketStatsHandler returns WebSocket connection statistics
// @Summary Get WebSocket connection statistics
// @Description Returns information about current WebSocket connections
// @Tags WebSocket
// @Produce json
// @Success 200 {object} map[string]interface{} "WebSocket statistics"
// @Router /vs [get]
func (h *HTTPHandler) GetWebSocketStatsHandler(c echo.Context) error {
	stats := map[string]interface{}{
		"connected_clients": h.service.WebSocketManager.GetConnectedClients(),
		"status":           "active",
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// SendTestNotificationHandler sends a test notification to a user
// @Summary Send test notification
// @Description Sends a test notification to a specific user via WebSocket
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param notification body TestNotificationRequest true "Test notification data"
// @Success 200 {object} map[string]interface{} "Notification sent successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /v1/ws/test-notification [post]
func (h *HTTPHandler) SendTestNotificationHandler(c echo.Context) error {
	var req TestNotificationRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind test notification request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Test notification validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Send test notification
	h.service.WebSocketManager.SendCustomNotification(
		req.UserID,
		"test_notification",
		req.Title,
		req.Message,
		map[string]interface{}{
			"test": true,
			"data": req.Data,
		},
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Test notification sent",
		"user_id": req.UserID,
	})
}

// TestNotificationRequest represents a test notification request
type TestNotificationRequest struct {
	UserID  string                 `json:"user_id" validate:"required"`
	Title   string                 `json:"title" validate:"required"`
	Message string                 `json:"message" validate:"required"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
