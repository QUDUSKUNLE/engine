package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/core/domain"
)

// PaymentRoutes registers all payment-related routes
func PaymentRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	paymentGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/payments",
			handler:     handler.CreatePayment,
			factory:     func() interface{} { return &domain.CreatePaymentDTO{} },
			description: "Process a new payment",
		},
		{
			method:      http.MethodGet,
			path:        "/payments/:payment_id",
			handler:     handler.GetPayment,
			factory:     func() interface{} { return &domain.GetPaymentDTO{} },
			description: "Get payment details",
		},
		{
			method:      http.MethodGet,
			path:        "/payments",
			handler:     handler.ListPayments,
			factory:     func() interface{} { return &domain.ListPaymentsDTO{} },
			description: "List payments with filters",
		},
		{
			method:      http.MethodPost,
			path:        "/payments/:payment_id/refund",
			handler:     handler.RefundPayment,
			factory:     func() interface{} { return &domain.RefundPaymentDTO{} },
			description: "Process payment refund",
		},
		{
			method:      http.MethodPost,
			path:        "/payments/webhook",
			handler:     handler.PaymentWebhook,
			factory:     func() interface{} { return &domain.PaymentWebhookDTO{} },
			description: "Handle payment provider webhook",
		},
	}

	registerRoutes(group, paymentGroup)
}
