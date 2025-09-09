package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// AuthRoutes registers all authentication-related routes
func AuthRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	authGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/register",
			handler:     handler.Register,
			factory:     func() interface{} { return &domain.RegisterationDTO{} },
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
			method:      http.MethodPut,
			path:        "/update_password",
			handler:     handler.UpdatePassword,
			factory:     func() interface{} { return &domain.UpdatePasswordDTO{} },
			description: "Update password",
		},
		{
			method:      http.MethodPost,
			path:        "/reset_password",
			handler:     handler.ResetPassword,
			factory:     func() interface{} { return &domain.ResetPasswordDTO{} },
			description: "Reset password",
		},
		{
			method:      http.MethodPost,
			path:        "/resend_verification",
			handler:     handler.ResendVerification,
			factory:     func() interface{} { return &domain.ResendVerificationDTO{} },
			description: "Resend verification",
		},
		{
			method:  http.MethodGet,
			path:    "/verify_email",
			handler: handler.VerifyEmail,
			factory: func() interface{} {
				return &domain.EmailVerificationDTO{}
			},
			description: "Verify email",
		},
		{
			method:  http.MethodPut,
			path:    "/account/profile",
			handler: handler.UpdateProfile,
			factory: func() interface{} {
				return &domain.UpdateUserProfileDTO{}
			},
			description: "Update user profile",
		},
		{
			method:  http.MethodGet,
			path:    "/account/profile",
			handler: handler.GetProfile,
			factory: func() interface{} {
				return &domain.GetProfileDTO{}
			},
			description: "Get a user profile",
		},
		{
			method:  http.MethodGet,
			path:    "/managers",
			handler: handler.ListManagersByAdmin,
			factory: func() interface{} {
				return &domain.GetManagerDTO{}
			},
			description: "List managers",
		},
	}

	registerRoutes(group, authGroup)
}
