package services

import (
	"errors"
	"net/http"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/labstack/echo/v4"
)

// GetMedicalRecord retrieves a single medical record.
func (service *ServicesHandler) GetMedicalRecord(cont echo.Context) error {
	// Authentication check
	user, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, utils.AuthenticationRequired)
	}

	// This validated at the middleware level
	param, _ := cont.Get(utils.ValidatedQueryParamDTO).(*domain.GetMedicalRecordParamsDTO)

	// Fetch medical record
	response, err := service.recordPort.GetMedicalRecord(
		cont.Request().Context(),
		db.GetMedicalRecordParams{
			ID:     param.RecordID.String(),
			UserID: user.UserID.String(),
		},
	)
	if err != nil {
		cont.Logger().Error("Failed to get medical record:", err)
		switch {
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Medical record not found")
		case errors.Is(err, utils.ErrPermissionDenied):
			return echo.NewHTTPError(http.StatusForbidden, "Access denied to medical record")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve medical record")
		}
	}

	if response == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Medical record not found")
	}

	return utils.ResponseMessage(http.StatusOK, response, cont)
}

// GetMedicalRecords retrieves multiple medical records for a user.
func (service *ServicesHandler) GetMedicalRecords(cont echo.Context) error {
	user, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, cont)
	}

	// This validated at the middleware level
	query, _ := cont.Get(utils.ValidatedQueryParamDTO).(*domain.PaginationQueryDTO)

	query = SetDefaultPagination(query).(*domain.PaginationQueryDTO)

	response, err := service.recordPort.GetMedicalRecords(cont.Request().Context(), db.GetMedicalRecordsParams{
		UserID: user.UserID.String(),
		Limit:  query.GetLimit(),
		Offset: query.GetOffset(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, cont)
	}
	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, cont)
	}

	return utils.ResponseMessage(http.StatusOK, response, cont)
}
