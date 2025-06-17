package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils/response"
	"go.uber.org/zap"
)

var (
	validate *validator.Validate
	logger   *zap.Logger

	validatorMu sync.RWMutex

	// JWT related errors
	ErrMissingSecretKey = errors.New("missing JWT secret key")
	ErrInvalidToken     = errors.New("invalid or expired token")
	ErrUnauthorized     = errors.New("unauthorized to perform this operation")

	// Minimum password length for security
	MinPasswordLength = 12
)

// InitValidator initializes the validator with custom validations
func InitValidator() error {
	validatorMu.Lock()
	defer validatorMu.Unlock()

	// Initialize validator
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validations
	if err := validate.RegisterValidation("min_one", validateMinOne); err != nil {
		return fmt.Errorf("failed to register min_one validator: %w", err)
	}

	return nil
}

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	if err := InitValidator(); err != nil {
		panic(fmt.Sprintf("failed to initialize validator: %v", err))
	}
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	validatorMu.RLock()
	defer validatorMu.RUnlock()
	return validate
}

// Add this new function to handle min_one validation
func validateMinOne(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return field.Len() > 0
	}
	return false
}

// ErrorResponse sends a structured error response with logging.
// Deprecated: Use response.Error instead
func ErrorResponse(status int, err error, c echo.Context) error {
	if c.Get("logger") == nil {
		c.Set("logger", logger) // Set logger if not present for backward compatibility
	}

	code := response.StatusToCode[status]
	if code == "" {
		code = response.CodeInternalError
	}

	return response.Error(status, err, c, code)
}

// ResponseMessage sends a structured success response with logging.
// Deprecated: Use response.Success instead
func ResponseMessage(status int, data interface{}, c echo.Context) error {
	if c.Get("logger") == nil {
		c.Set("logger", logger) // Set logger if not present for backward compatibility
	}
	return response.Success(status, data, c)
}

// GenerateToken generates a JWT token for the given user with enhanced security
func GenerateToken(user domain.CurrentUserDTO) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		logger.Error("jwt secret key missing")
		return "", ErrMissingSecretKey
	}

	// Get token expiration from env or use default
	tokenExpiration := 72 * time.Hour
	if expStr := os.Getenv("JWT_EXPIRATION_HOURS"); expStr != "" {
		if exp, err := time.ParseDuration(expStr + "h"); err == nil {
			tokenExpiration = exp
		}
	}

	claims := &domain.JwtCustomClaimsDTO{
		UserID:       user.UserID,
		DiagnosticID: user.DiagnosticID,
		UserType:     db.UserEnum(user.UserType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Error("failed to sign token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("token generated successfully",
		zap.String("userID", user.UserID.String()),
		zap.String("userType", string(user.UserType)),
	)
	return signedToken, nil
}

const (
	passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?/|"
)

// GenerateRandomPassword generates a cryptographically secure random password
func GenerateRandomPassword(length int) (string, error) {
	if length < MinPasswordLength {
		return "", fmt.Errorf("password length must be at least %d characters", MinPasswordLength)
	}

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			logger.Error("failed to generate random number", zap.Error(err))
			return "", fmt.Errorf("failed to generate password: %w", err)
		}
		password[i] = passwordChars[num.Int64()]
	}
	return string(password), nil
}

// CurrentUser extracts the current user from the JWT token in the context
func CurrentUser(c echo.Context) (*domain.CurrentUserDTO, error) {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		logger.Error("missing or invalid user token")
		return nil, ErrInvalidToken
	}

	claims, ok := userToken.Claims.(*domain.JwtCustomClaimsDTO)
	if !ok {
		logger.Error("invalid token claims")
		return nil, ErrInvalidToken
	}

	logger.Debug("user extracted from token",
		zap.String("userID", claims.UserID.String()),
		zap.String("userType", string(claims.UserType)),
	)

	return &domain.CurrentUserDTO{
		UserID:       claims.UserID,
		DiagnosticID: claims.DiagnosticID,
		UserType:     claims.UserType,
	}, nil
}

// PrivateMiddlewareContext validates user type access
func PrivateMiddlewareContext(c echo.Context, userTypes []db.UserEnum) (*domain.CurrentUserDTO, error) {
	user, err := CurrentUser(c)
	if err != nil {
		return nil, err
	}

	for _, allowedType := range userTypes {
		if user.UserType == allowedType {
			return user, nil
		}
	}

	logger.Warn("unauthorized access attempt",
		zap.Strings("requiredTypes", convertUserTypesToStrings(userTypes)),
		zap.String("actualType", string(user.UserType)),
		zap.String("userID", user.UserID.String()),
	)
	return nil, ErrUnauthorized
}

// Helper function to convert UserEnum slice to string slice for logging
func convertUserTypesToStrings(types []db.UserEnum) []string {
	strings := make([]string, len(types))
	for i, t := range types {
		strings[i] = string(t)
	}
	return strings
}

// MarshalJSONField marshals any struct to JSON with error handling
func MarshalJSONField(field interface{}, c echo.Context) ([]byte, error) {
	data, err := json.Marshal(field)
	if err != nil {
		logger.Error("json marshal failed", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return data, nil
}

// UnmarshalJSONField unmarshals JSON data with validation
func UnmarshalJSONField(data []byte, v interface{}, c echo.Context) error {
	if err := json.Unmarshal(data, v); err != nil {
		logger.Error("json unmarshal failed", zap.Error(err))
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if v, ok := v.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			logger.Error("validation failed", zap.Error(err))
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// ValidateParams validates URL or query params with detailed errors
func ValidateParams(c echo.Context, params interface{}) error {
	if err := c.Bind(params); err != nil {
		logger.Error("parameter binding failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameters format")
	}

	if err := validate.Struct(params); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		logger.Error("parameter validation failed",
			zap.Any("validationErrors", validationErrors),
		)
		return echo.NewHTTPError(http.StatusBadRequest, formatValidationError(validationErrors))
	}

	return nil
}

// formatValidationError formats validation errors into user-friendly messages
func formatValidationError(errs validator.ValidationErrors) string {
	var errMsg string
	for _, err := range errs {
		if errMsg != "" {
			errMsg += "; "
		}
		errMsg += fmt.Sprintf(
			"Field '%s' failed '%s' validation",
			err.Field(),
			err.Tag(),
		)
	}
	return errMsg
}

// GenerateRandomToken creates a cryptographically secure random token
func GenerateRandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
