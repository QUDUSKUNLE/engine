package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"Home": "Welcome to Diagnoxix Technologies"})
}

func Health(c echo.Context) error {
	health := map[string]string{
		"status":    "OK",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	health["database"] = "connected"
	return c.JSON(http.StatusOK, health)
}
