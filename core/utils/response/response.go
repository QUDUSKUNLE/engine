package response

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	// Response is the standard API response structure
	Response struct {
		Status  int           `json:"status"`
		Success bool          `json:"success"`
		Data    interface{}   `json:"data,omitempty"`
		Error   *ErrorDetails `json:"error,omitempty"`
		Meta    *Meta         `json:"meta,omitempty"`
	}

	// ErrorDetails holds error details
	ErrorDetails struct {
		Code    string `json:"code,omitempty"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	}

	// Meta holds metadata like pagination info
	Meta struct {
		Page      int `json:"page,omitempty"`
		PerPage   int `json:"per_page,omitempty"`
		Total     int `json:"total,omitempty"`
		TotalPage int `json:"total_page,omitempty"`
	}

	// PaginationParams holds pagination request parameters
	PaginationParams struct {
		Page    int `query:"page"`
		PerPage int `query:"per_page"`
	}
)
// Success sends a successful response
func Success(status int, data interface{}, c echo.Context) error {
	resp := Response{
		Status:  status,
		Success: true,
		Data:    data,
	}

	// Log success response
	zapLogger := c.Get("logger").(*zap.Logger)
	zapLogger.Info("success response",
		zap.Int("status", status),
		zap.String("path", c.Path()),
		zap.String("method", c.Request().Method),
	)

	return c.JSON(status, resp)
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(status int, data interface{}, c echo.Context, meta *Meta) error {
	resp := Response{
		Status:  status,
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	zapLogger := c.Get("logger").(*zap.Logger)
	zapLogger.Info("success response with meta",
		zap.Int("status", status),
		zap.String("path", c.Path()),
		zap.String("method", c.Request().Method),
	)

	return c.JSON(status, resp)
}

// Error sends an error response
func Error(status int, err error, c echo.Context, code string) error {
	resp := Response{
		Status:  status,
		Success: false,
		Error: &ErrorDetails{
			Code:    code,
			Message: err.Error(),
		},
	}

	zapLogger := c.Get("logger").(*zap.Logger)
	zapLogger.Error("error response",
		zap.Int("status", status),
		zap.String("path", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("error", err.Error()),
		zap.String("code", code),
	)

	return c.JSON(status, resp)
}

// ErrorWithDetails sends an error response with additional details
func ErrorWithDetails(status int, err error, c echo.Context, code string, details string) error {
	resp := Response{
		Status:  status,
		Success: false,
		Error: &ErrorDetails{
			Code:    code,
			Message: err.Error(),
			Details: details,
		},
	}

	zapLogger := c.Get("logger").(*zap.Logger)
	zapLogger.Error("error response with details",
		zap.Int("status", status),
		zap.String("path", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("error", err.Error()),
		zap.String("code", code),
		zap.String("details", details),
	)

	return c.JSON(status, resp)
}
