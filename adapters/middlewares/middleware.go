package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
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
				if err := bindAndValidateDTO(c, dtoFactory, c.Bind, "validatedQueryParamsDTO"); err != nil {
					return err
				}
			case http.MethodPost:
				bodyBytes, err := io.ReadAll(c.Request().Body)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Failed to read request body")
				}
				c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				bindFunc := func(dto interface{}) error {
					return json.Unmarshal(bodyBytes, dto)
				}
				if err := bindAndValidateDTO(c, dtoFactory, bindFunc, "validatedBodyDTO"); err != nil {
					return err
				}
			}
			return next(c)
		}
	}
}
