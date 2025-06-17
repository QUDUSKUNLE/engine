package handlers

import (
	"github.com/labstack/echo/v4"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user (patient) or diagnostic centre owner in the system
// @Tags User
// @Accept json
// @Produce json
// @Param user body domain.UserRegisterDTO true "User registration details"
// @Success 201 {object} db.User "User created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data/Email already exists"
// @Failure 422 {object} utils.ErrorResponse "Invalid user type/Validation failed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
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
// @Failure 400 {object} utils.ErrorResponse "Invalid credentials"
// @Failure 401 {object} utils.ErrorResponse "Invalid email or password"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/login [post]
func (handler *HTTPHandler) SignIn(context echo.Context) error {
	return handler.service.Login(context)
}

// CreateDiagnosticCentreManager godoc
// @Summary Create a diagnostic centre manager
// @Description Create a new diagnostic centre manager account. Only accessible by diagnostic centre owners.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param manager body domain.DiagnosticCentreManagerRegisterDTO true "Manager details"
// @Success 201 {object} db.User "Manager account created successfully"
// @Success 202 {object} map[string]string "Manager invite sent successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data/Email already exists"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Not a diagnostic centre owner"
// @Failure 422 {object} utils.ErrorResponse "Invalid manager type"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/diagnostic_centre_manager [post]
func (handler *HTTPHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	return handler.service.CreateDiagnosticCentreManager(context)
}

// RequestPasswordReset godoc
// @Summary Request password reset
// @Description Send a password reset link to user's email
// @Tags User
// @Accept json
// @Produce json
// @Param request body domain.RequestPasswordResetDTO true "Password reset request"
// @Success 200 {object} map[string]string "message: Reset link sent if email exists"
// @Failure 400 {object} utils.ErrorResponse "Invalid email format"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/request-password-reset [post]
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
// @Failure 400 {object} utils.ErrorResponse "Invalid input/Expired token"
// @Failure 422 {object} utils.ErrorResponse "Password validation failed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/reset-password [post]
func (handler *HTTPHandler) ResetPassword(context echo.Context) error {
	return handler.service.ResetPassword(context)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Allow user to change their current password
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param password body domain.ChangePasswordDTO true "Password change details"
// @Success 200 {object} map[string]string "message: Password changed successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required/Invalid current password"
// @Failure 422 {object} utils.ErrorResponse "Password validation failed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/account/password [put]
func (handler *HTTPHandler) ChangePassword(context echo.Context) error {
	// return handler.service.ChangePassword(context)
	return nil
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} domain.UserProfileResponse "User profile details"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
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
// @Success 200 {object} domain.UserProfileResponse "Updated profile details"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
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
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required/Invalid password"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
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
// @Param verification body domain.EmailVerificationDTO true "Email verification details"
// @Success 200 {object} map[string]string "message: Email verified successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/verify-email [post]
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
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/resend-verification [post]
func (handler *HTTPHandler) ResendVerification(context echo.Context) error {
	return handler.service.ResendVerification(context)
}

// GoogleLogin godoc
// @Summary Login with Google
// @Description Authenticate a user using Google OAuth
// @Tags User
// @Accept json
// @Produce json
// @Param body body domain.GoogleAuthDTO true "Google ID token"
// @Success 200 {object} map[string]string "token: JWT token for authentication"
// @Failure 400 {object} utils.ErrorResponse "Invalid token"
// @Failure 401 {object} utils.ErrorResponse "Authentication failed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/auth/google [post]
func (handler *HTTPHandler) GoogleLogin(context echo.Context) error {
	return handler.service.GoogleLogin(context)
}
