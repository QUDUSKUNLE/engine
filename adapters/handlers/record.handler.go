package handlers

import (
	"github.com/labstack/echo/v4"
)

func (handler *HTTPHandler) CreateMedicalRecord(context echo.Context) error {
	return handler.service.CreateMedicalRecord(context)
}

func (handler *HTTPHandler) GetMedicalRecord(context echo.Context) error {
	return handler.service.GetMedicalRecord(context)
}

func (handler *HTTPHandler) GetMedicalRecords(context echo.Context) error {
	return handler.service.GetMedicalRecords(context)
}

func (handler *HTTPHandler) GetUploaderMedicalRecord(context echo.Context) error {
	return handler.service.GetUploaderMedicalRecord(context)
}

func (handler *HTTPHandler) GetUploaderMedicalRecords(context echo.Context) error {
	return handler.service.GetUploaderMedicalRecords(context)
}

func (handler *HTTPHandler) UpdateMedicalRecord(context echo.Context) error {
	return handler.service.UpdateMedicalRecord(context)
}
