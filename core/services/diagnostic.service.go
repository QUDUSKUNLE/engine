package services

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/adapters/metrics"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) CreateDiagnosticCentre(context echo.Context) error {
	// Authentication & Authorization
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, utils.AuthenticationRequired)
	}

	// This validated at the middleware level
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.CreateDiagnosticDTO)

	// Build and validate parameters
	params, err := buildCreateDiagnosticCentreParams(context, dto)
	if err != nil {
		utils.Error("Failed to build diagnostic centre params",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid diagnostic centre data")
	}

	params.CreatedBy = currentUser.UserID.String()

	// Create diagnostic centre
	response, err := service.DiagnosticRepo.CreateDiagnosticCentre(
		context.Request().Context(), *params)
	if err != nil {
		utils.Error("Failed to create diagnostic centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()},
			utils.LogField{Key: "diagnostic_centre", Value: params})
		switch {
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, utils.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid diagnostic centre data")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create diagnostic centre")
		}
	}
	res, err := buildDiagnosticCentreResponseFromRow(response, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, res, context)
}

func (service *ServicesHandler) GetDiagnosticCentre(context echo.Context) error {
	// This validated at the middleware level
	params, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticParamDTO)

	// Get diagnostic centre
	response, err := service.DiagnosticRepo.GetDiagnosticCentre(context.Request().Context(), params.DiagnosticCentreID)
	if err != nil {
		utils.Error("Failed to get diagnostic centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: params.DiagnosticCentreID})
		switch {
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Diagnostic centre not found")
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve diagnostic centre")
		}
	}
	res, err := buildDiagnosticCentreResponseFromRow(response, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, res, context)
}

func (service *ServicesHandler) SearchDiagnosticCentre(context echo.Context) error {
	// Get and validate query parameters
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.SearchDiagnosticCentreQueryDTO)

	// Validate coordinates
	if !isValidLatitude(query.Latitude) || !isValidLongitude(query.Longitude) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid coordinates")
	}

	// Build search parameters
	params := db.Get_Nearest_Diagnostic_CentresParams{
		Radians:   query.Latitude,
		Radians_2: query.Longitude,
	}

	hasFilters := false
	// Add optional filters
	if query.Doctor != "" {
		params.Doctors = []string{query.Doctor}
		hasFilters = true
	}
	if query.Test != "" {
		params.AvailableTests = []string{query.Test}
		hasFilters = true
	}

	response, err := service.DiagnosticRepo.GetNearestDiagnosticCentres(context.Request().Context(), params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	// Record metrics for the search
	metrics.RecordDiagnosticSearch(hasFilters, len(response))

	result := make([]map[string]interface{}, 0, len(response))
	for _, v := range response {
		// Map v to a DiagnosticCentre struct
		diagnosticCentre := buildDiagnosticCentre(*v)
		item, err := buildDiagnosticCentreResponseFromRow(diagnosticCentre, context)
		if err != nil {
			return utils.ErrorResponse(http.StatusInternalServerError, err, context)
		}
		item["distance"] = v.DistanceKm
		item["distance_unit"] = "km"
		result = append(result, item)
	}
	return utils.ResponseMessage(http.StatusOK, result, context)
}

func (service *ServicesHandler) UpdateDiagnosticCentre(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	body, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateDiagnosticBodyDTO)
	param := context.Param(utils.DiagnosticCentreID)
	dto, err := buildUpdateDiagnosticCentreByOwnerParams(context, body)
	if err != nil {
		return err
	}
	dto.ID = param
	dto.CreatedBy = currentUser.UserID.String()
	response, err := service.DiagnosticRepo.UpdateDiagnosticCentreByOwner(context.Request().Context(), *dto)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotAcceptable, err, context)
	}
	return utils.ResponseMessage(http.StatusNoContent, response, context)
}

// buildDiagnosticCentre converts a database row to a domain diagnostic centre
func buildDiagnosticCentre(row db.Get_Nearest_Diagnostic_CentresRow) *db.DiagnosticCentre {
	return &db.DiagnosticCentre{
		ID: row.ID,
		// Name:        row.Name,
		Address:   row.Address,
		Latitude:  row.Latitude,
		Longitude: row.Longitude,
		// CreatedBy:   row.CreatedBy,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// isValidLatitude checks if the latitude is within valid range (-90 to 90)
func isValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

// isValidLongitude checks if the longitude is within valid range (-180 to 180)
func isValidLongitude(lon float64) bool {
	return lon >= -180 && lon <= 180
}
