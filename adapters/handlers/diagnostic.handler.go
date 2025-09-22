package handlers

import (
	"github.com/labstack/echo/v4"
)

// CreateDiagnostic godoc
// @Summary Create a new diagnostic centre
// @Description Create a new diagnostic centre with location, contact details, and available services
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param RequestBody body domain.CreateDiagnosticDTO true "Diagnostic centre details"
// @Success 201 {object} handlers.DiagnosticCentreSwagger "Diagnostic centre created successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres [post]
func (handler *HTTPHandler) CreateDiagnostic(context echo.Context) error {
	return handler.service.CreateDiagnosticCentre(context)
}

// CreateDiagnosticCentreManager godoc
// @Summary Create a diagnostic centre manager
// @Description Create a new diagnostic centre manager account. Only accessible by diagnostic centre owners.
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param RequestBody body domain.DiagnosticCentreManagerRegisterDTO true "Manager details"
// @Success 201 {object} handlers.ManagerSwagger "Manager account created successfully"
// @Success 202 {object} map[string]string "Manager invite sent successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 409 {object} handlers.DUPLICATE_ERROR "DUPLICATE_ERROR"
// @Failure 422 {object} handlers.UNPROCESSED_ERROR "UNPROCESSED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/manager [post]
func (handler *HTTPHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	return handler.service.CreateDiagnosticCentreManager(context)
}

// GetDiagnosticCentre godoc
// @Summary Get a diagnostic centre by ID
// @Description Retrieve detailed information about a specific diagnostic centre
// @Tags DiagnosticCentre
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID (UUID format)" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Success 200 {object} handlers.DiagnosticCentreSwagger "Diagnostic centre details retrieved successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id} [get]
func (handler *HTTPHandler) GetDiagnosticCentre(context echo.Context) error {
	return handler.service.GetDiagnosticCentre(context)
}

// SearchDiagnosticCentre godoc
// @Summary Search for diagnostic centres
// @Description Search for diagnostic centres based on location, available doctors, and test types
// @Tags DiagnosticCentre
// @Produce json
// @Param latitude query number true "Latitude (-90 to 90)" minimum(-90) maximum(90) default(25.06)
// @Param longitude query number true "Longitude (-180 to 180)" minimum(-180) maximum(180) default(56.67)
// @Param day_of_week query string true "Filter by day" Enums(monday,tuesday,wednesday,thursday,friday,saturday,sunday)
// @Param doctor query string false "Filter by doctor specialization" Enums(Male,Female)
// @Param available_tests query string true "Filter by available test type" Enums(BLOOD_TEST,URINE_TEST,X_RAY,MRI,CT_SCAN,ULTRASOUND,ECG,EEG,BIOPSY,SKIN_TEST,ALLERGY_TEST,GENETIC_TEST,IMMUNOLOGY_TEST,HORMONE_TEST,VIRAL_TEST,BACTERIAL_TEST,PARASITIC_TEST,FUNGAL_TEST,MOLECULAR_TEST,TOXICOLOGY_TEST,ECHO,COVID_19_TEST,OTHER,BLOOD_SUGAR_TEST,LIPID_PROFILE,HEMOGLOBIN_TEST,THYROID_TEST,LIVER_FUNCTION_TEST,KIDNEY_FUNCTION_TEST,URIC_ACID_TEST,VITAMIN_D_TEST,VITAMIN_B12_TEST,HEMOGRAM,COMPLETE_BLOOD_COUNT,BLOOD_GROUPING,HEPATITIS_B_TEST,HEPATITIS_C_TEST,HIV_TEST,MALARIA_TEST,DENGUE_TEST,TYPHOID_TEST,COVID_19_ANTIBODY_TEST,COVID_19_RAPID_ANTIGEN_TEST,COVID_19_RT_PCR_TEST,PREGNANCY_TEST)
// @Param page query integer false "Page number for pagination" default(1) minimum(1)
// @Param per_page query integer false "Number of results per page" default(10) minimum(1) maximum(100)
// @Success 200 {array} handlers.DiagnosticCentreSwagger "List of diagnostic centres"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres [get]
func (handler *HTTPHandler) SearchDiagnosticCentre(context echo.Context) error {
	return handler.service.SearchDiagnosticCentre(context)
}

// UpdateDiagnosticCentre godoc
// @Summary Update a diagnostic centre
// @Description Update an existing diagnostic centre's information (owner only)
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID (UUID format)" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param RequestBody body domain.UpdateDiagnosticBodyDTO true "Updated diagnostic centre details"
// @Success 200 {object} handlers.DiagnosticCentreSwagger "Diagnostic centre updated successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id} [put]
func (handler *HTTPHandler) UpdateDiagnosticCentre(context echo.Context) error {
	return handler.service.UpdateDiagnosticCentre(context)
}

