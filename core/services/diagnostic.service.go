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

// Helper to unmarshal address and contact, and build response
func buildDiagnosticCentreResponseFromRow(row *db.DiagnosticCentre, c echo.Context) (map[string]interface{}, error) {
	var address domain.Address
	if err := utils.UnmarshalJSONField(row.Address, &address, c); err != nil {
		return nil, err
	}
	var contact domain.Contact
	if err := utils.UnmarshalJSONField(row.Contact, &contact, c); err != nil {
		return nil, err
	}
	return buildDiagnosticCentreResponse(row, address, contact), nil
}

func (service *ServicesHandler) CreateDiagnosticCentre(context echo.Context) error {
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
	response, err := service.DiagnosticCentreRepo.CreateDiagnosticCentre(ctx, params)
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
	ctx := context.Request().Context()
	var params domain.GetDiagnosticParamDTO
	if err := utils.ValidateParams(context, &params); err != nil {
		fmt.Println(err.Error())
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	response, err := service.DiagnosticCentreRepo.GetDiagnosticCentre(ctx, params.DiagnosticCentreID)
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
	ctx := context.Request().Context()
	query, ok := context.Get("validatedQueryDTO").(*domain.SearchDiagnosticCentreQueryDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
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
	response, err := service.DiagnosticCentreRepo.GetNearestDiagnosticCentres(ctx, params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	result := make([]map[string]interface{}, 0, len(response))
	for _, v := range response {
		// Map v to a DiagnosticCentre struct
		diagnosticCentre := &db.DiagnosticCentre{
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
		}
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
	ctx := context.Request().Context()
	currentUser, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	body, ok := context.Get("validatedBodyDTO").(*domain.UpdateDiagnosticBodyDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequest), context)
	}
	param := context.Param("diagnostic_centre_id")
	addressBytes, err := utils.MarshalJSONField(body.Address, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	contactBytes, err := utils.MarshalJSONField(body.Contact, context)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	response, err := service.DiagnosticCentreRepo.UpdateDiagnosticCentreByOwner(ctx,
		db.Update_Diagnostic_Centre_ByOwnerParams{
			ID:                   param,
			CreatedBy:            currentUser.UserID.String(),
			DiagnosticCentreName: body.DiagnosticCentreName,
			Latitude:             pgtype.Float8{Float64: body.Latitude, Valid: true},
			Longitude:            pgtype.Float8{Float64: body.Longitude, Valid: true},
			Address:              addressBytes,
			Contact:              contactBytes,
			Doctors:              body.Doctors,
			AvailableTests:       body.AvailableTests,
			AdminID:              body.ADMINID.String(),
		})
	if err != nil {
		return utils.ErrorResponse(http.StatusNotAcceptable, err, context)
	}
	fmt.Println(body, response, "&&&&&&&&&&&&&&")
	return utils.ResponseMessage(http.StatusNoContent, "ok", context)
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
