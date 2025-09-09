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
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {array} domain.Notification
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/notifications [get]
func (h *HTTPHandler) GetNotifications(c echo.Context) error {
	return h.service.GetNotifications(c)
}

// MarkNotificationRead marks a notification as read
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags Notifications
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param Notification_ID path string true "Notification ID"
// @Success 200 {object} domain.Notification
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/notifications/{notification_id}/read [post]
func (h *HTTPHandler) MarkNotificationRead(c echo.Context) error {
	return h.service.MarkNotificationRead(c)
}
