package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/medivue/adapters/config"
	"github.com/medivue/core/utils"
	"golang.org/x/time/rate"
)

// SecureHeaders returns Echo middleware with strong default security headers.
func SecureHeaders() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'",
	})
}

func CORS(cfg *config.EnvConfiguration) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.AllowOrigins},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	})
}

func Logger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "id=${id} protocol=${protocol} time=${time}, remote_ip=${remote_ip}, latency=${latency}, method=${method}, uri=${uri}, status=${status}, host=${host}\n",
	})
}

func RateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/metrics" || c.Path() == "/health"
		},
		Store: middleware.NewRateLimiterMemoryStore(
			rate.Limit(10),
		),
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "rate limit exceeded",
			})
		},
	})
}

func Recover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			utils.Error("Panic recovered",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "stack", Value: string(stack)},
				utils.LogField{Key: "request_id", Value: requestID},
				utils.LogField{Key: "path", Value: c.Request().URL.Path},
			)
			return nil
		},
	})
}

func Timeout() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	})
}

func Gzip() echo.MiddlewareFunc {
	return middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5, // balanced compression (1 = fastest, 9 = best compression)
		Skipper: func(c echo.Context) bool {
			// Skip for metrics, health check, and file downloads (optional)
			return strings.HasPrefix(c.Path(), "/metrics") || strings.HasPrefix(c.Path(), "/health")
		},
	})
}
