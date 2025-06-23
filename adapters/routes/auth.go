package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/core/domain"
)

// AuthRoutes registers all authentication-related routes
func AuthRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	authGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/register",
			handler:     handler.Register,
			factory:     func() interface{} { return &domain.UserRegisterDTO{} },
			description: "User registration",
		},
		{
			method:      http.MethodPost,
			path:        "/login",
			handler:     handler.SignIn,
			factory:     func() interface{} { return &domain.UserSignInDTO{} },
			description: "User login",
		},
		{
			method:      http.MethodPost,
			path:        "/request_password_reset",
			handler:     handler.RequestPasswordReset,
			factory:     func() interface{} { return &domain.RequestPasswordResetDTO{} },
			description: "Request password reset",
		},
		{
			method:      http.MethodPost,
			path:        "/reset_password",
			handler:     handler.ResetPassword,
			factory:     func() interface{} { return &domain.ResetPasswordDTO{} },
			description: "Reset password",
		},
		{
			method:      http.MethodGet,
			path:        "/verify_email",
			handler:     handler.VerifyEmail,
			description: "Verify email",
		},
		{
			method:  http.MethodPost,
			path:    "/diagnostic_centre_manager",
			handler: handler.CreateDiagnosticCentreManager,
			factory: func() interface{} {
				return &domain.DiagnosticCentreManagerRegisterDTO{}
			},
			description: "Create a diagnostic centre manager",
		},
	}

	registerRoutes(group, authGroup)
}
