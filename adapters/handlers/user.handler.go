package handlers

import (
	"github.com/labstack/echo/v4"
)

// @Summary Register a new user
// @Description register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param body body domain.UserDto true "Register a user"
// @Failure 409 {object} domain.Response
// @Success 201 {object} domain.Response
// @Router /register [post]
func (handler *HTTPHandler) Register(context echo.Context) error {
	return handler.service.Create(context)
}

func (handler *HTTPHandler) SignIn(context echo.Context) error {
	return handler.service.Login(context)
}

func (handler *HTTPHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	return handler.service.CreateDiagnosticCentreManager(context)
}
