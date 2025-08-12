package handlers

import (
	"github.com/labstack/echo/v4"
)

// UserSwagger is used for Swagger documentation only
// @Description User response for Swagger
// @name UserSwagger
type UserSwagger struct {
	ID        string `json:"id" example:"user-001"`
	Email     string `json:"email" example:"user@example.com"`
	FullName  string `json:"full_name" example:"John Doe"`
	Role      string `json:"role" example:"PATIENT"`
	CreatedAt string `json:"created_at" example:"2025-06-27T12:00:00Z"` // format: date-time
	UpdatedAt string `json:"updated_at" example:"2025-06-27T12:00:00Z"` // format: date-time
	// ...add other fields as needed for docs
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user (patient) or diagnostic centre owner in the system
// @Tags User
// @Accept json
// @Produce json
// @Param user body domain.UserRegisterDTO true "User registration details"
// @Success 201 {object} handlers.UserSwagger "User created successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input data/Email already exists"
// @Failure 422 {object} handlers.ErrorResponse "Invalid user type/Validation failed"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/register [post]
func (handler *HTTPHandler) Register(context echo.Context) error {
	return handler.service.Create(context)
}

// SignIn godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token for API access
// @Tags User
// @Accept json
// @Produce json
// @Param credentials body domain.UserSignInDTO true "User credentials"
// @Success 200 {object} map[string]string "token: JWT token for authentication"
// @Failure 400 {object} handlers.ErrorResponse "Invalid credentials"
// @Failure 401 {object} handlers.ErrorResponse "Invalid email or password"
// @Failure 404 {object} handlers.ErrorResponse "User not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/login [post]
func (handler *HTTPHandler) SignIn(context echo.Context) error {
	return handler.service.Login(context)
}

// RequestPasswordReset godoc
// @Summary Request password reset
// @Description Send a password reset link to user's email
// @Tags User
// @Accept json
// @Produce json
// @Param request body domain.RequestPasswordResetDTO true "Password reset request"
// @Success 200 {object} map[string]string "message: Reset link sent if email exists"
// @Failure 400 {object} handlers.ErrorResponse "Invalid email format"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/request_password_reset [post]
func (handler *HTTPHandler) RequestPasswordReset(context echo.Context) error {
	return handler.service.RequestPasswordReset(context)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user's password using the token received via email
// @Tags User
// @Accept json
// @Produce json
// @Param reset body domain.ResetPasswordDTO true "Password reset details"
// @Success 200 {object} map[string]string "message: Password reset successful"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input/Expired token"
// @Failure 422 {object} handlers.ErrorResponse "Password validation failed"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/reset_password [post]
func (handler *HTTPHandler) ResetPassword(context echo.Context) error {
	return handler.service.ResetPassword(context)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} handlers.UserSwagger "User profile details"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 404 {object} handlers.ErrorResponse "User not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/account/profile [get]
func (handler *HTTPHandler) GetProfile(context echo.Context) error {
	return handler.service.GetProfile(context)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param profile body domain.UpdateUserProfileDTO true "Profile update details"
// @Success 200 {object} handlers.UserSwagger "Updated profile details"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input data"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/account/profile [put]
func (handler *HTTPHandler) UpdateProfile(context echo.Context) error {
	return handler.service.UpdateProfile(context)
}

// DeactivateAccount godoc
// @Summary Deactivate user account
// @Description Deactivate user's account (soft delete)
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param deactivate body domain.DeactivateAccountDTO true "Account deactivation details"
// @Success 200 {object} map[string]string "message: Account deactivated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input data"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required/Invalid password"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/account [delete]
func (handler *HTTPHandler) DeactivateAccount(context echo.Context) error {
	// return handler.service.DeactivateAccount(context)
	return nil
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify a user's email address using the token sent to their email
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "message: Email verified successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input data"
// @Failure 401 {object} handlers.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/verify_email [get]
func (handler *HTTPHandler) VerifyEmail(context echo.Context) error {
	return handler.service.VerifyEmail(context)
}

// ResendVerification godoc
// @Summary Resend verification email
// @Description Resend the email verification token to the user's email address
// @Tags User
// @Accept json
// @Produce json
// @Param email body domain.ResendVerificationDTO true "Email address"
// @Success 200 {object} map[string]string "message: Verification email sent"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input data"
// @Failure 404 {object} handlers.ErrorResponse "User not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/resend_verification [post]
func (handler *HTTPHandler) ResendVerification(context echo.Context) error {
	return handler.service.ResendVerification(context)
}


func (handler *HTTPHandler) GoogleLogin(context echo.Context) error {
	return handler.service.GoogleLogin(context)
}
