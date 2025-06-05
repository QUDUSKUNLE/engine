package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/adapters/middlewares"
	"github.com/medicue/core/domain"
)

func PublicRoutesAdaptor(public *echo.Group, handler *handlers.HTTPHandler) *echo.Group {
	public.POST("/register", handler.Register, middlewares.BodyValidationInterceptorFor(func() interface{} { return &domain.UserRegisterDTO{} }))
	public.POST("/login", handler.SignIn, middlewares.BodyValidationInterceptorFor(func() interface{} { return &domain.UserSignInDTO{} }))
	return public
}
