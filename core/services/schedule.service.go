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
	dto, ok := context.Get("validatedBodyDTO").(*domain.CreateScheduleDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	// Convert dto.ScheduleTime (string) to pgtype.Timestamptz
	var scheduleTime pgtype.Timestamptz
	err = scheduleTime.Scan(dto.ScheduleTime)
	if err != nil {
		// Try parsing as RFC3339 with or without milliseconds
		parsed, parseErr := time.Parse(time.RFC3339Nano, dto.ScheduleTime)
		if parseErr != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid schedule time format, expected RFC3339 or RFC3339Nano"), context)
		}
		err = scheduleTime.Scan(parsed)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid schedule time format after parsing"), context)
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
	response, err := service.ScheduleRepo.CreateDiagnosticSchedule(context.Request().Context(), params)
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
	param, ok := context.Get("validatedQueryParamDTO").(*domain.GetDiagnosticScheduleParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	req := db.Get_Diagnostic_ScheduleParams{
		ID:     param.ScheduleID.String(),
		UserID: currentUser.UserID.String(),
	}
	response, err := service.ScheduleRepo.GetDiagnosticSchedule(context.Request().Context(), req)
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
	req := db.Get_Diagnostic_SchedulesParams{
		UserID: currentUser.UserID.String(),
		Limit:  50,
		Offset: 0,
	}
	response, err := service.ScheduleRepo.GetDiagnosticSchedules(context.Request().Context(), req)
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
	dto, ok := context.Get("validatedBodyDTO").(*domain.UpdateScheduleDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	var scheduleTime pgtype.Timestamptz
	err = scheduleTime.Scan(dto.ScheduleTime)
	if err != nil {
		// Try parsing as RFC3339 with or without milliseconds
		parsed, parseErr := time.Parse(time.RFC3339Nano, dto.ScheduleTime)
		if parseErr != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid schedule time format, expected RFC3339 or RFC3339Nano"), context)
		}
		err = scheduleTime.Scan(parsed)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, errors.New("invalid schedule time format after parsing"), context)
		}
	}
	schedule_id := context.Param("schedule_id")
	body := db.Update_Diagnostic_ScheduleParams{
		ID:             schedule_id,
		ScheduleTime:   scheduleTime,
		TestType:       dto.TestType,
		ScheduleStatus: dto.ScheduleStatus,
		Notes:          pgtype.Text{String: dto.Notes},
		Doctor:         string(dto.Doctor),
		UserID:         currentUser.UserID.String(),
	}
	response, err := service.ScheduleRepo.UpdateDiagnosticSchedule(context.Request().Context(), body)
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
	param, ok := context.Get("validatedQueryParamDTO").(*domain.GetDiagnosticScheduleByCentreParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	req := db.Get_Diagnsotic_Schedule_By_CentreParams{
		ID:                 param.ScheduleID.String(),
		DiagnosticCentreID: param.DiagnosticCentreID.String(),
	}
	response, err := service.ScheduleRepo.GetDiagnosticScheduleByCentre(context.Request().Context(), req)
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

	param, ok := context.Get("validatedQueryParamDTO").(*domain.GetDiagnosticSchedulesByCentreParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}

	req := db.Get_Diagnsotic_Schedules_By_CentreParams{
		DiagnosticCentreID: param.DiagnosticCentreID.String(),
		Offset:             0,
		Limit:              50,
	}
	response, err := service.ScheduleRepo.GetDiagnosticSchedulesByCentre(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

func (service *ServicesHandler) UpdateDiagnosticScheduleByCentre(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	body, ok := context.Get("validatedBodyDTO").(*domain.UpdateDiagnosticScheduleByCentreDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	schedule_id := context.Param("schedule_id")
	req := db.Update_Diagnostic_Schedule_By_CentreParams{
		DiagnosticCentreID: currentUser.DiagnosticID.String(),
		AcceptanceStatus:   body.AcceptanceStatus,
		ID:                 schedule_id,
	}
	response, err := service.ScheduleRepo.UpdateDiagnosticScheduleByCentre(context.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}
