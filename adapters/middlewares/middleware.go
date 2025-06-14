package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

const (
	maxFileSize       = 32 << 20 // 32MB
	validatedBodyKey  = "validatedBodyDTO"
	validatedQueryKey = "validatedQueryParamDTO"
)

// Custom validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate handles struct validation with detailed error messages
func (c *CustomValidator) Validate(inter interface{}) error {
	if err := c.validator.Struct(inter); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			utils.Error("Invalid validation error",
				utils.LogField{Key: "error", Value: err.Error()})
			return echo.NewHTTPError(http.StatusBadRequest, err.(validator.ValidationErrors))
		}

		var errorMessage []map[string]string
		for _, er := range err.(validator.ValidationErrors) {
			errDetail := map[string]string{
				"field":   er.Field(),
				"message": fmt.Sprintf("%v is an invalid input for field: %s", er.Value(), er.Field()),
				"tag":     er.Tag(),
			}
			errorMessage = append(errorMessage, errDetail)

			utils.Error("Validation failed",
				utils.LogField{Key: "field", Value: er.Field()},
				utils.LogField{Key: "value", Value: er.Value()},
				utils.LogField{Key: "tag", Value: er.Tag()})
		}
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage)
	}
	return nil
}

// ValidationAdaptor initializes and sets up the validator
func ValidationAdaptor(e *echo.Echo) *echo.Echo {
	utils.Info("Initializing validation adapter")
	e.Validator = &CustomValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
	return e
}

// bindAndValidateDTO handles DTO binding and validation with proper error handling
func bindAndValidateDTO(c echo.Context, dtoFactory func() interface{}, bindFunc func(interface{}) error, setKey string) error {
	dto := dtoFactory()

	if err := bindFunc(dto); err != nil {
		utils.Error("Failed to bind request data",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "path", Value: c.Path()},
			utils.LogField{Key: "method", Value: c.Request().Method})
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", err))
	}

	// Handle medical record creation specially
	if v, ok := dto.(*domain.CreateMedicalRecordDTO); ok {
		if err := handleMedicalRecordDTO(c, v); err != nil {
			return err
		}
	}

	// Handle update medical records
	if v, ok := dto.(*domain.UpdateMedicalRecordDTO); ok {
		if err := handleUpdateMedicalRecordDTO(c, v); err != nil {
			return err
		}
	}

	if err := c.Validate(dto); err != nil {
		return err
	}

	utils.Info("Request data validated successfully",
		utils.LogField{Key: "path", Value: c.Path()},
		utils.LogField{Key: "method", Value: c.Request().Method},
		utils.LogField{Key: "dto_type", Value: fmt.Sprintf("%T", dto)})

	c.Set(setKey, dto)
	return nil
}

// handleMedicalRecordDTO processes form data for medical record creation
func handleMedicalRecordDTO(c echo.Context, dto *domain.CreateMedicalRecordDTO) error {
	fields := map[string]*uuid.UUID{
		"user_id":           &dto.UserID,
		"uploader_id":       &dto.UploaderID,
		"schedule_id":       &dto.ScheduleID,
		"uploader_admin_id": &dto.UploaderAdminID,
	}

	for field, ptr := range fields {
		if value := c.FormValue(field); value != "" {
			parsed, err := uuid.Parse(value)
			if err != nil {
				utils.Error("Failed to parse UUID",
					utils.LogField{Key: "field", Value: field},
					utils.LogField{Key: "value", Value: value},
					utils.LogField{Key: "error", Value: err.Error()})
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid %s format", field))
			}
			*ptr = parsed
		}
	}

	// Set other form values
	dto.Title = c.FormValue("title")
	dto.DocumentType = db.DocumentType(c.FormValue("document_type"))
	dto.UploadedAt = c.FormValue("uploaded_at")
	dto.ProviderName = c.FormValue("provider_name")
	dto.Specialty = c.FormValue("specialty")
	dto.IsShared = c.FormValue("is_shared") == "true"
	dto.SharedUntil = c.FormValue("shared_until")

	return nil
}

func handleUpdateMedicalRecordDTO(c echo.Context, dto *domain.UpdateMedicalRecordDTO) error {
	fields := map[string]*uuid.UUID{
		"uploader_id":       &dto.UploaderID,
		"uploader_admin_id": &dto.UploaderAdminID,
	}

	for field, ptr := range fields {
		if value := c.FormValue(field); value != "" {
			parsed, err := uuid.Parse(value)
			if err != nil {
				utils.Error("Failed to parse UUID",
					utils.LogField{Key: "field", Value: field},
					utils.LogField{Key: "value", Value: value},
					utils.LogField{Key: "error", Value: err.Error()})
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid %s format", field))
			}
			*ptr = parsed
		}
	}

	// Set other form values
	dto.Title = c.FormValue("title")
	dto.DocumentType = db.DocumentType(c.FormValue("document_type"))
	dto.UploadedAt = c.FormValue("uploaded_at")
	dto.ProviderName = c.FormValue("provider_name")
	dto.Specialty = c.FormValue("specialty")
	dto.IsShared = c.FormValue("is_shared") == "true"
	dto.SharedUntil = c.FormValue("shared_until")

	return nil
}

// BodyValidationInterceptorFor provides middleware for request validation
func BodyValidationInterceptorFor(dtoFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			utils.Info("Processing request",
				utils.LogField{Key: "path", Value: c.Path()},
				utils.LogField{Key: "method", Value: c.Request().Method})

			switch c.Request().Method {
			case http.MethodGet, http.MethodDelete:
				if err := bindAndValidateDTO(c, dtoFactory, c.Bind, validatedQueryKey); err != nil {
					return err
				}

			case http.MethodPost, http.MethodPut:
				contentType := c.Request().Header.Get("Content-Type")

				var bindFunc func(interface{}) error
				if strings.HasPrefix(contentType, "application/json") {
					bindFunc = handleJSONRequest(c)
				} else {
					if strings.HasPrefix(contentType, "multipart/form-data") {
						if err := c.Request().ParseMultipartForm(maxFileSize); err != nil {
							utils.Error("Failed to parse multipart form",
								utils.LogField{Key: "error", Value: err.Error()},
								utils.LogField{Key: "content_type", Value: contentType})
							return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse multipart form")
						}
					}
					bindFunc = c.Bind
				}

				if err := bindAndValidateDTO(c, dtoFactory, bindFunc, validatedBodyKey); err != nil {
					return err
				}
			}

			return next(c)
		}
	}
}

// handleJSONRequest processes JSON requests with proper error handling
func handleJSONRequest(c echo.Context) func(interface{}) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		utils.Error("Failed to read request body",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil
	}
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return func(dto interface{}) error {
		return json.Unmarshal(bodyBytes, dto)
	}
}
