package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/adapters/metrics"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) CreateDiagnosticCentre(context echo.Context) error {
	// Authentication & Authorization
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
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
		case strings.Contains(err.Error(), "unique_admin_id"):
			return echo.NewHTTPError(http.StatusConflict, "Manager is already assigned to another diagnostic centre")
		case strings.Contains(err.Error(), "diagnostic_centres_latitude_longitude_key"):
			return echo.NewHTTPError(http.StatusConflict, "Latitude and Longitude already assigned to a centre")
		case strings.Contains(err.Error(), "diagnostic_centres_diagnostic_centre_name_created_by_longit_key"):
			return echo.NewHTTPError(http.StatusConflict, "Diagnostic centre name with the location already exists")
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, utils.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid diagnostic centre data")
		case errors.Is(err, utils.ErrConflictError):
			return echo.NewHTTPError(http.StatusConflict, utils.ErrConflictError.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create diagnostic centre")
		}
	}
	fmt.Println(response, "okkk")
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
		params.Doctors = []db.Doctor{db.Doctor(query.Doctor)}
		hasFilters = true
	}
	if query.Test != "" {
		params.AvailableTests = []db.AvailableTests{db.AvailableTests(query.Test)}
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
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
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

// DeleteDiagnosticCentre deletes a diagnostic center (owner only)
func (service *ServicesHandler) DeleteDiagnosticCentre(context echo.Context) error {
	// Authentication & Authorization check for owner
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		if errors.Is(err, utils.ErrInvalidToken) || errors.Is(err, ErrUnauthorized) {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: utils.AuthenticationRequired,
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	param := context.Param(utils.DiagnosticCentreID)
	parsedUUID, err := uuid.Parse(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid diagnostic centre ID")
	}
	params := db.Delete_Diagnostic_Centre_ByOwnerParams{
		ID:        param,
		CreatedBy: currentUser.UserID.String(),
	}

	// First check if the diagnostic center exists and is owned by this user
	_, err = service.DiagnosticRepo.GetDiagnosticCentreByOwner(context.Request().Context(), db.Get_Diagnostic_Centre_ByOwnerParams{
		ID:        parsedUUID.String(),
		CreatedBy: currentUser.UserID.String(),
	})
	if err != nil {
		utils.Error("Failed to find diagnostic centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()},
			utils.LogField{Key: "diagnostic_centre_id", Value: param})

		if errors.Is(err, utils.ErrNotFound) {
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "diagnostic centre not found",
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	// Attempt to delete
	response, err := service.DiagnosticRepo.DeleteDiagnosticCentreByOwner(context.Request().Context(), params)
	if err != nil {
		utils.Error("Failed to delete diagnostic centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()},
			utils.LogField{Key: "diagnostic_centre_id", Value: param})
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete diagnostic centre",
		}
	}

	// Return success response
	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"message": "Diagnostic centre deleted successfully",
		"id":      response.ID,
	}, context)
}

// GetDiagnosticCentresByOwner retrieves all diagnostic centers owned by the authenticated user
func (service *ServicesHandler) GetDiagnosticCentresByOwner(context echo.Context) error {
	// Authentication & Authorization check for owner
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get and validate pagination parameters
	params, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.PaginationQueryDTO)
	params = SetDefaultPagination(params).(*domain.PaginationQueryDTO)

	dbParams := db.List_Diagnostic_Centres_ByOwnerParams{
		CreatedBy: currentUser.UserID.String(),
		Limit:     params.GetLimit(),
		Offset:    params.GetOffset(),
	}

	response, err := service.DiagnosticRepo.ListDiagnosticCentresByOwner(context.Request().Context(), dbParams)
	if err != nil {
		utils.Error("Failed to list diagnostic centres",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	result := make([]map[string]interface{}, 0, len(response))
	for _, centre := range response {
		item, err := buildDiagnosticCentreResponseFromRow(centre, context)
		if err != nil {
			return utils.ErrorResponse(http.StatusInternalServerError, err, context)
		}
		result = append(result, item)
	}

	return utils.ResponseMessage(http.StatusOK, result, context)
}

// GetDiagnosticCentreStats retrieves statistical information about a diagnostic centre
func (service *ServicesHandler) GetDiagnosticCentreStats(context echo.Context) error {
	// Authenticate and authorize user - first try owner, then manager
	_, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get diagnostic centre ID
	param := context.Param(utils.DiagnosticCentreID)

	centre, err := service.DiagnosticRepo.GetDiagnosticCentre(context.Request().Context(), param)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found"), context)
	}

	// TODO: Get statistics from schedule and record repositories
	// For now, return basic info
	stats := map[string]interface{}{
		"diagnostic_centre_id":   centre.ID,
		"diagnostic_centre_name": centre.DiagnosticCentreName,
		"total_doctors":          len(centre.Doctors),
		"total_tests":            len(centre.AvailableTests),
		// TODO: Add more statistics from schedules and records
	}

	return utils.ResponseMessage(http.StatusOK, stats, context)
}

// GetDiagnosticCentresByManager retrieves all diagnostic centres managed by the authenticated manager
func (service *ServicesHandler) GetDiagnosticCentresByManager(context echo.Context) error {
	// Authentication & Authorization check for manager
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get and validate pagination parameters
	// params, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.PaginationQueryDTO)
	// params = SetDefaultPagination(params).(*domain.PaginationQueryDTO)

	// Build query parameters
	dbParams := db.Get_Diagnostic_Centre_ByManagerParams{
		ID:      currentUser.UserID.String(), // Will be filled by DB query
		AdminID: currentUser.UserID.String(),
	}

	// Get diagnostic centres
	centre, err := service.DiagnosticRepo.GetDiagnosticCentreByManager(context.Request().Context(), dbParams)
	if err != nil {
		utils.Error("Failed to get diagnostic centres by manager",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Build response
	result, err := buildDiagnosticCentreResponseFromRow(centre, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, []map[string]interface{}{result}, context)
}

// UpdateDiagnosticCentreManager updates the manager of a diagnostic centre
func (service *ServicesHandler) UpdateDiagnosticCentreManager(context echo.Context) error {
	// Authentication & Authorization check for owner
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get diagnostic centre ID and manager details
	centreID := context.Param(utils.DiagnosticCentreID)
	managerDetails, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateDiagnosticManagerDTO)

	// Verify ownership
	_, err = service.DiagnosticRepo.GetDiagnosticCentreByOwner(context.Request().Context(), db.Get_Diagnostic_Centre_ByOwnerParams{
		ID:        centreID,
		CreatedBy: currentUser.UserID.String(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found or not owned by user"), context)
	}

	// Update manager
	updateParams := db.Update_Diagnostic_Centre_ByOwnerParams{
		ID:        centreID,
		CreatedBy: currentUser.UserID.String(),
		AdminID:   managerDetails.ManagerID,
	}

	response, err := service.DiagnosticRepo.UpdateDiagnosticCentreByOwner(context.Request().Context(), updateParams)
	if err != nil {
		utils.Error("Failed to update diagnostic centre manager",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: centreID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	result, err := buildDiagnosticCentreResponseFromRow(response, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, result, context)
}

// GetDiagnosticCentreSchedules retrieves all schedules for a diagnostic centre
func (service *ServicesHandler) GetDiagnosticCentreSchedules(context echo.Context) error {
	// Authentication & Authorization check - try owner first, then manager
	_, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get diagnostic centre ID and schedule params
	centreID := context.Param(utils.DiagnosticCentreID)
	params := &domain.GetDiagnosticSchedulesByCentreParamDTO{
		DiagnosticCentreID: uuid.Must(uuid.Parse(centreID)),
	}
	params = SetDefaultPagination(params).(*domain.GetDiagnosticSchedulesByCentreParamDTO)

	// Verify centre exists
	if _, err := service.DiagnosticRepo.GetDiagnosticCentre(context.Request().Context(), centreID); err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found"), context)
	}

	// Get schedules from schedule repository
	req := db.Get_Diagnsotic_Schedules_By_CentreParams{
		DiagnosticCentreID: centreID,
		Offset:             params.GetOffset(),
		Limit:              params.GetLimit(),
	}

	schedules, err := service.ScheduleRepo.GetDiagnosticSchedulesByCentre(context.Request().Context(), req)
	if err != nil {
		utils.Error("Failed to retrieve schedules",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: centreID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, schedules, context)
}

// GetDiagnosticCentreRecords retrieves medical records for a diagnostic centre
func (service *ServicesHandler) GetDiagnosticCentreRecords(context echo.Context) error {
	// Authentication & Authorization check - try owner first, then manager
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		currentUser, err = PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
		if err != nil {
			return utils.ErrorResponse(http.StatusUnauthorized, err, context)
		}
	}

	_ = currentUser // We don't need the user object right now

	// Get diagnostic centre ID and filter params
	centreID := context.Param(utils.DiagnosticCentreID)
	params := &domain.GetDiagnosticRecordsParamDTO{
		DiagnosticCentreID: uuid.Must(uuid.Parse(centreID)).String(),
	}
	params = SetDefaultPagination(params).(*domain.GetDiagnosticRecordsParamDTO)

	if _, err := service.DiagnosticRepo.GetDiagnosticCentre(context.Request().Context(), centreID); err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found"), context)
	}

	// TODO: Implement medical records repository methods and querying
	// For now, return empty result with proper format
	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"records":  []interface{}{},
		"total":    0,
		"page":     params.Page,
		"per_page": params.PerPage,
	}, context)
}

// buildDiagnosticCentre converts a database row to a domain diagnostic centre
func buildDiagnosticCentre(row db.Get_Nearest_Diagnostic_CentresRow) *db.DiagnosticCentre {
	return &db.DiagnosticCentre{
		ID:                   row.ID,
		DiagnosticCentreName: row.DiagnosticCentreName,
		Latitude:             row.Latitude,
		Longitude:            row.Longitude,
		Address:              row.Address,
		Contact:              row.Contact,
		Doctors:              row.Doctors,
		AvailableTests:       row.AvailableTests,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
}
