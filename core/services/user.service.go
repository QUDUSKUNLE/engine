package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/adapters/ex/templates/emails"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
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
	newUser.EmailVerified = pgtype.Bool{Bool: false, Valid: true}
	createdUser, err := service.createUserHelper(
		context, newUser, db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumPATIENT)
	if err != nil {
		return utils.ErrorResponse(http.StatusConflict, err, context)
	}

	// Generate verification token
	token := GenerateRandomToken()
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

	// // Send verification email
	escapedEmail := url.QueryEscape(createdUser.Email.String)

	emaildata := &emails.EmailVerificationData{
		Name:             newUser.Fullname.String,
		VerificationLink: fmt.Sprintf("%s/v1/verify_email?token=%s&email=%s", service.Config.AppUrl, verificationToken.Token, escapedEmail),
		ExpiryDuration:   "24 hours",
	}

	go service.emailGoroutine(emaildata, createdUser.Email.String, emails.SubjectEmailVerification, emails.TemplateEmailVerification)

	return utils.ResponseMessage(http.StatusCreated, createdUser, context)
}

func (service *ServicesHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	// Check for permission to add a diagnostic manager
	owner, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
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
	password, err := GenerateRandomPassword(15)
	if err != nil {
		utils.Error("Failed to generate random password",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	userDto := domain.UserRegisterDTO{
		Email:        dto.Email,
		LastName:     dto.LastName,
		FirstName:    dto.FirstName,
		Password:     password,
		UserType:     dto.UserType,
		CreatedAdmin: owner.UserID,
	}

	newUser, err := domain.BuildNewUser(userDto)
	if err != nil {
		utils.Error("Failed to build diagnostic centre manager",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	newUser.EmailVerified = pgtype.Bool{Bool: true, Valid: true}
	createdUser, err := service.createUserHelper(context, newUser, db.UserEnumDIAGNOSTICCENTREMANAGER)
	if err != nil {
		return utils.ErrorResponse(http.StatusConflict, err, context)
	}
	// Send registration email
	emaildata := &emails.DiagnosticCentreManager{
		ManagerName: newUser.Fullname.String,
		Email:       newUser.Email.String,
		Password:    password,
	}

	go service.emailGoroutine(
		emaildata,
		createdUser.Email.String,
		emails.SubjectDiagnosticCentreManager,
		emails.TemplateDiagnosticCentreManager,
	)

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
		if err.Error() == "no rows in result set" {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidLoginRequest), context)
		}
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	// Verify password
	if err := domain.ComparePassword(*user, dto.Password); err != nil {
		utils.Error("Login failed - invalid password",
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidLoginRequest), context)
	}

	// Generate token
	userClaims := domain.CurrentUserDTO{
		UserID:   uuid.MustParse(user.ID),
		UserType: user.UserType,
	}
	token, err := GenerateToken(userClaims)
	if err != nil {
		utils.Error("Failed to generate JWT token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("User logged in successfully",
		utils.LogField{Key: "user_id", Value: user.ID},
		utils.LogField{Key: "user_type", Value: user.UserType})

	return utils.ResponseMessage(
		http.StatusOK,
		map[string]string{
			"token":     token,
			"user_type": string(user.UserType),
		}, context)
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
	expiresAt := time.Now().Add(10 * time.Minute)

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

	// Send password reset email

	emailData := &emails.PasswordResetData{
		Name:      user.Fullname.String,
		ResetLink: fmt.Sprintf("%s/v1/reset_password?token=%s&email=%s", service.Config.AppUrl, token, url.QueryEscape(user.Email.String)),
		ExpiresIn: "15 minutes",
	}

	go service.emailGoroutine(
		emailData,
		user.Email.String,
		emails.SubjectResetPassword,
		emails.TemplateResetPassword,
	)

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "If your email exists in our system, you will receive password reset instructions",
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
	dto, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.EmailVerificationDTO)

	// Get user by email
	user, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		pgtype.Text{String: dto.Email, Valid: true},
	)
	if err != nil {
		utils.Error("Email verification failed - user not found",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusNotFound, errors.New("user not found"), context)
	}

	// Get and verify token
	token, err := service.UserRepo.GetEmailVerificationToken(context.Request().Context(), dto.Token)
	if err != nil {
		utils.Error("Invalid verification token",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid verification token"), context)
	}

	// Check if token is expired or used
	if token.ExpiresAt.Time.Before(time.Now()) || (token.Used.Valid && token.Used.Bool) {
		utils.Error("Expired or used verification token",
			utils.LogField{Key: "token_id", Value: token.ID})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("verification token expired or already used"), context)
	}

	// Verify email matches token
	if token.Email != dto.Email {
		utils.Error("Email mismatch for verification token",
			utils.LogField{Key: "token_id", Value: token.ID})
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid verification token"), context)
	}

	// Marked user as verified
	err = service.UserRepo.MarkEmailAsVerified(context.Request().Context(), dto.Email)
	if err != nil {
		utils.Error("Failed to mark user as verified",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Mark token as used
	if err := service.UserRepo.MarkEmailVerificationTokenUsed(context.Request().Context(), token.ID); err != nil {
		utils.Error("Failed to mark verification token as used",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "token_id", Value: token.ID})
		// Don't return error since verification was successful
	}

	utils.Info("Email verified successfully",
		utils.LogField{Key: "user_id", Value: user.ID},
		utils.LogField{Key: "email", Value: user.Email.String})

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "Email verified successfully",
	}, context)
}

