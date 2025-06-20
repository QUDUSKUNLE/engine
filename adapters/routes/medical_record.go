package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/core/domain"
)

// MedicalRecordRoutes registers all medical record-related routes
func MedicalRecordRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	recordGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/medical_records",
			handler:     handler.CreateMedicalRecord,
			factory:     func() interface{} { return &domain.CreateMedicalRecordDTO{} },
			description: "Create medical record",
		},
		{
			method:      http.MethodGet,
			path:        "/medical_records/:record_id",
			handler:     handler.GetMedicalRecord,
			factory:     func() interface{} { return &domain.GetMedicalRecordParamsDTO{} },
			description: "Get medical record details",
		},
		{
			method:      http.MethodGet,
			path:        "/medical_records",
			handler:     handler.GetMedicalRecords,
			factory:     func() interface{} { return &domain.GetMedicalRecordsParamQueryDTO{} },
			description: "Get all medical records",
		},
		{
			method:      http.MethodGet,
			path:        "/medical_records/:record_id/diagnostic_centre/:diagnostic_centre_id",
			handler:     handler.GetUploaderMedicalRecord,
			factory:     func() interface{} { return &domain.GetUploaderMedicalRecordParamsDTO{} },
			description: "Get uploader medical record",
		},
		{
			method:      http.MethodGet,
			path:        "/medical_records/diagnostic_centre/:diagnostic_centre_id",
			handler:     handler.GetUploaderMedicalRecords,
			factory:     func() interface{} { return &domain.GetUploaderMedicalRecordsParamQueryDTO{} },
			description: "Get uploader medical records by centre",
		},
		{
			method:      http.MethodPut,
			path:        "/medical_records",
			handler:     handler.UpdateMedicalRecord,
			factory:     func() interface{} { return &domain.UpdateMedicalRecordDTO{} },
			description: "Update a medical record",
		},
	}

	registerRoutes(group, recordGroup)
}
