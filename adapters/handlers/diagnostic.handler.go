package handlers

import "github.com/labstack/echo/v4"

func (handler *HTTPHandler) CreateDiagnostic(context echo.Context) error {
	return handler.service.CreateDiagnsoticCentre(context)
}
