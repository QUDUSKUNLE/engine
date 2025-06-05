package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
)

func ErrorResponse(status int, message error, context echo.Context) error {
	return context.JSON(status, echo.Map{"error": message.Error()})
}

func ResponseMessage(status int, message interface{}, context echo.Context) error {
	return context.JSON(status, echo.Map{"result": message})
}

// GenerateToken generates a JWT token for the given user
func GenerateToken(user domain.CurrentUserDTO) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		ErrMissingSecretKey := errors.New("missing JWT secret key")
		return "", ErrMissingSecretKey
	}

	var diagnosticID uuid.UUID
	if user.DiagnosticID != "" {
		parsedUUID, err := uuid.Parse(user.DiagnosticID)
		if err != nil {
			return "", err
		}
		diagnosticID = parsedUUID
	}
	var userID uuid.UUID
	if user.UserID != "" {
		parsedUserID, err := uuid.Parse(user.UserID)
		if err != nil {
			return "", err
		}
		userID = parsedUserID
	}
	claims := &domain.JwtCustomClaimsDTO{
		UserID:       userID,
		DiagnosticID: diagnosticID,
		UserType:     db.UserEnum(user.UserType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?/|"

func GenerateRandomPassword(length int) (string) {
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			return ""
		}
		password[i] = passwordChars[num.Int64()]
	}
	return string(password)
}