func (service *ServicesHandler) ResendVerification(context echo.Context) error {
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.ResendVerificationDTO)

	// Get user by email
	user, err := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		pgtype.Text{String: dto.Email, Valid: true},
	)
	if err != nil {
		utils.Error("Resend verification failed - user not found",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusNotFound, errors.New("user not found"), context)
	}

	// Generate new verification token
	token := GenerateRandomToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	verificationParams := db.CreateEmailVerificationTokenParams{
		Email:     user.Email.String,
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}

	// Save token to database
	_, err = service.UserRepo.CreateEmailVerificationToken(
		context.Request().Context(),
		verificationParams,
	)
	if err != nil {
		utils.Error("Failed to create verification token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "email", Value: user.Email.String})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Send verification email
	escapedEmail := url.QueryEscape(user.Email.String)
	emaildata := &emails.EmailVerificationData{
		Name:             user.Fullname.String,
		VerificationLink: fmt.Sprintf("%s/v1/verify_email?token=%s&email=%s", service.Config.AppUrl, token, escapedEmail),
		ExpiryDuration:   "24 hours",
	}

	go service.emailGoroutine(
		emaildata,
		user.Email.String,
		emails.SubjectEmailVerification,
		emails.TemplateEmailVerification,
	)

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"message": "Verification email sent",
	}, context)
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
			Fullname: pgtype.Text{String: tokenInfo.Email, Valid: true},
			// ID:            uuid.New().String(),
			Email:    pgtype.Text{String: tokenInfo.Email, Valid: true},
			Password: "", // No password for Google users
			UserType: db.UserEnumPATIENT,
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
	token, err := GenerateToken(userClaims)
	if err != nil {
		utils.Error("Failed to generate JWT token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("User logged in with Google",
		utils.LogField{Key: "user_id", Value: user.ID})

	return utils.ResponseMessage(http.StatusOK, map[string]string{
		"token":     token,
		"user_type": string(user.UserType),
	}, context)
}

