package handlers

import "github.com/labstack/echo/v4"

// CreateDiagnostic godoc
// @Summary Create a new diagnostic centre
// @Description Create a new diagnostic centre with location, contact details, and available services
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre body domain.CreateDiagnosticDTO true "Diagnostic centre details"
// @Success 201 {object} domain.DiagnosticCentreResponse "Diagnostic centre created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required/invalid token"
// @Failure 403 {object} utils.ErrorResponse "User is not a diagnostic centre owner"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres [post]
func (handler *HTTPHandler) CreateDiagnostic(context echo.Context) error {
	return handler.service.CreateDiagnosticCentre(context)
}

// GetDiagnosticCentre godoc
// @Summary Get a diagnostic centre by ID
// @Description Retrieve detailed information about a specific diagnostic centre
// @Tags DiagnosticCentre
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID (UUID format)" format(uuid)
// @Success 200 {object} domain.DiagnosticCentreResponse "Diagnostic centre details retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid diagnostic centre ID format"
// @Failure 404 {object} utils.ErrorResponse "Diagnostic centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id} [get]
func (handler *HTTPHandler) GetDiagnosticCentre(context echo.Context) error {
	return handler.service.GetDiagnosticCentre(context)
}

// SearchDiagnosticCentre godoc
// @Summary Search for diagnostic centres
// @Description Search for diagnostic centres based on location, available doctors, and test types
// @Tags DiagnosticCentre
// @Produce json
// @Param latitude query number true "Latitude (-90 to 90)" minimum(-90) maximum(90)
// @Param longitude query number true "Longitude (-180 to 180)" minimum(-180) maximum(180)
// @Param doctor query string false "Filter by doctor specialization"
// @Param available_tests query string false "Filter by available test type"
// @Param page query integer false "Page number for pagination" default(1) minimum(1)
// @Param per_page query integer false "Number of results per page" default(10) minimum(1) maximum(100)
// @Success 200 {array} domain.DiagnosticCentreResponse "List of diagnostic centres"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres [get]
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID (UUID format)" format(uuid)
// @Param diagnostic_centre body domain.UpdateDiagnosticBodyDTO true "Updated diagnostic centre details"
// @Success 200 {object} domain.DiagnosticCentreResponse "Diagnostic centre updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {object} utils.ErrorResponse "Authentication required/invalid token"
// @Failure 403 {object} utils.ErrorResponse "User is not the owner of this diagnostic centre"
// @Failure 404 {object} utils.ErrorResponse "Diagnostic centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id} [put]
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Success 200 {object} utils.SuccessResponse "Diagnostic centre deleted successfully"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "User is not the owner"
// @Failure 404 {object} utils.ErrorResponse "Diagnostic centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id} [delete]
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
// @Success 200 {array} domain.DiagnosticCentreResponse "List of diagnostic centres"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/owner [get]
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Success 200 {object} domain.DiagnosticCentreStats "Centre statistics"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Access denied"
// @Failure 404 {object} utils.ErrorResponse "Diagnostic centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id}/stats [get]
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
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} domain.DiagnosticCentreResponse "List of diagnostic centres"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "User is not a manager"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/manager [get]
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param manager_details body domain.UpdateDiagnosticManagerDTO true "Manager details"
// @Success 200 {object} domain.DiagnosticCentreResponse "Manager updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Not authorized"
// @Failure 404 {object} utils.ErrorResponse "Centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id}/manager [put]
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)" format(date)
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)" format(date)
// @Param status query string false "Filter by status" Enums(PENDING,ACCEPTED,REJECTED)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} domain.ScheduleResponse "List of schedules"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Access denied"
// @Failure 404 {object} utils.ErrorResponse "Centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id}/schedules [get]
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)" format(date)
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)" format(date)
// @Param document_type query string false "Filter by document type" Enums(LAB_REPORT,PRESCRIPTION,IMAGING,DISCHARGE_SUMMARY,OTHER)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param per_page query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {array} db.MedicalRecord "List of medical records"
// @Failure 401 {object} utils.ErrorResponse "Authentication required"
// @Failure 403 {object} utils.ErrorResponse "Access denied"
// @Failure 404 {object} utils.ErrorResponse "Centre not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /diagnostic_centres/{diagnostic_centre_id}/records [get]
func (handler *HTTPHandler) GetDiagnosticCentreRecords(context echo.Context) error {
	return handler.service.GetDiagnosticCentreRecords(context)
}
