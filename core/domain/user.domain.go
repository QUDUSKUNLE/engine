package domain

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medicue/adapters/db"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	UserRegisterDTO struct {
		Email            string      `json:"email" validate:"email,required"`
		Password         string      `json:"password" validate:"gte=6,lte=20,required"`
		ConfirmPassword  string      `json:"confirm_password" validate:"eqfield=Password,gte=6,lte=20,required"`
		UserType         db.UserEnum `json:"user_type" validate:"required"`
		DiagnosticCentre string      `json:"diagnostic_Centre"`
	}
	UserSignInDTO struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"gte=6,lte=20,required"`
	}
	CurrentUserDTO struct {
		UserID       string  `json:"user_id"`
		DiagnosticID string `json:"diagnostic_id"`
		UserType     string  `json:"user_type"`
	}
	JwtCustomClaimsDTO struct {
		UserID       uuid.UUID `json:"user_id"`
		DiagnosticID uuid.UUID `json:"diagnostic_id"`
		UserType     db.UserEnum   `json:"user_type"`
		jwt.RegisteredClaims
	}
)

func BuildNewUser(user UserRegisterDTO) (*db.User, error) {
	Password, err := HashPassword(user.Password)
	if err != nil {
		return &db.User{}, err
	}
	return &db.User{
		Email:    pgtype.Text{String: user.Email, Valid: true},
		Password: Password,
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
