package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/adapters/ex/templates/emails"
	"github.com/diagnoxix/adapters/metrics"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
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
	params, err := buildCreateDiagnosticCentreParams(dto)
	if err != nil {
		utils.Error("Failed to build diagnostic centre params",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid diagnostic centre data")
	}

	params.CreatedBy = currentUser.UserID.String()

	// Start transaction
	tx, err := service.diagnosticPort.BeginDiagnostic(context.Request().Context())
	if err != nil {
		utils.Error("Failed to start diagnostic transaction",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	defer tx.Rollback(context.Request().Context())

	// Make this a transaction
	// Create diagnostic centre
	diagnostic_centre, err := tx.CreateDiagnosticCentre(
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
	// Implement test Price here
	buildPrice, err := buildTestPrice(dto, diagnostic_centre.ID)
	if err != nil {
		utils.Error("Failed to build test prices",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	test_price, err := tx.CreateTestPrice(context.Request().Context(), *buildPrice)
	if err != nil {
		utils.Error("Failed to submit test price",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "admin_id", Value: dto.AdminId.String()})
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	response, _ := utils.MarshalJSONField(test_price)

	centreRow := &db.List_Diagnostic_Centres_ByOwnerRow{
		ID:                   diagnostic_centre.ID,
		DiagnosticCentreName: diagnostic_centre.DiagnosticCentreName,
		Latitude:             diagnostic_centre.Latitude,
		Longitude:            diagnostic_centre.Longitude,
		Address:              diagnostic_centre.Address,
		Contact:              diagnostic_centre.Contact,
		Doctors:              diagnostic_centre.Doctors,
		AvailableTests:       diagnostic_centre.AvailableTests,
		CreatedAt:            diagnostic_centre.CreatedAt,
		UpdatedAt:            diagnostic_centre.UpdatedAt,
		TestPrices:           response,
	}
	res, err := buildDiagnosticCentreResponseFromRow(centreRow)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Commit transaction
	if err := tx.Commit(context.Request().Context()); err != nil {
		utils.Error("Failed to commit diagnostic transaction",
			utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, errors.New("failed to commit transaction"), context)
	}

	return utils.ResponseMessage(http.StatusCreated, res, context)
}

func (service *ServicesHandler) GetDiagnosticCentre(context echo.Context) error {
	// Get validated params
	params, ok := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid query parameters"), context)
	}

	ctx := context.Request().Context()
	diagnosticCentre, err := service.diagnosticPort.GetDiagnosticCentre(ctx, params.DiagnosticCentreID)

	if err != nil {
		utils.Error("Failed to get diagnostic centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: params.DiagnosticCentreID},
		)

		switch {
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusNotFound, "Diagnostic centre not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve diagnostic centre")
		}
	}

	// Convert response to common format
	centreRow := &db.List_Diagnostic_Centres_ByOwnerRow{
		ID:                   diagnosticCentre.ID,
		DiagnosticCentreName: diagnosticCentre.DiagnosticCentreName,
		Latitude:             diagnosticCentre.Latitude,
		Longitude:            diagnosticCentre.Longitude,
		Address:              diagnosticCentre.Address,
		Contact:              diagnosticCentre.Contact,
		Doctors:              diagnosticCentre.Doctors,
		AvailableTests:       diagnosticCentre.AvailableTests,
		CreatedAt:            diagnosticCentre.CreatedAt,
		UpdatedAt:            diagnosticCentre.UpdatedAt,
		Availability:         diagnosticCentre.Availability,
		TestPrices:           diagnosticCentre.TestPrices,
		AdminID:              diagnosticCentre.AdminID,
		AdminAssignedAt:      diagnosticCentre.AdminAssignedAt,
		AdminAssignedBy:      diagnosticCentre.AdminAssignedBy,
		AdminStatus:          diagnosticCentre.AdminStatus,
	}

	// Build response
	res, err := buildDiagnosticCentreResponseFromRow(centreRow)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, res, context)
}

func (service *ServicesHandler) GetDiagnosticCentreByManagerOrOwner(context echo.Context) error {
	// Authentication & Authorization
	user, err := PrivateMiddlewareContext(context, []db.UserEnum{
		db.UserEnumDIAGNOSTICCENTREMANAGER,
		db.UserEnumDIAGNOSTICCENTREOWNER,
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get validated params
	params, ok := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid query parameters"), context)
	}

	ctx := context.Request().Context()
	var diagnosticCentre interface{}

	// Handle different user types
	switch user.UserType {
	case db.UserEnumDIAGNOSTICCENTREMANAGER:
		diagnosticCentre, err = service.diagnosticPort.GetDiagnosticCentreByManager(ctx, db.Get_Diagnostic_Centre_ByManagerParams{
			ID:      params.DiagnosticCentreID,
			AdminID: pgtype.UUID{Bytes: user.UserID, Valid: true},
		})
	case db.UserEnumDIAGNOSTICCENTREOWNER:
		diagnosticCentre, err = service.diagnosticPort.GetDiagnosticCentreByOwner(ctx, db.Get_Diagnostic_Centre_ByOwnerParams{
			ID:        params.DiagnosticCentreID,
			CreatedBy: user.UserID.String(),
		})
	default:
		return utils.ErrorResponse(
			http.StatusForbidden,
			errors.New("user type not authorized to perform this action"),
			context,
		)
	}

	// Handle errors
	if err != nil {
		utils.Error("Failed to get diagnostic centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: user.UserID.String()},
			utils.LogField{Key: "diagnostic_centre_id", Value: params.DiagnosticCentreID},
			utils.LogField{Key: "user_type", Value: user.UserType},
		)

		switch {
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Diagnostic centre not found")
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusNotFound, "Diagnostic centre not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve diagnostic centre")
		}
	}

	// Convert response to common format
	var centreRow *db.List_Diagnostic_Centres_ByOwnerRow

	switch v := diagnosticCentre.(type) {
	case *db.Get_Diagnostic_CentreRow:
		centreRow = &db.List_Diagnostic_Centres_ByOwnerRow{
			ID:                   v.ID,
			DiagnosticCentreName: v.DiagnosticCentreName,
			Latitude:             v.Latitude,
			Longitude:            v.Longitude,
			Address:              v.Address,
			Contact:              v.Contact,
			Doctors:              v.Doctors,
			AvailableTests:       v.AvailableTests,
			CreatedAt:            v.CreatedAt,
			UpdatedAt:            v.UpdatedAt,
			Availability:         v.Availability,
			TestPrices:           v.TestPrices,
			AdminID:              v.AdminID,
			AdminAssignedAt:      v.AdminAssignedAt,
			AdminAssignedBy:      v.AdminAssignedBy,
			AdminStatus:          v.AdminStatus,
		}
	case *db.DiagnosticCentre:
		centreRow = &db.List_Diagnostic_Centres_ByOwnerRow{
			ID:                   v.ID,
			DiagnosticCentreName: v.DiagnosticCentreName,
			Latitude:             v.Latitude,
			Longitude:            v.Longitude,
			Address:              v.Address,
			Contact:              v.Contact,
			Doctors:              v.Doctors,
			AvailableTests:       v.AvailableTests,
			CreatedAt:            v.CreatedAt,
			UpdatedAt:            v.UpdatedAt,
			AdminID:              v.AdminID,
			AdminAssignedAt:      v.AdminAssignedAt,
			AdminAssignedBy:      v.AdminAssignedBy,
			AdminStatus:          v.AdminStatus,
		}
	default:
		return utils.ErrorResponse(
			http.StatusInternalServerError,
			errors.New("invalid diagnostic centre data type"),
			context,
		)
	}

	// Build response
	res, err := buildDiagnosticCentreResponseFromRow(centreRow)
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

	param, _ := SetDefaultPagination(&query.PaginationQueryDTO).(*domain.PaginationQueryDTO)

	// Build search parameters
	params := db.Get_Nearest_Diagnostic_CentresParams{
		Radians:   query.Latitude,
		Radians_2: query.Longitude,
		Column4:   query.DayOfWeek,
		Column5:   query.Test, // Empty string will match all days in SQL query
		Limit:     param.GetLimit(),
		Offset:    param.GetOffset(),
	}

	hasFilters := false
	// // Add optional filters with case normalization
	if query.Doctor != "" {
		params.Doctors = []string{query.Doctor}
		hasFilters = true
	}

	response, err := service.diagnosticPort.GetNearestDiagnosticCentres(context.Request().Context(), params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	// Record metrics for the search
	metrics.RecordDiagnosticSearch(hasFilters, len(response))

	result := make([]map[string]interface{}, 0, len(response))
	for _, v := range response {
		// Map v to a DiagnosticCentre struct
		diagnosticCentre := buildDiagnosticCentre(*v)
		item, err := buildDiagnosticCentreResponseFromRow(diagnosticCentre)
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
	param_id := context.Param(utils.DiagnosticCentreID)
	dto, err := buildUpdateDiagnosticCentreByOwnerParams(body)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	dto.ID = param_id
	dto.CreatedBy = currentUser.UserID.String()
	if body.ADMINID != uuid.Nil {
		dto.AdminAssignedBy = pgtype.UUID{Bytes: currentUser.UserID, Valid: true}
	}
	response, err := service.diagnosticPort.UpdateDiagnosticCentreByOwner(context.Request().Context(), *dto)
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
	_, err = service.diagnosticPort.GetDiagnosticCentreByOwner(context.Request().Context(), db.Get_Diagnostic_Centre_ByOwnerParams{
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
	response, err := service.diagnosticPort.DeleteDiagnosticCentreByOwner(context.Request().Context(), params)
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

	response, err := service.diagnosticPort.ListDiagnosticCentresByOwner(context.Request().Context(), dbParams)
	if err != nil {
		utils.Error("Failed to list diagnostic centres",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	result := make([]map[string]interface{}, 0, len(response))
	for _, centre := range response {
		// Convert DiagnosticCentre to Get_Nearest_Diagnostic_CentresRow
		centreRow := &db.List_Diagnostic_Centres_ByOwnerRow{
			ID:                   centre.ID,
			DiagnosticCentreName: centre.DiagnosticCentreName,
			Latitude:             centre.Latitude,
			Longitude:            centre.Longitude,
			Address:              centre.Address,
			Contact:              centre.Contact,
			Doctors:              centre.Doctors,
			AvailableTests:       centre.AvailableTests,
			TestPrices:           centre.TestPrices,
			AdminID:              centre.AdminID,
			AdminAssignedAt:      centre.AdminAssignedAt,
			AdminAssignedBy:      centre.AdminAssignedBy,
			AdminUnassignedAt:    centre.AdminUnassignedAt,
			AdminStatus:          centre.AdminStatus,
			CreatedAt:            centre.CreatedAt,
			UpdatedAt:            centre.UpdatedAt,
		}
		item, err := buildDiagnosticCentreResponseFromRow(centreRow)
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
	admin, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// Get diagnostic centre ID
	centre_id := context.Param(utils.DiagnosticCentreID)

	centre, err := service.diagnosticPort.GetDiagnosticCentreByManager(context.Request().Context(), db.Get_Diagnostic_Centre_ByManagerParams{
		ID:      centre_id,
		AdminID: pgtype.UUID{Bytes: admin.UserID, Valid: true},
	})
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

	// Get diagnostic centres
	centre, err := service.diagnosticPort.GetDiagnosticCentreByManager(
		context.Request().Context(),
		db.Get_Diagnostic_Centre_ByManagerParams{
			ID:      currentUser.UserID.String(), // Will be filled by DB query
			AdminID: pgtype.UUID{Bytes: currentUser.UserID, Valid: true},
		},
	)
	if err != nil {
		utils.Error("Failed to get diagnostic centres by manager",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: currentUser.UserID.String()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	// Convert to Get_Nearest_Diagnostic_CentresRow
	centreRow := &db.List_Diagnostic_Centres_ByOwnerRow{
		ID:                   centre.ID,
		DiagnosticCentreName: centre.DiagnosticCentreName,
		Latitude:             centre.Latitude,
		Longitude:            centre.Longitude,
		Address:              centre.Address,
		Contact:              centre.Contact,
		Doctors:              centre.Doctors,
		AvailableTests:       centre.AvailableTests,
		TestPrices:           centre.TestPrices,
		CreatedAt:            centre.CreatedAt,
		UpdatedAt:            centre.UpdatedAt,
	}

	// Build response
	result, err := buildDiagnosticCentreResponseFromRow(centreRow)
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
	_, err = service.diagnosticPort.GetDiagnosticCentreByOwner(context.Request().Context(), db.Get_Diagnostic_Centre_ByOwnerParams{
		ID:        centreID,
		CreatedBy: currentUser.UserID.String(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found or not owned by user"), context)
	}

	updateParams := db.Update_Diagnostic_Centre_ByOwnerParams{
		ID:        centreID,
		CreatedBy: currentUser.UserID.String(),
		AdminID:   pgtype.UUID{Bytes: managerDetails.ManagerID, Valid: true},
	}

	response, err := service.diagnosticPort.UpdateDiagnosticCentreByOwner(context.Request().Context(), updateParams)
	if err != nil {
		utils.Error("Failed to update diagnostic centre manager",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: centreID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	centreRow := &db.List_Diagnostic_Centres_ByOwnerRow{
		ID:                   response.ID,
		DiagnosticCentreName: response.DiagnosticCentreName,
		Latitude:             response.Latitude,
		Longitude:            response.Longitude,
		Address:              response.Address,
		Contact:              response.Contact,
		Doctors:              response.Doctors,
		AvailableTests:       response.AvailableTests,
		CreatedAt:            response.CreatedAt,
		UpdatedAt:            response.UpdatedAt,
	}
	result, err := buildDiagnosticCentreResponseFromRow(centreRow)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	return utils.ResponseMessage(http.StatusOK, result, context)
}

// GetDiagnosticCentreSchedules retrieves all schedules for a diagnostic centre
func (service *ServicesHandler) GetDiagnosticCentreSchedules(context echo.Context) error {
	// Authentication & Authorization check - try owner first, then manager
	admin, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
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
	if _, err := service.diagnosticPort.GetDiagnosticCentreByManager(context.Request().Context(), db.Get_Diagnostic_Centre_ByManagerParams{ID: centreID, AdminID: pgtype.UUID{
		Bytes: admin.UserID,
		Valid: true}}); err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found"), context)
	}

	// Get schedules from schedule repository
	req := db.Get_Diagnsotic_Schedules_By_CentreParams{
		DiagnosticCentreID: centreID,
		Offset:             params.GetOffset(),
		Limit:              params.GetLimit(),
	}

	schedules, err := service.schedulePort.GetDiagnosticSchedulesByCentre(context.Request().Context(), req)
	if err != nil {
		utils.Error("Failed to retrieve schedules",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: centreID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	if len(schedules) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, context)
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

	if _, err := service.diagnosticPort.GetDiagnosticCentreByOwner(context.Request().Context(), db.Get_Diagnostic_Centre_ByOwnerParams{ID: centreID, CreatedBy: currentUser.UserID.String()}); err != nil {
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

func (service *ServicesHandler) AssignAdmin(context echo.Context) error {
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	// Get diagnostic centre ID and manager details
	manager, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateDiagnosticManagerDTO)

	ctx := context.Request().Context()
	// Verify ownership
	diagnostic_centre, err := service.diagnosticPort.GetDiagnosticCentreByOwner(ctx, db.Get_Diagnostic_Centre_ByOwnerParams{
		ID:        manager.ID.String(),
		CreatedBy: currentUser.UserID.String(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found or not owned by owner"), context)
	}

	assignedAdmin, err := service.diagnosticPort.AssignAdmin(ctx, db.AssignAdminParams{
		ID:      manager.ID.String(),
		AdminID: pgtype.UUID{Bytes: manager.ManagerID, Valid: true}, AdminAssignedBy: pgtype.UUID{Bytes: currentUser.UserID, Valid: true},
	})
	if err != nil {
		utils.Error("Failed to assign manager to centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: manager.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	var address domain.Address
	if err := utils.UnmarshalJSONField(diagnostic_centre.Address, &address); err != nil {
		utils.Error("Error unmarshal diagnostic centre address",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: manager.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}

	var managerDetails domain.ManagerDetails
	if err := utils.UnmarshalJSONField(diagnostic_centre.ManagerDetails, &managerDetails); err != nil {
		utils.Error("Error unmarshal diagnostic manager details",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: manager.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	// Send Notification email
	add := fmt.Sprintf("%s %s %s %s", address.Street, address.City, address.State, address.Country)
	emailData := &emails.DiagnosticCentreManagement{
		Name:          managerDetails.FullName,
		CentreName:    diagnostic_centre.DiagnosticCentreName,
		CentreAddress: add,
	}
	go service.emailGoroutine(
		emailData,
		managerDetails.Email,
		emails.SubjectDiagnosticCentreManagement,
		emails.TemplateDiagnosticCentreManagement,
	)
	return utils.ResponseMessage(http.StatusOK, assignedAdmin, context)
}

func (service *ServicesHandler) UnAssignAdmin(context echo.Context) error {
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	// Get diagnostic centre ID and manager details
	centre, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UnAssignDiagnosticManagerDTO)

	ctx := context.Request().Context()
	// Verify ownership
	_, err = service.diagnosticPort.GetDiagnosticCentreByOwner(ctx, db.Get_Diagnostic_Centre_ByOwnerParams{
		ID:        centre.ID.String(),
		CreatedBy: currentUser.UserID.String(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, errors.New("diagnostic centre not found or not owned by user"), context)
	}

	assignedAdmin, err := service.diagnosticPort.UnAssignAdmin(ctx, db.UnassignAdminParams{
		ID:                centre.ID.String(),
		AdminUnassignedBy: pgtype.UUID{Bytes: currentUser.UserID, Valid: true},
	})
	if err != nil {
		utils.Error("Failed to unassign manager to centre",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: centre.ID})
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, assignedAdmin, context)
}

// buildDiagnosticCentre converts a database row to a domain diagnostic centre
func buildDiagnosticCentre(row db.Get_Nearest_Diagnostic_CentresRow) *db.List_Diagnostic_Centres_ByOwnerRow {
	return &db.List_Diagnostic_Centres_ByOwnerRow{
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
		Availability:         row.Availability,
		TestPrices:           row.TestPrices,
	}
}
