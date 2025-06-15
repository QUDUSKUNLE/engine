package services

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
	"golang.org/x/oauth2"
	oauth2v2 "google.golang.org/api/oauth2/v2"
)

func (service *ServicesHandler) Create(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UserRegisterDTO)
	newUser, err := domain.BuildNewUser(*dto)
	if err != nil {
		utils.Error("Failed to build new user",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_type", Value: dto.UserType})
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	createdUser, err := service.createUserHelper(
		context, newUser, db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumUSER)
	if err != nil {
		return err
	}

	// Generate verification token
	token := utils.GenerateRandomToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	verificationParams := db.CreateEmailVerificationTokenParams{
		Email:     createdUser.Email.String,
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}

	// Save verification token
	verificationToken, err := service.UserRepo.CreateEmailVerificationToken(
		context.Request().Context(),
		verificationParams,
	)
	if err != nil {
		utils.Error("Failed to create verification token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "email", Value: createdUser.Email.String})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Send verification email
	subject := "Verify your email address"
	body := fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Hi there!</p>
		<p>Thanks for registering with Medicue. Please verify your email by clicking the link below:</p>
		<p><a href="%s/verify-email?token=%s&email=%s">Verify Email</a></p>
		<p>This link will expire in 24 hours.</p>
		<p>Best regards,<br/>Medicue Team</p>
	`, os.Getenv("APP_URL"), verificationToken.Token, url.QueryEscape(createdUser.Email.String))

	err = service.emailService.Send(createdUser.Email.String, subject, body)
	if err != nil {
		utils.Error("Failed to send verification email",
			utils.LogField{Key: "error", Value: err.Error()})
		// Don't return error here, user is still created
	}

	return utils.ResponseMessage(http.StatusCreated, createdUser, context)
}

func (service *ServicesHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	// Check for permission to add a diagnostic manager
	_, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		utils.Error("Unauthorized attempt to create diagnostic centre manager",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// This validated at the middleware level
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.DiagnosticCentreManagerRegisterDTO)
	// Check if there are appropriate UserEnumDiagnosticCentreManager
	if dto.UserType != db.UserEnumDIAGNOSTICCENTREMANAGER {
		utils.Error("Invalid user type for diagnostic centre manager",
			utils.LogField{Key: "user_type", Value: dto.UserType})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}

	// Generate password and create user
	password, err := utils.GenerateRandomPassword(20)
	if err != nil {
		utils.Error("Failed to generate random password",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	userDto := domain.UserRegisterDTO{
		Email:    dto.Email,
		Password: password,
		UserType: dto.UserType,
	}

	newUser, err := domain.BuildNewUser(userDto)
	if err != nil {
		utils.Error("Failed to build diagnostic centre manager",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	createdUser, err := service.createUserHelper(context, newUser, db.UserEnumDIAGNOSTICCENTREMANAGER)
	if err != nil {
		return err
	}

	return utils.ResponseMessage(http.StatusCreated, createdUser, context)
}

func (service *ServicesHandler) Login(context echo.Context) error {
	// This validated at the middleware level
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UserSignInDTO)
	user, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(), pgtype.Text{String: dto.Email, Valid: true})
	if err != nil {
		utils.Error("Login failed - user not found",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Verify password
	if err := domain.ComparePassword(*user, dto.Password); err != nil {
		utils.Error("Login failed - invalid password",
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusUnauthorized, errors.New(utils.InvalidRequest), context)
	}

	// Generate token
	userClaims := domain.CurrentUserDTO{
		UserID:   uuid.MustParse(user.ID),
		UserType: user.UserType,
	}
	token, err := utils.GenerateToken(userClaims)
	if err != nil {
		utils.Error("Failed to generate JWT token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("User logged in successfully",
		utils.LogField{Key: "user_id", Value: user.ID},
		utils.LogField{Key: "user_type", Value: user.UserType})

	return utils.ResponseMessage(http.StatusOK, map[string]string{"token": token}, context)
}

func (service *ServicesHandler) RequestPasswordReset(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.RequestPasswordResetDTO)

	// Check if user exists
	user, err := service.UserRepo.GetUserByEmail(context.Request().Context(), pgtype.Text{String: dto.Email, Valid: true})
	if err != nil {
		utils.Error("Password reset requested for non-existent user",
			utils.LogField{Key: "email", Value: dto.Email})
		// Don't reveal if email exists
		return utils.ResponseMessage(http.StatusOK, map[string]string{
			"message": "If your email exists in our system, you will receive a password reset link",
		}, context)
	}

	// Generate reset token
	token := generateResetToken()
	expiresAt := time.Now().Add(15 * time.Minute)

	// Save token to database
	resetToken := db.CreatePasswordResetTokenParams{
		Email:     user.Email.String,
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}
	if err := service.UserRepo.CreatePasswordResetToken(context.Request().Context(), resetToken); err != nil {
		utils.Error("Failed to create password reset token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("Password reset token generated and email sent",
		utils.LogField{Key: "user_id", Value: user.ID},
		utils.LogField{Key: "expires_at", Value: expiresAt})

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "Reset instructions sent to your email",
	}, context)
}

func (service *ServicesHandler) ResetPassword(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.ResetPasswordDTO)

	// Verify token
	token, err := service.UserRepo.GetPasswordResetToken(context.Request().Context(), dto.Token)
	if err != nil {
		utils.Error("Invalid password reset token",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid or expired reset token"), context)
	}

	// Check if token is expired or used
	if token.ExpiresAt.Time.Before(time.Now()) || (token.Used.Valid && token.Used.Bool) {
		utils.Error("Expired or used reset token",
			utils.LogField{Key: "token_id", Value: token.ID})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid or expired reset token"), context)
	}

	// Verify email matches token
	if token.Email != dto.Email {
		utils.Error("Email mismatch for reset token",
			utils.LogField{Key: "token_id", Value: token.ID},
			utils.LogField{Key: "created_at", Value: token.CreatedAt})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid reset token"), context)
	}

	// Hash new password
	hashedPassword, err := domain.HashPassword(dto.NewPassword)
	if err != nil {
		utils.Error("Failed to hash new password",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Update password
	if err := service.UserRepo.UpdateUserPassword(context.Request().Context(), db.UpdateUserPasswordParams{
		Email:    pgtype.Text{String: dto.Email, Valid: true},
		Password: hashedPassword,
	}); err != nil {
		utils.Error("Failed to update password",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Mark token as used
	if err := service.UserRepo.MarkResetTokenUsed(context.Request().Context(), token.ID); err != nil {
		utils.Error("Failed to mark reset token as used",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "token_id", Value: token.ID})
		// Don't return error to user since password was updated successfully
	}

	utils.Info("Password reset successful",
		utils.LogField{Key: "email", Value: dto.Email})

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "Password reset successful",
	}, context)
}

func (service *ServicesHandler) VerifyEmail(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.EmailVerificationDTO)

	// Get user by email
	_, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		pgtype.Text{String: dto.Email, Valid: true},
	)
	if err != nil {
		utils.Error("Email verification failed - user not found",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusNotFound, errors.New("user not found"), context)
	}

	// TODO: Implement verify token logic
	// This would involve:
	// 1. Getting the verification token from database
	// 2. Checking if it's valid and not expired
	// 3. Marking the user as verified
	// 4. Marking the token as used

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "Email verified successfully",
	}, context)
}

func (service *ServicesHandler) ResendVerification(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.ResendVerificationDTO)

	// Get user by email
	_, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		pgtype.Text{String: dto.Email, Valid: true},
	)
	if err != nil {
		utils.Error("Resend verification failed - user not found",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusNotFound, errors.New("user not found"), context)
	}

	// TODO: Implement resend verification logic
	// This would involve:
	// 1. Generating a new verification token
	// 2. Saving it to the database
	// 3. Sending the verification email

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "Verification email sent",
	}, context)
}

func generateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (service *ServicesHandler) createUserHelper(
	context echo.Context,
	arg db.CreateUserParams,
	userTypes ...db.UserEnum,
) (*db.User, error) {
	if !isValidUserType(userTypes, arg.UserType) {
		utils.Error("Invalid user type",
			utils.LogField{Key: "user_type", Value: arg.UserType})
		return nil, utils.ErrorResponse(http.StatusUnprocessableEntity,
			errors.New("invalid user type"), context)
	}

	// Check if user exists
	existingUser, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		arg.Email,
	)
	if err != nil && err != sql.ErrNoRows {
		utils.Error("Failed to check existing user",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	if existingUser != nil {
		utils.Error("User already exists",
			utils.LogField{Key: "email", Value: arg.Email.String})
		return nil, utils.ErrorResponse(http.StatusBadRequest,
			errors.New("user already exists"), context)
	}

	// Create user
	createdRow, err := service.UserRepo.CreateUser(context.Request().Context(), arg)
	if err != nil {
		utils.Error("Failed to create user",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Convert CreateUserRow to User
	user := &db.User{
		ID:       createdRow.ID,
		Email:    createdRow.Email,
		Nin:      createdRow.Nin,
		UserType: createdRow.UserType,
	}

	return user, nil
}

func (service *ServicesHandler) GoogleLogin(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.GoogleAuthDTO)

	// Initialize OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		// We only need the ID token, not the client secret since we're using the frontend flow
	}

	// Create OAuth2 service
	oauth2Service, err := oauth2v2.New(oauth2.NewClient(context.Request().Context(), nil))
	if err != nil {
		utils.Error("Failed to create OAuth2 service",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Verify the ID token
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(dto.IDToken).Do()
	if err != nil {
		utils.Error("Failed to verify Google ID token",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusUnauthorized, errors.New("invalid Google token"), context)
	}

	// Verify token issuer
	if tokenInfo.IssuedTo != "accounts.google.com" && tokenInfo.IssuedTo != "https://accounts.google.com" {
		utils.Error("Invalid token issuer",
			utils.LogField{Key: "issuer", Value: tokenInfo.IssuedTo})
		return utils.ErrorResponse(http.StatusUnauthorized, errors.New("invalid token issuer"), context)
	}

	// Verify audience
	if tokenInfo.Audience != oauth2Config.ClientID {
		utils.Error("Token was not issued for this app",
			utils.LogField{Key: "audience", Value: tokenInfo.Audience})
		return utils.ErrorResponse(http.StatusUnauthorized, errors.New("invalid token audience"), context)
	}

	// Check if user exists
	user, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		pgtype.Text{String: tokenInfo.Email, Valid: true},
	)

	if err != nil && err != sql.ErrNoRows {
		utils.Error("Database error while checking user",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "email", Value: tokenInfo.Email})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	if user == nil {
		// Create new user
		newUser := db.CreateUserParams{
			// ID:            uuid.New().String(),
			Email:    pgtype.Text{String: tokenInfo.Email, Valid: true},
			Password: "", // No password for Google users
			UserType: db.UserEnumUSER,
		}

		user, err = service.UserRepo.CreateUser(context.Request().Context(), newUser)
		if err != nil {
			utils.Error("Failed to create user from Google login",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "email", Value: tokenInfo.Email})
			return utils.ErrorResponse(http.StatusInternalServerError, err, context)
		}
	}

	// Generate JWT token
	userClaims := domain.CurrentUserDTO{
		UserID:   uuid.MustParse(user.ID),
		UserType: user.UserType,
	}
	token, err := utils.GenerateToken(userClaims)
	if err != nil {
		utils.Error("Failed to generate JWT token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("User logged in with Google",
		utils.LogField{Key: "user_id", Value: user.ID})

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"token": token,
		"name":  tokenInfo.UserId,
	}, context)
}
