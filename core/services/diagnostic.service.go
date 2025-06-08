package services

import (
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
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	dto, ok := context.Get("validatedBodyDTO").(*domain.CreateDiagnosticDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}

	addressBytes, err := utils.MarshalJSONField(dto.Address, context)
	if err != nil {
		return err
	}
	contactBytes, err := utils.MarshalJSONField(dto.Contact, context)
	if err != nil {
		return err
	}

	params := db.Create_Diagnostic_CentreParams{
		DiagnosticCentreName: dto.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: dto.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: dto.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              dto.Doctors,
		AvailableTests:       dto.AvailableTests,
		CreatedBy:            currentUser.UserID.String(),
		AdminID:              dto.AdminId.String(),
	}
	response, err := service.repositoryService.CreateDiagnosticCentre(ctx, params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	var address domain.Address
	if err := utils.UnmarshalJSONField(response.Address, &address, context); err != nil {
		return err
	}
	var contact domain.Contact
	if err := utils.UnmarshalJSONField(response.Contact, &contact, context); err != nil {
		return err
	}
	res := buildDiagnosticCentreResponse(response, address, contact)
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
	var address domain.Address
	if err := utils.UnmarshalJSONField(response.Address, &address, context); err != nil {
		return err
	}
	var contact domain.Contact
	if err := utils.UnmarshalJSONField(response.Contact, &contact, context); err != nil {
		return err
	}
	res := buildDiagnosticCentreResponse(response, address, contact)
	return utils.ResponseMessage(http.StatusOK, res, context)
}

func (service *ServicesHandler) SearchDiagnosticCentre(context echo.Context) error {
	ctx := context.Request().Context()
	queryParams, ok := context.Get("validatedQueryParamsDTO").(*domain.SearchDiagnosticCentreParamDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	params := db.Get_Nearest_Diagnostic_CentresParams{
		Radians:   queryParams.Latitude,
		Radians_2: queryParams.Longitude,
	}
	if queryParams.Doctor != "" {
		params.Doctors = []string{queryParams.Doctor}
	}
	if queryParams.Test != "" {
		params.AvailableTests = []string{queryParams.Test}
	}
	response, err := service.repositoryService.GetNearestDiagnosticCentres(ctx, params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	result := make([]map[string]interface{}, 0, len(response))
	for _, v := range response {
		var address domain.Address
		if err := utils.UnmarshalJSONField(v.Address, &address, context); err != nil {
			return err
		}
		var contact domain.Contact
		if err := utils.UnmarshalJSONField(v.Contact, &contact, context); err != nil {
			return err
		}
		// Map v to a DiagnosticCentre struct
		diagnosticCentre := &db.DiagnosticCentre{
			ID:                  v.ID,
			DiagnosticCentreName: v.DiagnosticCentreName,
			Latitude:            v.Latitude,
			Longitude:           v.Longitude,
			Address:             v.Address,
			Contact:             v.Contact,
			Doctors:             v.Doctors,
			AvailableTests:      v.AvailableTests,
			CreatedAt:           v.CreatedAt,
			UpdatedAt:           v.UpdatedAt,
		}
		item := buildDiagnosticCentreResponse(diagnosticCentre, address, contact)
		item["distance"] = v.DistanceKm
		item["distance_unit"] = "km"
		result = append(result, item)
	}
	return utils.ResponseMessage(http.StatusOK, result, context)
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
