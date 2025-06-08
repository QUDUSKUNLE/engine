package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Centralized HTTP error handler for Echo
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code             = http.StatusInternalServerError
		msg  interface{} = "Internal Server Error"
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else if ve, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		msg = ve.Error()
	} else if err != nil {
		msg = err.Error()
	}

	// You can add logging here if needed
	_ = c.JSON(code, map[string]interface{}{
		"error": msg,
	})
}
