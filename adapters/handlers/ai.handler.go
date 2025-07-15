package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) InterpretLabHandler(context echo.Context) error {
	return h.service.InterpretLabHandler(context)
}
