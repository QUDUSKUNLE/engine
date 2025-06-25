package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) CreatePayment(c echo.Context) error {
	// TODO: Implement payment creation
	return h.service.CreatePayment(c)
}

func (h *HTTPHandler) GetPayment(c echo.Context) error {
	// TODO: Implement get payment
	return h.service.GetPayment(c)
}

func (h *HTTPHandler) ListPayments(c echo.Context) error {
	// TODO: Implement list payments
	return h.service.ListPayments(c)
}

func (h *HTTPHandler) RefundPayment(c echo.Context) error {
	// TODO: Implement refund payment
	return h.service.RPayment(c)
}

func (h *HTTPHandler) PaymentWebhook(c echo.Context) error {
	// TODO: Implement payment webhook handler
	return h.service.HandlePaymentWebhook(c)
}

func (h *HTTPHandler) VerifyPayment(c echo.Context) error {
	reference := c.Param("reference")
	if reference == "" {
		return c.JSON(400, map[string]string{
			"error": "payment reference is required",
		})
	}

	return h.service.VerifyPayment(c, reference)
}