func (service *ServicesHandler) UpdateProfile(context echo.Context) error {
	// Get validated DTO from context
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateUserProfileDTO)

	// Get current user from context
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumPATIENT, db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		utils.Error("Failed to get current user",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Update user profile
	updateParams := db.UpdateUserParams{
		ID:       currentUser.UserID.String(),
		Fullname: pgtype.Text{String: fmt.Sprintf("%s %s", dto.FirstName, dto.LastName), Valid: true},
		Nin:      pgtype.Text{String: dto.Nin, Valid: true},
	}
	if dto.PhoneNumber != "" {
		updateParams.PhoneNumber = pgtype.Text{String: dto.PhoneNumber, Valid: true}
	}

	updatedUser, err := service.UserRepo.UpdateUser(context.Request().Context(), updateParams)
	if err != nil {
		if err.Error() == "no rows in result set" {
			utils.Error("Failed to update user profile",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "user_id", Value: currentUser.UserID})
			return utils.ErrorResponse(http.StatusUnprocessableEntity, err, context)
		}
		utils.Error("Failed to update user profile",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	utils.Info("User profile updated successfully",
		utils.LogField{Key: "user_id", Value: updatedUser.ID})

	return utils.ResponseMessage(http.StatusOK, updatedUser, context)
}

func (service *ServicesHandler) GetProfile(context echo.Context) error {
	// Get current user from context
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumPATIENT, db.UserEnumDIAGNOSTICCENTREMANAGER, db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		utils.Error("Failed to get current user",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get user profile from database
	user, err := service.UserRepo.GetUser(context.Request().Context(), currentUser.UserID.String())
	if err != nil {
		if err.Error() == "no rows in result set" {
			utils.Error("Failed to get user profile",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "user_id", Value: currentUser.UserID})
			return utils.ErrorResponse(http.StatusNotFound, err, context)
		}
		utils.Error("Failed to get user profile",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	db_user := &db.User{
		ID:              user.ID,
		Email:           user.Email,
		Nin:             user.Nin,
		UserType:        user.UserType,
		Fullname:        user.Fullname,
		EmailVerified:   user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneNumber:     user.PhoneNumber,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}

	utils.Info("User profile retrieved successfully",
		utils.LogField{Key: "user_id", Value: db_user.ID})

	return utils.ResponseMessage(http.StatusOK, db_user, context)
}

func (service *ServicesHandler) ListManagers(context echo.Context) error {
	// Get validated DTO from context
	dto, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetManagerDTO)
	// Get current user from context
	admin, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		utils.Error("Failed to get current user",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	pagination := SetDefaultPagination(&domain.PaginationQueryDTO{
		Page:    dto.PaginationQueryDTO.Page,
		PerPage: dto.PaginationQueryDTO.PerPage,
	})
	response, err := service.UserRepo.ListManagersByadmin(context.Request().Context(), db.ListUsersByAdminParams{
		CreatedAdmin: pgtype.UUID{Bytes: admin.UserID, Valid: true},
		Limit:        pagination.GetLimit(),
		Offset:       pagination.GetOffset(),
		Column4:      !dto.Assigned,
	})
	if err != nil {
		utils.Error("Failed to list managers",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: admin.UserID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	if len(response) == 0 {
		response = []*db.ListUsersByAdminRow{}
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) OwnerKYC(context echo.Context) error {
	return utils.ResponseMessage(http.StatusOK, map[string]string{"kyc": "KYC processed successfully"}, context)
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
	existingUser, _ := service.UserRepo.GetUserByEmail(
		context.Request().Context(),
		arg.Email,
	)
	if existingUser != nil {
		utils.Error("User already exists",
			utils.LogField{Key: "email", Value: arg.Email.String})
		return nil, errors.New("user already exists")
	}

	// Create user
	createdRow, err := service.UserRepo.CreateUser(context.Request().Context(), arg)
	if err != nil {
		utils.Error("Failed to create user",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	// Convert CreateUserRow to User
	user := &db.User{
		ID:              createdRow.ID,
		Email:           createdRow.Email,
		Nin:             createdRow.Nin,
		UserType:        createdRow.UserType,
		EmailVerified:   createdRow.EmailVerified,
		EmailVerifiedAt: createdRow.EmailVerifiedAt,
		CreatedAt:       createdRow.CreatedAt,
		UpdatedAt:       createdRow.UpdatedAt,
		Fullname:        createdRow.Fullname,
	}
	return user, nil
}
