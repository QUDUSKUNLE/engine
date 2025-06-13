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
