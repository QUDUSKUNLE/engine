package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) CreateDiagnsoticCentre(context echo.Context) error {
	ctx := context.Request().Context()
	current_user, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	dto, ok := context.Get("validatedDTO").(*domain.CreateDiagnosticDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequestBody), context)
	}

	addressBytes, err := marshalJSONField(dto.Address, context)
	if err != nil {
		return err
	}
	contactBytes, err := marshalJSONField(dto.Contact, context)
	if err != nil {
		return err
	}

	diagnostic_centre_param := db.Create_Diagnostic_CentreParams{
		DiagnosticCentreName: dto.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: dto.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: dto.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              dto.Doctors,
		AvailableTests:       dto.AvailableTests,
		CreatedBy:            current_user.UserID.String(),
		AdminID:              dto.AdminId.String(),
	}
	response, err := service.repositoryService.CreateDiagnosticCentre(ctx, diagnostic_centre_param)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, response, context)
}

func (service *ServicesHandler) GetDiagnosticCentre(context echo.Context) error {
	ctx := context.Request().Context()
	dto, ok := context.Get("validatedDTO").(*domain.GetDiagnosticParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequestBody), context)
	}
	response, err := service.repositoryService.GetDiagnosticCentre(ctx, dto.DiagnosticCentreID.String())
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

// Helper to marshal JSON fields and handle errors
func marshalJSONField(field interface{}, context echo.Context) ([]byte, error) {
	data, err := json.Marshal(field)
	if err != nil {
		utils.ErrorResponse(http.StatusInternalServerError, err, context)
		return nil, err
	}
	return data, nil
}
