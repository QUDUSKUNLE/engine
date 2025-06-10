package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/medicue/adapters/config"
	"github.com/medicue/adapters/db"
	"github.com/medicue/adapters/db/repository"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/adapters/middlewares"
	"github.com/medicue/adapters/routes"
	"github.com/medicue/core/services"
	"github.com/medicue/core/utils"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
)

// @title Medicue
// @version 1.0
// @description Medicue API
// @host localhost:8080
// @BasePath /v1
func main() {
	if err := utils.InitLogger(); err != nil {
		panic(err)
	}
	defer utils.Logger.Sync()
	// Get Port number from the loaded .env file
	cfg, err := config.LoadConfig("MEDICUE")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// Create a new echo instance
	e := echo.New()
	// Plug echo int validationAdaptor
	e = middlewares.ValidationAdaptor(e)

	store, err := db.DatabaseConnection(cfg.DB_URL)
	if err != nil {
		log.Fatalf("Error connecting to the database")
	}

	userRepo := repository.NewUserRepository(store)
	scheduleRepo := repository.NewScheduleRepository(store)
	diagnosticRepo := repository.NewDiagnosticCentreRepository(store)
	core := services.ServicesAdapter(userRepo, scheduleRepo, diagnosticRepo)
	httpHandler := handlers.HttpAdapter(core)

	v1 := e.Group("/v1")
	// Add a middleware to skip JWT validation for specific routes under /v1
	v1.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// List of /v1 routes that should NOT require JWT
			key := fmt.Sprintf("%s %s", c.Request().Method, c.Path())
			noAuthRoutes := map[string]bool{
				"POST /v1/register":                                true,
				"POST /v1/login":                                   true,
				"GET /v1/diagnostic_centres":                       true,
				"GET /v1/diagnostic_centres/:diagnostic_centre_id": true,
			}
			if noAuthRoutes[key] {
				return next(c)
			}
			conn := middlewares.JWTConfig(cfg.JwtKey)
			conn.ErrorHandler = func(c echo.Context, err error) error {
				if c.Path() == "/v1/*" {
					return c.JSON(http.StatusNotFound, map[string]string{"error": "Ouch!!! Page not found"})
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or malformed jwt"})
			}
			return echojwt.WithConfig(conn)(next)(c)
		}
	})
	// Plug echo into PublicRoutesAdaptor
	routes.RoutesAdaptor(v1, httpHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/metrics", echoprometheus.NewHandler())

	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

	// Add secure headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	// Limit request body size to 2MB
	e.Use(middleware.BodyLimit("2M"))
	// Improved CORS config: restrict to trusted origins and methods
	e.Use(echoprometheus.NewMiddleware("Medicue"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.AllowOrigins},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time}, remote_ip=${remote_ip}, latency=${latency}, method=${method}, uri=${uri}, status=${status}, host=${host}\n",
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))

	// Start the server on port 8080
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
