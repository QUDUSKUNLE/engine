package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/medivue/core/utils/response"
)

// Centralized HTTP error handler for Echo
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		status = http.StatusInternalServerError
		errMsg = err.Error()
		code   = response.CodeInternalError
	)

	// Set logger if not already set
	if c.Get("logger") == nil {
		c.Set("logger", logger)
	}

	// Handle different error types
	switch e := err.(type) {
	case *echo.HTTPError:
		status = e.Code
		if msg, ok := e.Message.(string); ok {
			errMsg = msg
		} else if msg, ok := e.Message.(error); ok {
			errMsg = msg.Error()
		} else {
			errMsg = fmt.Sprintf("%v", e.Message)
		}
		code = response.StatusToCode[status]

	case validator.ValidationErrors:
		status = http.StatusBadRequest
		errMsg = formatValidationError(e)
		code = response.CodeValidationError
	}

	if code == "" {
		code = response.StatusToCode[status]
		if code == "" {
			code = response.CodeInternalError
		}
	}

	// Log error with full context
	Error("HTTP request failed",
		LogField{Key: "error", Value: errMsg},
		LogField{Key: "error_code", Value: code},
		LogField{Key: "status_code", Value: status},
		LogField{Key: "method", Value: c.Request().Method},
		LogField{Key: "path", Value: c.Request().URL.Path},
		LogField{Key: "remote_ip", Value: c.RealIP()})

	// Convert error to our standardized format
	_ = response.Error(status, errors.New(errMsg), c, code)
}
