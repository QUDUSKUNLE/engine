package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type (
	// GoogleAuthDTO represents the data sent by the frontend for Google auth
	GoogleAuthDTO struct {
		IDToken string `json:"id_token" validate:"required"`
	}
	// GoogleUserInfo represents the verified user info from Google
	GoogleUserInfo struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		GoogleID      string `json:"sub"`
	}
	UserRegisterDTO struct {
		FirstName       string      `json:"first_name" validate:"gte=3"`
		LastName        string      `json:"last_name" validate:"gte=3"`
		Email           string      `json:"email" validate:"email,required"`
		Password        string      `json:"password" validate:"gte=6,lte=20,required"`
		ConfirmPassword string      `json:"confirm_password" validate:"eqfield=Password,gte=6,lte=20,required"`
		UserType        db.UserEnum `json:"user_type" validate:"required,oneof=PATIENT DIAGNOSTIC_CENTRE_OWNER"`
		CreatedAdmin    uuid.UUID   `json:"created_admin" validate:"omitempty"`
	}
	RegisterationDTO struct {
		FirstName       string      `json:"first_name" validate:"gte=3"`
		LastName        string      `json:"last_name" validate:"gte=3"`
		Email           string      `json:"email" validate:"email,required"`
		Password        string      `json:"password" validate:"gte=6,lte=20,required"`
		ConfirmPassword string      `json:"confirm_password" validate:"eqfield=Password,gte=6,lte=20,required"`
		UserType        db.UserEnum `json:"user_type" validate:"required,oneof=PATIENT DIAGNOSTIC_CENTRE_OWNER"`
	}
	DiagnosticCentreManagerRegisterDTO struct {
		Email     string      `json:"email" validate:"email,required"`
		LastName  string      `json:"last_name" validate:"gte=3,required"`
		FirstName string      `json:"first_name" validate:"gte=3,required"`
		UserType  db.UserEnum `json:"user_type" validate:"required,oneof=DIAGNOSTIC_CENTRE_MANAGER"`
	}
	UserSignInDTO struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"gte=6,lte=20,required"`
	}
	CurrentUserDTO struct {
		UserID       uuid.UUID   `json:"user_id"`
		UserType     db.UserEnum `json:"user_type"`
		DiagnosticID uuid.UUID   `json:"diagnostic_id"`
	}
	JwtCustomClaimsDTO struct {
		UserID       uuid.UUID   `json:"user_id"`
		UserType     db.UserEnum `json:"user_type"`
		DiagnosticID uuid.UUID   `json:"diagnostic_id"`
		jwt.RegisteredClaims
	}
	ResetPasswordDTO struct {
		Email           string `json:"email" validate:"email,required"`
		Token           string `json:"token" validate:"required"`
		NewPassword     string `json:"new_password" validate:"gte=6,lte=20,required"`
		ConfirmPassword string `json:"confirm_password" validate:"eqfield=NewPassword,required"`
	}
	GetProfileDTO struct{}

	RequestPasswordResetDTO struct {
		Email string `json:"email" validate:"email,required"`
	}

	PasswordResetTokenDTO struct {
		Email     string    `json:"email"`
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	// PasswordResetToken represents a password reset token
	PasswordResetToken struct {
		ID        string
		Email     string
		Token     string
		Used      bool
		ExpiresAt time.Time
		CreatedAt time.Time
	}
	ChangePasswordDTO struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"gte=6,lte=20,required"`
		ConfirmPassword string `json:"confirm_password" validate:"eqfield=NewPassword,required"`
	}
	UpdateUserProfileDTO struct {
		FirstName   string `json:"first_name" validate:"required,min=3"`
		LastName    string `json:"last_name" validate:"required,min=3"`
		Nin         string `json:"nin" validate:"omitempty,min=11"`
		PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
	}
	GetUserProfileParamDTO struct {
		UserID uuid.UUID `json:"user_id" validate:"required"`
	}
	DeactivateAccountDTO struct {
		Password string `json:"password" validate:"required"`
		Reason   string `json:"reason" validate:"omitempty"`
	}
	EmailVerificationDTO struct {
		Email string `query:"email" validate:"email,required"`
		Token string `query:"token" validate:"required"`
	}
	ResendVerificationDTO struct {
		Email string `json:"email" validate:"email,required"`
	}
	EmailVerificationToken struct {
		ID        string
		Email     string
		Token     string
		Used      bool
		ExpiresAt time.Time
		CreatedAt time.Time
	}
	UpgradeDTO struct {
		Nin           string `json:"nin" validate:"at_least_one"`
		Passport      string `json:"passport" validate:"at_least_one"`
		DriverLicence string `json:"driver_licence" validate:"at_least_one"`
	}
)

func BuildNewUser(user UserRegisterDTO) (db.CreateUserParams, error) {
	password, err := HashPassword(user.Password)
	if err != nil {
		return db.CreateUserParams{}, err
	}

	// Handle created_admin as nullable UUID
	var createdAdmin pgtype.UUID
	if user.CreatedAdmin != uuid.Nil {
		createdAdmin = pgtype.UUID{
			Bytes: user.CreatedAdmin,
			Valid: true,
		}
	} else {
		createdAdmin = pgtype.UUID{Valid: false}
	}
	params := db.CreateUserParams{
		Email:    pgtype.Text{String: user.Email, Valid: true},
		Password: password,
		UserType: user.UserType,
		Fullname: pgtype.Text{
			String: strings.Trim(user.FirstName, "") + " " + strings.Trim(user.LastName, ""),
			Valid:  true,
		},
		PhoneNumber:  pgtype.Text{String: "", Valid: true},
		CreatedAdmin: createdAdmin,
	}
	return params, nil
}

func ComparePassword(user db.User, pass string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		return errors.New("incorrect log in credentials")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	password = string(hashPassword)
	return password, nil
}
