package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Initialize logger with custom configuration
	logConfig := utils.LogConfig{
		Level:       "debug", // Set to debug in development, info in production
		OutputPath:  "logs/medicue.log",
		Development: true, // Set to false in production
	}
	if err := utils.InitLogger(logConfig); err != nil {
		panic(err)
	}
	defer utils.Logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig("MEDICUE")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize DB connection
	store, err := db.DatabaseConnection(cfg.DB_URL)
	if err != nil {
		log.Fatalf("Error connecting to the database")
	}

	// Create a new echo instance
	e := echo.New()

	// Setup Prometheus metrics endpoint first, before any other middleware
	e.GET("/metrics", echoprometheus.NewHandler())
	e.Use(middlewares.PrometheusMiddleware)
	e.Use(echoprometheus.NewMiddleware("Medicue"))

	// Plug echo into validationAdaptor
	e = middlewares.ValidationAdaptor(e)

	userRepo := repository.NewUserRepository(store)
	scheduleRepo := repository.NewScheduleRepository(store)
	diagnosticRepo := repository.NewDiagnosticCentreRepository(store)
	recordRepo := repository.NewRecordRepository(store)
	availabilityRepo := repository.NewAvailabilityRepository(store)
	paymentRepo := repository.NewPaymentRepository(store)
	appointmentRepo := repository.NewApppointmentRepository(store)
	core := services.ServicesAdapter(
		userRepo,
		scheduleRepo,
		diagnosticRepo,
		availabilityRepo,
		recordRepo,
		paymentRepo,
		appointmentRepo,
		*cfg,
	)
	// Initialize CronConfig
	cronConfig := config.GetConfig(userRepo, diagnosticRepo, appointmentRepo, *cfg)
	err = cronConfig.Start()
	if err != nil {
		log.Printf("Warning: Failed to start background services: %v", err)
	}
	defer cronConfig.Cleanup()

	httpHandler := handlers.HttpAdapter(core)

	v1 := e.Group("/v1")
	// Add a middleware to skip JWT validation for specific routes under /v1
	v1.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// List of /v1 routes that should NOT require JWT
			key := fmt.Sprintf("%s %s", c.Request().Method, c.Path())
			noAuthRoutes := map[string]bool{
				"POST /v1/login":                                   true,
				"POST /v1/register":                                true,
				"POST /v1/verify_email":                            true,
				"POST /v1/reset_password":                          true,
				"POST /v1/resend_verification":                     true,
				"POST /v1/request_password_reset":                  true,
				"GET /v1/diagnostic_centres":                       true,
				"GET /v1/diagnostic_centres/:diagnostic_centre_id": true,
				"POST /v1/auth/google":                             true,
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

	// Register routes
	routes.RoutesAdaptor(v1, httpHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

	// Add secure headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'",
	}))

	// Limit request body size to 25MB
	e.Use(middleware.BodyLimit("25M"))

	// Improved CORS config: restrict to trusted origins and methods
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.AllowOrigins},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time}, remote_ip=${remote_ip}, latency=${latency}, method=${method}, uri=${uri}, status=${status}, host=${host}\n",
	}))

	// Configure rate limiter with metrics endpoint excluded
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/metrics" // Skip rate limiting for metrics endpoint
		},
		Store: middleware.NewRateLimiterMemoryStore(rate.Limit(10)),
	}))

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server error:", err)
		} else {
			log.Printf("Server gracefully shutting down...")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
