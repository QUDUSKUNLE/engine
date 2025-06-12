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
)

// Custom validator
type CustomValidator struct {
	validator *validator.Validate
}

// Custom validator
func (c *CustomValidator) Validate(inter interface{}) error {
	if err := c.validator.Struct(inter); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.(validator.ValidationErrors))
		}
		var errorMessage []map[string]string
		for _, er := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, map[string]string{"field": er.Field(), "message": fmt.Sprintf("%s is an invalid input for field: %s", er.Value(), er.Field())})
		}
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage)
	}
	return nil
}

// ValidationAdaptor
func ValidationAdaptor(xx *echo.Echo) *echo.Echo {
	xx.Validator = &CustomValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
	return xx
}

// Helper to handle DTO binding and validation
func bindAndValidateDTO(c echo.Context, dtoFactory func() interface{}, bindFunc func(interface{}) error, setKey string) error {
	dto := dtoFactory()
	if err := bindFunc(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", err))
	}

	// Special handling for CreateMedicalRecordDTO: manually parse UUIDs and other fields from form-data
	if v, ok := dto.(*domain.CreateMedicalRecordDTO); ok {
		if userID := c.FormValue("user_id"); userID != "" {
			v.UserID, _ = uuid.Parse(userID)
		}
		if uploaderID := c.FormValue("uploader_id"); uploaderID != "" {
			v.UploaderID, _ = uuid.Parse(uploaderID)
		}
		if scheduleID := c.FormValue("schedule_id"); scheduleID != "" {
			v.ScheduleID, _ = uuid.Parse(scheduleID)
		}
		if uploaderAdminID := c.FormValue("uploader_admin_id"); uploaderAdminID != "" {
			v.UploaderAdminID, _ = uuid.Parse(uploaderAdminID)
		}
		v.Title = c.FormValue("title")
		v.DocumentType = db.DocumentType(c.FormValue("document_type"))
		v.UploadedAt = c.FormValue("uploaded_at")
		v.ProviderName = c.FormValue("provider_name")
		v.Specialty = c.FormValue("specialty")
		v.IsShared = c.FormValue("is_shared") == "true"
		v.SharedUntil = c.FormValue("shared_until")
	}

	if err := c.Validate(dto); err != nil {
		return err
	}
	c.Set(setKey, dto)
	return nil
}

// Generic body validation interceptor for any DTO
func BodyValidationInterceptorFor(dtoFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch c.Request().Method {
			case http.MethodGet, http.MethodDelete:
				if err := bindAndValidateDTO(c, dtoFactory, c.Bind, "validatedQueryParamDTO"); err != nil {
					return err
				}
			case http.MethodPost, http.MethodPut:
				contentType := c.Request().Header.Get("Content-Type")
				var bindFunc func(interface{}) error
				if strings.HasPrefix(contentType, "application/json") {
					bodyBytes, err := io.ReadAll(c.Request().Body)
					if err != nil {
						return echo.NewHTTPError(http.StatusBadRequest, "Failed to read request body")
					}
					c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					bindFunc = func(dto interface{}) error {
						return json.Unmarshal(bodyBytes, dto)
					}
				} else {
					if strings.HasPrefix(contentType, "multipart/form-data") {
						if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
							return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse multipart form")
						}
					}
					bindFunc = c.Bind // Use c.Bind for form-data and urlencoded
				}
				if err := bindAndValidateDTO(c, dtoFactory, bindFunc, "validatedBodyDTO"); err != nil {
					return err
				}
			}
			return next(c)
		}
	}
}
