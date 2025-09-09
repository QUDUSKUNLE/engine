package services

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
)

func (srv *ServicesHandler) GetNotifications(c echo.Context) error {
	// Authenticate and authorize user - owner or manager only
	currentUser, err := PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT, db.UserEnumDIAGNOSTICCENTREMANAGER, db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}

	query, _ := c.Get(utils.ValidatedQueryParamDTO).(*domain.GetNotificationsDTO)

	param, _ := SetDefaultPagination(&query.PaginationQueryDTO).(*domain.PaginationQueryDTO)

	params := db.GetUserNotificationsParams{
		UserID: currentUser.UserID.String(),
		Limit:  param.GetLimit(),
		Offset: param.GetOffset(),
	}
	response, err := srv.notificationRepo.GetNotifications(c.Request().Context(), params)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, c)
	}
	return utils.ResponseMessage(http.StatusOK, response, c)

}

func (srv *ServicesHandler) MarkNotificationRead(c echo.Context) error {
	// Authenticate and authorize user - owner or manager only
	currentUser, err := PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT, db.UserEnumDIAGNOSTICCENTREMANAGER, db.UserEnumDIAGNOSTICCENTREOWNER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}

	param, _ := c.Get(utils.ValidatedQueryParamDTO).(*domain.MarkNotificationReadDTO)

	response, err := srv.notificationRepo.MarkNotificationRead(c.Request().Context(), db.MarkAsReadParams{ID: param.NotificationID.String(), UserID: currentUser.UserID.String()})

	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, c)
	}
	return utils.ResponseMessage(http.StatusOK, response, c)

}
