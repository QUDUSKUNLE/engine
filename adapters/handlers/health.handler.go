package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Service   string    `json:"service"`
}

func (h *HTTPHandler) HealthCheck(c echo.Context) error {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Service:   "DiagnoxixAI API",
	}
	return c.JSON(http.StatusOK, response)
}
