package handlers

import (
	"github.com/labstack/echo/v4"
)

// GetNotifications retrieves user notifications
// @Summary Get notifications
// @Description Get user's notifications with pagination
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {array} domain.Notification
// @Router /v1/notifications [get]
func (h *HTTPHandler) GetNotifications(c echo.Context) error {
	return h.service.GetNotifications(c)
}

// MarkNotificationRead marks a notification as read
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags Notifications
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} domain.Notification
// @Router /v1/notifications/{notification_id}/read [post]
func (h *HTTPHandler) MarkNotificationRead(c echo.Context) error {
	return h.service.MarkNotificationRead(c)
}
