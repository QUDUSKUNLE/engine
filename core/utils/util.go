package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	claims := &domain.JwtCustomClaimsDTO{
		UserID:       user.UserID,
		DiagnosticID: user.DiagnosticID,
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

func GenerateRandomPassword(length int) string {
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

func currentUser(context echo.Context) (*domain.CurrentUserDTO, error) {
	user := context.Get("user").(*jwt.Token)
	claims, ok := user.Claims.(*domain.JwtCustomClaimsDTO)
	if !ok {
		return nil, errors.New("token is required")
	}
	return &domain.CurrentUserDTO{
		UserID:       claims.UserID,
		DiagnosticID: claims.DiagnosticID,
		UserType:     claims.UserType,
	}, nil
}

func PrivateMiddlewareContext(context echo.Context, userType string) (*domain.CurrentUserDTO, error) {
	user, err := currentUser(context)
	if err != nil {
		return nil, err
	}
	// Check user type
	if user.UserType != db.UserEnum(userType) {
		return nil, errors.New("unauthorized to perform this operation")
	}
	return user, nil
}
