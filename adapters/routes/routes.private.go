package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/adapters/middlewares"
	"github.com/medicue/core/domain"
)

func PrivateRoutesAdaptor(
	private *echo.Group,
	handler *handlers.HTTPHandler,
) *echo.Group {
	private.POST(
		"/diagnostic_centre_manager",
		handler.CreateDiagnosticCentreManager,
		middlewares.BodyValidationInterceptorFor(func() interface{} {
			return &domain.DiagnosticCentreManagerRegisterDTO{}
		}),
	)
	// private.GET("/shipments", handler.GetShippings)

	// private.PUT("/pickups", handler.PutPickUp)
	// private.GET("/pickups", handler.GetPickUps)
	// private.GET("/pickups/:pick_up_id", handler.GetPickUp)

	// private.POST("/addresses", handler.PostAddress)
	// private.GET("/addresses", handler.GetAddresses)
	// private.GET("/addresses/:address_id", handler.GetAddress)
	// private.PUT("/addresses/:address_id", handler.PutAddress)
	// private.DELETE("/addresses/:address_id", handler.DeleteAddress)

	// private.POST("/packagings", handler.PostPackaging)
	// private.GET("/packagings", handler.GetPackagings)
	// private.GET("/packagings/:packaging_id", handler.GetPackaging)
	// private.PUT("/packagings/:packaging_id", handler.PutPackaging)
	// private.DELETE("/packagings/:packaging_id", handler.DeletePackaging)

	// private.POST("/parcels", handler.PostParcel)
	// private.GET("/parcels", handler.GetParcels)
	// private.GET("/parcels/:parcel_id", handler.GetParcel)
	// private.PUT("/parcels/:parcel_id", handler.PutParcel)
	// private.DELETE("/parcels/:parcel_id", handler.DeleteParcel)

	// private.GET("/rates", handler.Rates)

	return private
}