// DeleteDiagnosticCentre godoc
// @Summary Delete a diagnostic centre
// @Description Delete an existing diagnostic centre (owner only)
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Success 200 {object} handlers.DiagnosticCentreSwagger "Diagnostic centre deleted successfully"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id} [delete]
func (handler *HTTPHandler) DeleteDiagnosticCentre(context echo.Context) error {
	return handler.service.DeleteDiagnosticCentre(context)
}

// GetDiagnosticCentresByOwner godoc
// @Summary List owner's diagnostic centres
// @Description Get all diagnostic centres owned by the authenticated user
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} handlers.DiagnosticCentreSwagger "List of diagnostic centres"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/owner [get]
func (handler *HTTPHandler) GetDiagnosticCentresByOwner(context echo.Context) error {
	return handler.service.GetDiagnosticCentresByOwner(context)
}

// GetDiagnosticCentreStats godoc
// @Summary Get diagnostic centre statistics
// @Description Get statistical information about a diagnostic centre (appointments, tests, etc.)
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Success 200 {object} handlers.DiagnosticCentreSwagger "Centre statistics"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/stats [get]
func (handler *HTTPHandler) GetDiagnosticCentreStats(context echo.Context) error {
	return handler.service.GetDiagnosticCentreStats(context)
}

// GetDiagnosticCentresByManager godoc
// @Summary List manager's diagnostic centres
// @Description Get all diagnostic centres managed by the authenticated diagnostic centre manager
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param assigned query boolean false "Assigned" 
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(50) default(20)
// @Success 200 {array} handlers.DiagnosticCentreSwagger "List of diagnostic centres"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/manager [get]
func (handler *HTTPHandler) GetDiagnosticCentresByManager(context echo.Context) error {
	return handler.service.GetDiagnosticCentresByManager(context)
}

// UpdateDiagnosticCentreManager godoc
// @Summary Update diagnostic centre manager
// @Description Update or assign a new manager to a diagnostic centre (owner only)
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param RequestBody body domain.UpdateDiagnosticManagerDTO true "Manager details"
// @Success 200 {object} handlers.DiagnosticCentreSwagger "Manager updated successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/manager [put]
func (handler *HTTPHandler) UpdateDiagnosticCentreManager(context echo.Context) error {
	return handler.service.UpdateDiagnosticCentreManager(context)
}

// GetDiagnosticCentreSchedules godoc
// @Summary Get diagnostic centre schedules
// @Description Get all schedules for a diagnostic centre with pagination and filtering
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)" format(date)
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)" format(date)
// @Param status query string false "Filter by status" Enums(PENDING,ACCEPTED,REJECTED)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} handlers.DiagnosticCentreSwagger "List of schedules"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/schedules [get]
func (handler *HTTPHandler) GetDiagnosticCentreSchedules(context echo.Context) error {
	return handler.service.GetDiagnosticCentreSchedules(context)
}

// GetDiagnosticCentreRecords godoc
// @Summary Get diagnostic centre medical records
// @Description Get all medical records uploaded by a diagnostic centre
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" default(123e4567-e89b-12d3-a456-426614174000)
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)" format(date)
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)" format(date)
// @Param document_type query string false "Filter by document type" Enums(LAB_REPORT,PRESCRIPTION,IMAGING,DISCHARGE_SUMMARY,OTHER)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} handlers.MedicalRecordSwagger "List of medical records"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/records [get]
func (handler *HTTPHandler) GetDiagnosticCentreRecords(context echo.Context) error {
	return handler.service.GetDiagnosticCentreRecords(context)
}

func (handler *HTTPHandler) GetDiagnosticCentreManagers(context echo.Context) error {
	return nil
}

// ListManagersByAdmin godoc
// @Summary List managers by admin
// @Tags DiagnosticCentre
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param assigned query boolean false "Assigned" 
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(50) default(20)
// @Success 200 {array} handlers.ManagerSwagger "List of managers"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/managers [get]
func (handler *HTTPHandler) ListManagersByAdmin(context echo.Context) error {
	return handler.service.ListManagers(context)
}

// AssignManagerByOwner godoc
// @Summary Assign manager to a centre
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre body domain.UpdateDiagnosticManagerDTO true "Assign Manager to a diagnostic centre"
// @Success 200 {object} handlers.ManagerSwagger "Manager assigned successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/assign [post]
func (handler *HTTPHandler) AssignManager(context echo.Context) error {
	return handler.service.AssignAdmin(context)
}

// UnAssignManagerByOwner godoc
// @Summary Unassign manager to a centre
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param RequestBody body domain.UnAssignDiagnosticManagerDTO true "Un Assign Manager to a diagnostic centre"
// @Success 200 {object} handlers.ManagerSwagger "Manager unassigned successfully"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/unassign [post]
func (handler *HTTPHandler) UnAssignManager(context echo.Context) error {
	return handler.service.UnAssignAdmin(context)
}

func (handler *HTTPHandler) SubmitKYC(context echo.Context) error {
	return handler.service.OwnerKYC(context)
}
