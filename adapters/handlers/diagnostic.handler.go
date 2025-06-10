package handlers

import "github.com/labstack/echo/v4"

// CreateDiagnostic godoc
// @Summary Create a new diagnostic centre
// @Tags DiagnosticCentre
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /diagnostic_centres [post]
func (handler *HTTPHandler) CreateDiagnostic(context echo.Context) error {
	return handler.service.CreateDiagnosticCentre(context)
}

// GetDiagnosticCentre godoc
// @Summary Get a diagnostic centre by ID
// @Tags DiagnosticCentre
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /diagnostic_centres/{diagnostic_centre_id} [get]
func (handler *HTTPHandler) GetDiagnosticCentre(context echo.Context) error {
	return handler.service.GetDiagnosticCentre(context)
}

// SearchDiagnosticCentre godoc
// @Summary Search for diagnostic centres
// @Tags DiagnosticCentre
// @Produce json
// @Param latitude query number true "Latitude"
// @Param longitude query number true "Longitude"
// @Param doctor query string false "Doctor"
// @Param available_tests query string false "Available Test"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /diagnostic_centres [get]
func (handler *HTTPHandler) SearchDiagnosticCentre(context echo.Context) error {
	return handler.service.SearchDiagnosticCentre(context)
}

func (handler *HTTPHandler) UpdateDiagnosticCentre(context echo.Context) error {
	return handler.service.UpdateDiagnosticCentre(context)
}
