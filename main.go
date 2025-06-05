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
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
)

// @title Bahsoon Africa API
// @version 1.0
// @description Bahsoon Africa API
// @host localhost:8080
// @BasePath /v1
func main() {
	// Get Port number from the loaded .env file
	cfg, err := config.LoadConfig("MEDICUE")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// Create a new echo instance
	e := echo.New()
	// Plug echo int validationAdaptor
	e = middlewares.ValidationAdaptor(e)

	e.Use(echoprometheus.NewMiddleware("Medicue"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time}, remote_ip=${remote_ip}, latency=${latency}, method=${method}, uri=${uri}, status=${status}, host=${host}\n",
	}))
	// Recover servers when break down
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.AllowOrigins},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))

	store, err := db.DatabaseConnection(cfg.DB_URL)
	if err != nil {
		log.Fatalf("Error connecting to the database")
	}

	repo := repository.NewPostgresRepository(store)
	core := services.ServicesAdapter(*repo)
	httpHandler := handlers.HttpAdapter(core)

	// Plug echo into PublicRoutesAdaptor
	public := e.Group("/v1")
	routes.PublicRoutesAdaptor(public, httpHandler)

	privateRoutes := e.Group("/v1")
	// Set JWT Configuration 
	conn := middlewares.JWTConfig(cfg.JwtKey)
	privateRoutes.Use(echojwt.WithConfig(conn))

	// Plug echo into PrivateRoutesAdaptor
	routes.PrivateRoutesAdaptor(privateRoutes, httpHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/metrics", echoprometheus.NewHandler())
	// Start the server on port 8080
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
