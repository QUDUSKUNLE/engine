package handlers

import (
	"github.com/labstack/echo/v4"
)

func (handler *HTTPHandler) CreateSchedule(context echo.Context) error {
	return handler.service.CreateSchedule(context)
}

func (handler *HTTPHandler) GetSchedule(context echo.Context) error {
	return handler.service.GetDiagnosticSchedule(context)
}

func (handler *HTTPHandler) GetSchedules(context echo.Context) error {
	return handler.service.GetDiagnosticSchedules(context)
}

func (handler *HTTPHandler) UpdateSchedule(context echo.Context) error {
	return handler.service.UpdateDiagnosticSchedule(context)
}

func (handler *HTTPHandler) DeleteDiagnosticSchedule(context echo.Context) error {
	return nil
}

func (handler *HTTPHandler) GetDiagnosticScheduleByCentre(context echo.Context) error {
	return handler.service.GetDiagnosticScheduleByCentre(context)
}

func (handler *HTTPHandler) GetDiagnosticSchedulesByCentre(context echo.Context) error {
	return handler.service.GetDiagnosticSchedulesByCentre(context)
}

func (handler *HTTPHandler) UpdateDiagnosticScheduleByCentre(context echo.Context) error {
	return handler.service.UpdateDiagnosticScheduleByCentre(context)
}
