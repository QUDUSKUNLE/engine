package handlers

import (
	"github.com/labstack/echo/v4"
)

// GetPayment retrieves details of a specific payment
// @Summary Get payment
// @Description Get information about a single payment by its ID
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param payment_id path string true "Payment ID"
// @Success 200 {object} domain.GetPaymentDTO
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/payments/{payment_id} [get]
func (h *HTTPHandler) GetPayment(c echo.Context) error {
	return h.service.GetPayment(c)
}

// ListPayments retrieves a paginated list of payments
// @Summary List payments
// @Description Get a paginated list of all payments
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id query uuid true "Filter by diagnostic centre ID"
// @Param patient_id query uuid true "Filter by patient ID"
// @Param status query string true "Filter by payment status" Enums(pending,success,failed,refunded,cancelled)
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {array} domain.ListPaymentsDTO
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/payments [get]
func (h *HTTPHandler) ListPayments(c echo.Context) error {
	return h.service.ListPayments(c)
}

// RefundPayment refunds a specific payment
// @Summary Refund payment
// @Description Issue a refund for a specific payment by its ID
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param payment_id path uuid true "Payment ID to refund"
// @Param refund_request body domain.RefundPaymentDTO true "Refund details (amount, reason, etc.)"
// @Success 200 {object} domain.GetPaymentDTO "Refunded payment details"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/payments/{payment_id}/refund [post]
func (h *HTTPHandler) RefundPayment(c echo.Context) error {
	return h.service.RPayment(c)
}

// PaymentWebhook handles incoming payment gateway webhooks
// @Summary Payment gateway webhook
// @Description Receives asynchronous payment status updates from the payment provider
// @Tags Payments
// @Accept json
// @Produce json
// @Param X-Signature header string false "Optional signature header for webhook verification"
// @Param webhook body domain.PaymentWebhookDTO true "Webhook payload from the payment gateway"
// @Success 200 {object} object "Acknowledgement of webhook received"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/payments/webhook [post]
func (h *HTTPHandler) PaymentWebhook(c echo.Context) error {
	return h.service.HandlePaymentWebhook(c)
}

// VerifyPayment verifies the status of a specific payment
// @Summary Verify payment
// @Description Check the current status of a payment using its ID or reference
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param reference path string true "Payment reference to verify"
// @Success 200 {object} domain.GetPaymentDTO "Verified payment details"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/payments/verify/{reference} [get]
func (h *HTTPHandler) VerifyPayment(c echo.Context) error {
	return h.service.VerifyPayment(c)
}
