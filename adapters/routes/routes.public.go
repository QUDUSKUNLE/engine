package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
)

func PublicRoutesAdaptor(public *echo.Group, handler *handlers.HTTPHandler) *echo.Group {
	public.POST("/register", handler.Register)
	public.POST("/login", handler.SignIn)
	// public.POST("/delivery", handler.DeliveryProduct)
	// public.POST("/resetpassword", handler.ResetPassword)
	return public
}
