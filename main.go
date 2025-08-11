package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diagnoxix/adapters/config"
	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/adapters/db/repository"
	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/adapters/middlewares"
	"github.com/diagnoxix/adapters/routes"
	"github.com/diagnoxix/core/services"
	"github.com/diagnoxix/core/utils"
	_ "github.com/diagnoxix/swaggerdocs"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Diagnoxix
// @version 1.0
// @description Diagnoxix API
// @host 127.0.0.1:7556
// @BasePath /
func main() {
	// Initialize logger with custom configuration
	logConfig := utils.LogConfig{
		Level:       "info", // Set to debug in development, info in production
		OutputPath:  "logs/medivue.log",
		Development: false, // Set to false in production
	}
	if err := utils.InitLogger(logConfig); err != nil {
		panic(err)
	}
	defer func() {
		if err := utils.Logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Load configuration
	cfg, err := config.LoadEnvironmentVariables("MEDIVUE")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// Initialize DB connection
	store, conn, err := db.DatabaseConnection(context.Background(), cfg.DB_URL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Create a new echo instance
	e := echo.New()

	// Setup Prometheus metrics endpoint first, before any other middleware
	// ---- Metrics and Observability ----
	e.GET("/metrics", echoprometheus.NewHandler())
	e.Use(middlewares.PrometheusMiddleware)
	e.Use(echoprometheus.NewMiddleware("Diagnoxix"))

	// ---- Middleware Stack ----
	// Add basic observability middleware
	e.Use(middleware.RequestID())
	// Limit request body size to 25MB
	e.Use(middleware.BodyLimit("25M"))
	// Add gzip compression
	e.Use(middlewares.Gzip())
	e.Use(middlewares.Logger())
	// Add secure headers
	e.Use(middlewares.SecureHeaders())
	// Improved CORS config: restrict to trusted origins and methods
	e.Use(middlewares.CORS(cfg))
	// Add request timing middleware
	e.Use(middlewares.Timeout())
	// Configure rate limiter with metrics endpoint excluded
	e.Use(middlewares.RateLimiter())

	e.Use(middlewares.Recover())
	// Initialize all repositories
	repos := repository.InitializeRepositories(store, conn)

	// Plug echo into validationAdaptor
	e = middlewares.ValidationAdaptor(e)

	// Initialize all services
	services := services.InitializeServices(repos, cfg)

	// Initialize CronConfig with repositories
	cronConfig := config.GetConfig(repos.User, repos.Diagnostic, repos.Appointment, *cfg)
	if err := cronConfig.Start(); err != nil {
		log.Printf("Warning: Failed to start background services: %v", err)
	}
	defer cronConfig.Cleanup()

	// Initialize HTTP handlers with core services
	httpHandler := handlers.HttpAdapter(services.Core)

	// Add a middleware to skip JWT validation for specific routes under /v1
	v1 := e.Group("/v1")
	v1.Use(middlewares.ConditionalJWTMiddleware(cfg.JwtKey))

	// Register routes
	routes.RoutesAdaptor(v1, httpHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

	// Update health check endpoint
	e.GET("/health", handlers.Health)
	e.GET("", handlers.Home)

	// Get port from environment (Railway and most PaaS set PORT)
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(":" + port); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
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
