package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
)

var validate = validator.New()

func ErrorResponse(status int, err error, c echo.Context) error {
	// Optionally log the error here
	return c.JSON(status, echo.Map{"error": err.Error()})
}

func ResponseMessage(status int, message interface{}, c echo.Context) error {
	return c.JSON(status, echo.Map{"data": message})
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

func CurrentUser(c echo.Context) (*domain.CurrentUserDTO, error) {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return nil, errors.New("user token is missing or invalid")
	}
	claims, ok := userToken.Claims.(*domain.JwtCustomClaimsDTO)
	if !ok {
		return nil, errors.New("token claims are invalid")
	}
	return &domain.CurrentUserDTO{
		UserID:       claims.UserID,
		DiagnosticID: claims.DiagnosticID,
		UserType:     claims.UserType,
	}, nil
}

func PrivateMiddlewareContext(c echo.Context, userType string) (*domain.CurrentUserDTO, error) {
	user, err := CurrentUser(c)
	if err != nil {
		return nil, err
	}
	if user.UserType != db.UserEnum(userType) {
		return nil, errors.New("unauthorized to perform this operation")
	}
	return user, nil
}

// MarshalField marshals any struct to JSON and returns the bytes or an error
func MarshalField(field interface{}) ([]byte, error) {
	return json.Marshal(field)
}

// ValidateParams validates URL or query params bound to a struct
func ValidateParams(c echo.Context, params interface{}) error {
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid params")
	}
	if err := validate.Struct(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Helper to marshal JSON fields and handle errors
func MarshalJSONField(field interface{}, c echo.Context) ([]byte, error) {
	data, err := json.Marshal(field)
	if err != nil {
		return nil, ErrorResponse(http.StatusInternalServerError, err, c)
	}
	return data, nil
}

// Unmarshal to JSON
func UnmarshalJSONField(data []byte, v interface{}, c echo.Context) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err, c)
	}
	return nil
}
