package handlers

import "github.com/labstack/echo/v4"

func (handler *HTTPHandler) CreateDiagnostic(context echo.Context) error {
	return handler.service.CreateDiagnsoticCentre(context)
}

func (handler *HTTPHandler) GetDiagnosticCentre(context echo.Context) error {
	return handler.service.GetDiagnosticCentre(context)
}

func (handler *HTTPHandler) SearchDiagnosticCentre(context echo.Context) error {
	return handler.service.SearchDiagnosticCentre(context)
}
