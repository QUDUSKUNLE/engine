package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	validate *validator.Validate
	logger   *zap.Logger

	validatorMu sync.RWMutex

	// JWT related errors
	ErrMissingSecretKey = errors.New("missing JWT secret key")
	ErrInvalidToken     = errors.New("invalid or expired token")

	// Minimum password length for security
	MinPasswordLength = 12
)

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

	// Register validation for comparing times
	validate.RegisterStructValidation(validateTimeComparison, domain.Slots{})

	return nil
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

// ValidateTime validates if a string is in either HH:MM format or ISO 8601
func ValidateTime(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	if timeStr == "" {
		return true // Allow empty strings, use 'required' tag if needed
	}

	// First try ISO 8601 format
	if _, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return true
	}

	// Try just the date part of ISO 8601
	if _, err := time.Parse("2006-01-02T15:04:05", timeStr); err == nil {
		return true
	}

	// Try HH:MM format
	if _, err := time.Parse("15:04", timeStr); err == nil {
		return true
	}

	logger.Error("invalid time format",
		zap.String("time", timeStr),
		zap.String("expected_formats", "HH:MM or ISO 8601"))
	return false
}

// ErrorResponse sends a structured error response with logging.
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
func ResponseMessage(status int, data interface{}, c echo.Context) error {
	if c.Get("logger") == nil {
		c.Set("logger", logger) // Set logger if not present for backward compatibility
	}
	return response.Success(status, data, c)
}

// MarshalJSONField marshals any struct to JSON with error handling
func MarshalJSONField(field interface{}) ([]byte, error) {
	data, err := json.Marshal(field)
	if err != nil {
		logger.Error("json marshal failed", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return data, nil
}

// UnmarshalJSONField unmarshals JSON data with validation
func UnmarshalJSONField(data []byte, v interface{}) error {
	// Sanity check: handle empty or null values
	s := strings.TrimSpace(string(data))
	if len(s) == 0 || s == "null" {
		return nil // no error, just don't populate the target
	}

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

// validateTimeComparison implements custom validation for comparing time strings
func validateTimeComparison(sl validator.StructLevel) {
	slots := sl.Current().Interface().(domain.Slots)
	startTime := slots.StartTime
	endTime := slots.EndTime

	// Extract and compare only the time part of the day
	startTimeOnly := startTime.Sub(time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location()))
	endTimeOnly := endTime.Sub(time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location()))

	if endTimeOnly <= startTimeOnly || !endTime.After(startTime) {
		sl.ReportError(slots.EndTime, "EndTime", "end_time", "gtfield", "must be after start_time")
	}
}

// ParseTimeString converts a time string in either HH:MM or ISO 8601 format to time.Time
func ParseTimeString(timeStr string, date time.Time) (time.Time, error) {
	// Try ISO 8601 format first
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t, nil
	}

	// Try just the date part of ISO 8601
	if t, err := time.Parse("2006-01-02T15:04:05", timeStr); err == nil {
		return t, nil
	}

	// Try HH:MM format
	if t, err := time.Parse("15:04", timeStr); err == nil {
		// Combine with the provided date
		return time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			t.Hour(),
			t.Minute(),
			0, 0,
			date.Location(),
		), nil
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s (expected HH:MM or ISO 8601)", timeStr)
}
