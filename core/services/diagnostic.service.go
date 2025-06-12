package services

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) CreateDiagnosticCentre(context echo.Context) error {
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	dto, _ := context.Get(utils.ValidatedBodyDTO).(*domain.CreateDiagnosticDTO)
	params, err := buildCreateDiagnosticCentreParams(context, dto)
	if err != nil {
		return err
	}
	params.CreatedBy = currentUser.UserID.String()
	response, err := service.diagnosticRepo.CreateDiagnosticCentre(
		context.Request().Context(), *params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	res, err := buildDiagnosticCentreResponseFromRow(response, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, res, context)
}

func (service *ServicesHandler) GetDiagnosticCentre(context echo.Context) error {
	var params domain.GetDiagnosticParamDTO
	if err := utils.ValidateParams(context, &params); err != nil {
		fmt.Println(err.Error())
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	response, err := service.diagnosticRepo.GetDiagnosticCentre(context.Request().Context(), params.DiagnosticCentreID)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	res, err := buildDiagnosticCentreResponseFromRow(response, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, res, context)
}

func (service *ServicesHandler) SearchDiagnosticCentre(context echo.Context) error {
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.SearchDiagnosticCentreQueryDTO)
	params := db.Get_Nearest_Diagnostic_CentresParams{
		Radians:   query.Latitude,
		Radians_2: query.Longitude,
	}
	if query.Doctor != "" {
		params.Doctors = []string{query.Doctor}
	}
	if query.Test != "" {
		params.AvailableTests = []string{query.Test}
	}
	response, err := service.diagnosticRepo.GetNearestDiagnosticCentres(context.Request().Context(), params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
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
	response, err := service.diagnosticRepo.UpdateDiagnosticCentreByOwner(context.Request().Context(), *dto)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotAcceptable, err, context)
	}
	return utils.ResponseMessage(http.StatusNoContent, response, context)
}
