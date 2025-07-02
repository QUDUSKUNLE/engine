package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medivue/adapters/handlers"
	"github.com/medivue/core/domain"
)

// AuthRoutes registers all authentication-related routes
func AuthRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	authGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/register",
			handler:     handler.Register,
			factory:     func() interface{} { return &domain.UserRegisterDTO{} },
			description: "User Registration",
		},
		{
			method:      http.MethodPost,
			path:        "/login",
			handler:     handler.SignIn,
			factory:     func() interface{} { return &domain.UserSignInDTO{} },
			description: "User Login",
		},
		{
			method:      http.MethodPost,
			path:        "/request_password_reset",
			handler:     handler.RequestPasswordReset,
			factory:     func() interface{} { return &domain.RequestPasswordResetDTO{} },
			description: "Request Password Reset",
		},
		{
			method:      http.MethodPost,
			path:        "/reset_password",
			handler:     handler.ResetPassword,
			factory:     func() interface{} { return &domain.ResetPasswordDTO{} },
			description: "Reset Password",
		},
		{
			method:  http.MethodGet,
			path:    "/verify_email",
			handler: handler.VerifyEmail,
			factory: func() interface{} {
				return &domain.EmailVerificationDTO{}
			},
			description: "Verify Email",
		},
		{
			method:  http.MethodPut,
			path:    "/account/profile",
			handler: handler.UpdateProfile,
			factory: func() interface{} {
				return &domain.UpdateUserProfileDTO{}
			},
			description: "Update User Profile",
		},
		{
			method:  http.MethodGet,
			path:    "/account/profile",
			handler: handler.GetProfile,
			factory: func() interface{} {
				return &domain.GetProfileDTO{}
			},
			description: "Get A User Profile",
		},
		{
			method:  http.MethodPost,
			path:    "/diagnostic_centre_manager",
			handler: handler.CreateDiagnosticCentreManager,
			factory: func() interface{} {
				return &domain.DiagnosticCentreManagerRegisterDTO{}
			},
			description: "Create A Diagnostic Centre Manager",
		},
	}

	registerRoutes(group, authGroup)
}
