package handlers

import (
	"github.com/labstack/echo/v4"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/register [post]
func (handler *HTTPHandler) Register(context echo.Context) error {
	return handler.service.Create(context)
}

// SignIn godoc
// @Summary User login
// @Description Authenticate a user and return a token
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/login [post]
func (handler *HTTPHandler) SignIn(context echo.Context) error {
	return handler.service.Login(context)
}

// CreateDiagnosticCentreManager godoc
// @Summary Create a diagnostic centre manager
// @Description Create a new diagnostic centre manager
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/diagnostic_centre_manager [post]
func (handler *HTTPHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	return handler.service.CreateDiagnosticCentreManager(context)
}
