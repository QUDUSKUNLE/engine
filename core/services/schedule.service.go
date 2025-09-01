package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (service *ServicesHandler) CreateSchedule(context echo.Context) error {
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.CreateScheduleDTO)
	// Convert dto.ScheduleTime (string) to pgtype.Timestamptz
	var scheduleTime pgtype.Timestamptz
	err = scheduleTime.Scan(dto.ScheduleTime)
	if err != nil {
		// Try parsing as RFC3339 with or without milliseconds
		parsed, parseErr := time.Parse(time.RFC3339Nano, dto.ScheduleTime)
		if parseErr != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.ErrScheduleTimeFormatRFC3339), context)
		}
		err = scheduleTime.Scan(parsed)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.ErrScheduleTimeFormatParsing), context)
		}
	}
	params := db.Create_Diagnostic_ScheduleParams{
		UserID:             currentUser.UserID.String(),
		DiagnosticCentreID: dto.DiagnosticCentreID.String(),
		ScheduleTime:       scheduleTime,
		TestType:           dto.TestType,
		AcceptanceStatus:   db.ScheduleAcceptanceStatusPENDING,
		Doctor:             string(dto.Doctor),
		Notes:              pgtype.Text{String: dto.Notes},
	}
	response, err := service.schedulePort.CreateDiagnosticSchedule(context.Request().Context(), params)
	if err != nil {
		utils.Error("Failed to create diagnostic schedule",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: params.UserID},
			utils.LogField{Key: "diagnostic_centre_id", Value: params.DiagnosticCentreID})
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, response, context)
}

func (service *ServicesHandler) GetDiagnosticSchedule(context echo.Context) error {
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	param, _ := context.Get(
		utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticScheduleParamDTO)
	req := db.Get_Diagnostic_ScheduleParams{
		ID:     param.ScheduleID.String(),
		UserID: currentUser.UserID.String(),
	}
	response, err := service.schedulePort.GetDiagnosticSchedule(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) GetDiagnosticSchedules(context echo.Context) error {
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.PaginationQueryDTO)

	query = SetDefaultPagination(query).(*domain.PaginationQueryDTO)

	req := db.Get_Diagnostic_SchedulesParams{
		UserID: currentUser.UserID.String(),
		Limit:  query.GetLimit(),
		Offset: query.GetOffset(),
	}
	response, err := service.schedulePort.GetDiagnosticSchedules(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, context)
	}

	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) UpdateDiagnosticSchedule(context echo.Context) error {
	currentUser, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateScheduleDTO)
	var scheduleTime pgtype.Timestamptz
	err = scheduleTime.Scan(dto.ScheduleTime)
	if err != nil {
		// Try parsing as RFC3339 with or without milliseconds
		parsed, parseErr := time.Parse(time.RFC3339Nano, dto.ScheduleTime)
		if parseErr != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.ErrScheduleTimeFormatRFC3339), context)
		}
		err = scheduleTime.Scan(parsed)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.ErrScheduleTimeFormatParsing), context)
		}
	}
	schedule_id := context.Param(utils.ScheduleID)
	body := db.Update_Diagnostic_ScheduleParams{
		ID:             schedule_id,
		ScheduleTime:   scheduleTime,
		TestType:       dto.TestType,
		ScheduleStatus: dto.ScheduleStatus,
		Notes:          pgtype.Text{String: dto.Notes},
		Doctor:         string(dto.Doctor),
		UserID:         currentUser.UserID.String(),
	}
	response, err := service.schedulePort.UpdateDiagnosticSchedule(context.Request().Context(), body)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusNoContent, response, context)
}

func (service *ServicesHandler) DeleteDiagnosticSchedule(context echo.Context) error {
	return nil
}

func (service *ServicesHandler) GetDiagnosticScheduleByCentre(context echo.Context) error {
	_, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	param, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticScheduleByCentreParamDTO)
	req := db.Get_Diagnsotic_Schedule_By_CentreParams{
		ID:                 param.ScheduleID.String(),
		DiagnosticCentreID: param.DiagnosticCentreID.String(),
	}
	response, err := service.schedulePort.GetDiagnosticScheduleByCentre(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) GetDiagnosticSchedulesByCentre(context echo.Context) error {
	// Authentication check
	_, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	// Get and validate query parameters
	param, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticSchedulesByCentreParamDTO)

	param = SetDefaultPagination(param).(*domain.GetDiagnosticSchedulesByCentreParamDTO)
	req := db.Get_Diagnsotic_Schedules_By_CentreParams{
		DiagnosticCentreID: param.DiagnosticCentreID.String(),
		Offset:             param.GetOffset(),
		Limit:              param.GetLimit(),
	}
	response, err := service.schedulePort.GetDiagnosticSchedulesByCentre(context.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "No schedules found")
		case errors.Is(err, utils.ErrDatabaseError):
			utils.Error("Database error retrieving schedules",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "diagnostic_centre_id", Value: req.DiagnosticCentreID})
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		default:
			utils.Error("Failed to retrieve schedules",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "diagnostic_centre_id", Value: req.DiagnosticCentreID})
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve schedules")
		}
	}

	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, context)
	}

	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) UpdateDiagnosticScheduleByCentre(context echo.Context) error {
	// Verify authentication and authorization
	_, err := PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, utils.AuthenticationRequired)
	}

	// Get and validate request body
	body, ok := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateDiagnosticScheduleByCentreDTO)
	if !ok || body == nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.InvalidRequestBody)
	}

	// Get and validate path parameters
	schedule_id := context.Param(utils.ScheduleID)
	diagnostic_centre_id := context.Param(utils.DiagnosticCentreID)
	if schedule_id == "" || diagnostic_centre_id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, utils.MissingParameters)
	}

	req := db.Update_Diagnostic_Schedule_By_CentreParams{
		DiagnosticCentreID: diagnostic_centre_id,
		AcceptanceStatus:   body.AcceptanceStatus,
		ID:                 schedule_id,
	}
	response, err := service.schedulePort.UpdateDiagnosticScheduleByCentre(context.Request().Context(), req)
	if err != nil {
		// Handle specific error cases
		switch {
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, utils.ScheduleNotFound)
		case errors.Is(err, utils.ErrPermissionDenied):
			return echo.NewHTTPError(http.StatusForbidden, utils.PermissionDenied)
		default:
			// Log unexpected errors for debugging
			utils.Error("Failed to update diagnostic schedule",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "diagnostic_centre_id", Value: diagnostic_centre_id},
				utils.LogField{Key: "schedule_id", Value: schedule_id})
			return echo.NewHTTPError(http.StatusInternalServerError, utils.FailedToUpdateSchedule)
		}
	}

	return utils.ResponseMessage(http.StatusOK, response, context)
}
