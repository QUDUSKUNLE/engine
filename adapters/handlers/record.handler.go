package handlers

import (
	"github.com/labstack/echo/v4"
)

// CreateMedicalRecord godoc
// @Summary Upload a new medical record
// @Description Upload a new medical record with metadata and file attachment
// @Tags MedicalRecord
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param user_id formData string true "User ID" format(uuid)
// @Param uploader_id formData string false "Uploader ID (for diagnostic centres)" format(uuid)
// @Param schedule_id formData string false "Associated Schedule ID" format(uuid)
// @Param title formData string true "Record title"
// @Param document_type formData string true "Type of medical document" Enums(LAB_REPORT, PRESCRIPTION, IMAGING, DISCHARGE_SUMMARY, OTHER)
// @Param document_date formData string true "Date of the document" format(date)
// @Param provider_name formData string false "Healthcare provider name"
// @Param specialty formData string false "Medical specialty"
// @Param file formData file true "Medical record file (PDF/Image)"
// @Success 201 {object} handlers.MedicalRecordSwagger "Medical record created successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "Invalid input data/file format"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 403 {object} handlers.ErrorResponse "Insufficient permissions"
// @Failure 409 {object} handlers.DUPLICATE_ERROR "DUPLICATE_ERROR"
// @Failure 413 {object} handlers.ErrorResponse "File too large"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/medical_records [post]
func (handler *HTTPHandler) CreateMedicalRecord(context echo.Context) error {
	return handler.service.CreateMedicalRecord(context)
}

// GetMedicalRecord godoc
// @Summary Get a medical record by ID
// @Description Retrieve a specific medical record. Access limited to record owner or authorized diagnostic centre.
// @Tags MedicalRecord
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param record_id path string true "Medical Record ID" format(uuid)
// @Success 200 {object} handlers.MedicalRecordSwagger "Medical record details"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/medical_records/{record_id} [get]
func (handler *HTTPHandler) GetMedicalRecord(context echo.Context) error {
	return handler.service.GetMedicalRecord(context)
}

// GetMedicalRecords godoc
// @Summary List user's medical records
// @Description Get all medical records for the authenticated user with pagination and filtering
// @Tags MedicalRecord
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Records per page" minimum(1) maximum(100) default(10)
// @Param document_type query string false "Filter by document type" Enums(LAB_REPORT, PRESCRIPTION, IMAGING, DISCHARGE_SUMMARY, OTHER)
// @Param from_date query string false "Filter by date from (YYYY-MM-DD)" format(date)
// @Param to_date query string false "Filter by date to (YYYY-MM-DD)" format(date)
// @Success 200 {array} handlers.MedicalRecordSwagger "List of medical records"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/medical_records [get]
func (handler *HTTPHandler) GetMedicalRecords(context echo.Context) error {
	return handler.service.GetMedicalRecords(context)
}

// GetUploaderMedicalRecord godoc
// @Summary Get an uploaded medical record
// @Description Retrieve a medical record uploaded by a diagnostic centre
// @Tags MedicalRecord
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param record_id path string true "Medical Record ID" format(uuid)
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Success 200 {object} handlers.MedicalRecordSwagger "Medical record details"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/medical_records/{record_id}/diagnostic_centre/{diagnostic_centre_id} [get]
func (handler *HTTPHandler) GetUploaderMedicalRecord(context echo.Context) error {
	return handler.service.GetUploaderMedicalRecord(context)
}

// GetUploaderMedicalRecords godoc
// @Summary List uploaded medical records
// @Description Get all medical records uploaded by a diagnostic centre with pagination
// @Tags MedicalRecord
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Records per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} handlers.MedicalRecordSwagger "List of medical records"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/medical_records/diagnostic_centre/{diagnostic_centre_id} [get]
func (handler *HTTPHandler) GetUploaderMedicalRecords(context echo.Context) error {
	return handler.service.GetUploaderMedicalRecords(context)
}

// UpdateMedicalRecord godoc
// @Summary Update a medical record
// @Description Update metadata of an existing medical record. File content cannot be updated.
// @Tags MedicalRecord
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param record_id path string true "Medical Record ID" format(uuid)
// @Param record body domain.UpdateMedicalRecordDTO true "Updated record details"
// @Success 200 {object} handlers.MedicalRecordSwagger "Record updated successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/medical_records/{record_id} [put]
func (handler *HTTPHandler) UpdateMedicalRecord(context echo.Context) error {
	return handler.service.UpdateMedicalRecord(context)
}
