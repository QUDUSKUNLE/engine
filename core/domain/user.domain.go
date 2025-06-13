package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medicue/adapters/db"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserRegisterDTO struct {
		Email            string      `json:"email" validate:"email,required"`
		Password         string      `json:"password" validate:"gte=6,lte=20,required"`
		ConfirmPassword  string      `json:"confirm_password" validate:"eqfield=Password,gte=6,lte=20,required"`
		UserType         db.UserEnum `json:"user_type" validate:"required"`
		DiagnosticCentre string      `json:"diagnostic_Centre"`
	}
	DiagnosticCentreManagerRegisterDTO struct {
		Email    string      `json:"email" validate:"email,required"`
		UserType db.UserEnum `json:"user_type" validate:"required"`
	}
	UserSignInDTO struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"gte=6,lte=20,required"`
	}
	CurrentUserDTO struct {
		UserID       uuid.UUID   `json:"user_id"`
		DiagnosticID uuid.UUID   `json:"diagnostic_id"`
		UserType     db.UserEnum `json:"user_type"`
	}
	JwtCustomClaimsDTO struct {
		UserID       uuid.UUID   `json:"user_id"`
		DiagnosticID uuid.UUID   `json:"diagnostic_id"`
		UserType     db.UserEnum `json:"user_type"`
		jwt.RegisteredClaims
	}
	ResetPasswordDTO struct {
		Email           string `json:"email" validate:"email,required"`
		Token           string `json:"token" validate:"required"`
		NewPassword     string `json:"new_password" validate:"gte=6,lte=20,required"`
		ConfirmPassword string `json:"confirm_password" validate:"eqfield=NewPassword,required"`
	}

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
		FullName string `json:"full_name" validate:"required,min=2"`
		NIN      string `json:"nin" validate:"omitempty,min=11"`
	}

	GetUserProfileParamDTO struct {
		UserID uuid.UUID `json:"user_id" validate:"required"`
	}

	DeactivateAccountDTO struct {
		Password string `json:"password" validate:"required"`
		Reason   string `json:"reason" validate:"omitempty"`
	}

	EmailVerificationDTO struct {
		Email string `json:"email" validate:"email,required"`
		Token string `json:"token" validate:"required"`
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
)

func BuildNewUser(user UserRegisterDTO) (db.CreateUserParams, error) {
	password, err := HashPassword(user.Password)
	if err != nil {
		return db.CreateUserParams{}, err
	}
	return db.CreateUserParams{
		Email:    pgtype.Text{String: user.Email, Valid: true},
		Password: password,
		UserType: user.UserType,
	}, nil
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
