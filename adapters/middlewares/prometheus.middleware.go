package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/medivue/adapters/metrics"
)

// PrometheusMiddleware adds custom metrics tracking to Echo
func PrometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		duration := time.Since(start).Seconds()
		metrics.ObserveHTTPRequest(c, duration)

		return err
	}
}
