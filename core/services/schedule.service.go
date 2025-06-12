package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) CreateSchedule(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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
	response, err := service.scheduleRepo.CreateDiagnosticSchedule(context.Request().Context(), params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, response, context)
}

func (service *ServicesHandler) GetDiagnosticSchedule(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	param, _ := context.Get(
		utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticScheduleParamDTO)
	req := db.Get_Diagnostic_ScheduleParams{
		ID:     param.ScheduleID.String(),
		UserID: currentUser.UserID.String(),
	}
	response, err := service.scheduleRepo.GetDiagnosticSchedule(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) GetDiagnosticSchedules(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticSchedulesQueryDTO)
	var Limit, Offset int32
	if query.Limit == 0 {
		Limit = utils.Limit
	}
	if query.Offset == 0 {
		Offset = utils.Offset
	}
	req := db.Get_Diagnostic_SchedulesParams{
		UserID: currentUser.UserID.String(),
		Limit:  Limit,
		Offset: Offset,
	}
	response, err := service.scheduleRepo.GetDiagnosticSchedules(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) UpdateDiagnosticSchedule(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumUSER))
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
	response, err := service.scheduleRepo.UpdateDiagnosticSchedule(context.Request().Context(), body)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusNoContent, response, context)
}

func (service *ServicesHandler) DeleteDiagnosticSchedule(context echo.Context) error {
	return nil
}

func (service *ServicesHandler) GetDiagnosticScheduleByCentre(context echo.Context) error {
	_, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	param, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticScheduleByCentreParamDTO)
	req := db.Get_Diagnsotic_Schedule_By_CentreParams{
		ID:                 param.ScheduleID.String(),
		DiagnosticCentreID: param.DiagnosticCentreID.String(),
	}
	response, err := service.scheduleRepo.GetDiagnosticScheduleByCentre(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) GetDiagnosticSchedulesByCentre(context echo.Context) error {
	_, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	param, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetDiagnosticSchedulesByCentreParamDTO)
	var Limit, Offset int32
	if param.Limit == 0 {
		Limit = utils.Limit
	}
	if param.Offset == 0 {
		Offset = utils.Offset
	}
	req := db.Get_Diagnsotic_Schedules_By_CentreParams{
		DiagnosticCentreID: param.DiagnosticCentreID.String(),
		Offset:             Offset,
		Limit:              Limit,
	}
	response, err := service.scheduleRepo.GetDiagnosticSchedulesByCentre(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) UpdateDiagnosticScheduleByCentre(context echo.Context) error {
	_, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	body, _ := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateDiagnosticScheduleByCentreDTO)
	schedule_id := context.Param(utils.ScheduleID)
	diagnostic_centre_id := context.Param(utils.DiagnosticCentreID)
	req := db.Update_Diagnostic_Schedule_By_CentreParams{
		DiagnosticCentreID: diagnostic_centre_id,
		AcceptanceStatus:   body.AcceptanceStatus,
		ID:                 schedule_id,
	}
	response, err := service.scheduleRepo.UpdateDiagnosticScheduleByCentre(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}
