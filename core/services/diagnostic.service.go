package services

import (
	"encoding/json"
	"errors"
	"fmt"

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
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
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
	var Address domain.Address
	if err := unmarshalJSONField(response.Address, &Address, context); err != nil {
		return err
	}
	var Contact domain.Contact
	if err := unmarshalJSONField(response.Contact, &Contact, context); err != nil {
		return err
	}
	res := buildDiagnosticCentreResponse(response, Address, Contact)
	return utils.ResponseMessage(http.StatusCreated, res, context)
}

func (service *ServicesHandler) GetDiagnosticCentre(context echo.Context) error {
	ctx := context.Request().Context()
	var params domain.GetDiagnosticParamDTO
	if err := utils.ValidateParams(context, &params); err != nil {
		fmt.Println(err.Error())
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	response, err := service.repositoryService.GetDiagnosticCentre(ctx, params.DiagnosticCentreID)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	var Address domain.Address
	if err := unmarshalJSONField(response.Address, &Address, context); err != nil {
		return err
	}
	var Contact domain.Contact
	if err := unmarshalJSONField(response.Contact, &Contact, context); err != nil {
		return err
	}
	res := buildDiagnosticCentreResponse(response, Address, Contact)
	return utils.ResponseMessage(http.StatusOK, res, context)
}

func (service *ServicesHandler) SearchDiagnosticCentre(context echo.Context) error {
	return nil
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

// Unmarshal to JSON
func unmarshalJSONField(data []byte, v interface{}, context echo.Context) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		utils.ErrorResponse(http.StatusInternalServerError, err, context)
		return err
	}
	return nil
}

// Helper to build diagnostic centre response
func buildDiagnosticCentreResponse(response *db.DiagnosticCentre, address domain.Address, contact domain.Contact) map[string]interface{} {
	return map[string]interface{}{
		"diagnostic_centre_id":   response.ID,
		"diagnostic_centre_name": response.DiagnosticCentreName,
		"latitude":               response.Latitude,
		"longitude":              response.Longitude,
		"address":                address,
		"contact":                contact,
		"doctors":                response.Doctors,
		"available_tests":        response.AvailableTests,
		"created_by":             response.CreatedBy,
		"admin_id":               response.AdminID,
		"created_at":             response.CreatedAt,
		"updated_at":             response.UpdatedAt,
	}
}
