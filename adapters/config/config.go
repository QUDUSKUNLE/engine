package config

import (
	"fmt"
	"os"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/utils"
	"github.com/joho/godotenv"
)

// loadDotEnv loads variables from a .env file if present. It only warns on failure.
func loadDotEnv() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}
}

func DBConfig() db.DBConfig {
	return db.DBConfig{
		MaxConns:          20, // Increased for production load
		MinConns:          5,  // Increased minimum connections
		ConnTimeout:       15 * time.Second,
		MaxConnLifetime:   30 * time.Minute,
		MaxConnIdleTime:   5 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

type EnvConfiguration struct {
	PORT          string
	DATABASE_URL  string
	JWT_KEY       string
	ALLOW_ORIGINS string
	APP_URL       string
	// Google OAuth Configuration
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string
	// Generated token
	TOKEN_EXPIRED string
	// Email Service for SendGRID
	SEND_GRID_API_KEY string
	// Email Service
	EMAIL_USERNAME     string
	EMAIL_APP_PASSWORD string
	EMAIL_FROM_ADDRESS string
	EMAIL_TYPE         string
	EMAIL_HOST         string
	EMAIL_PORT         string
	// Paystack Service
	PAYSTACK_BASE_URL   string
	PAYSTACK_SECRET_KEY string
	PAYSTACK_PUBLIC_KEY string
	// OPEN API
	OPEN_API_KEY string
	// MONGODB_URL
	MONGODB_URL string
	// Redis Cache
	REDIS_URL string
}

// LoadEnvironmentVariables loads configuration from env and .env for the MEDIVUE service.
func LoadEnvironmentVariables() (*EnvConfiguration, error) {
	loadDotEnv()

	cfg := &EnvConfiguration{
		PORT:                 os.Getenv("PORT"),
		DATABASE_URL:         os.Getenv("DATABASE_URL"),
		APP_URL:              os.Getenv("APP_URL"),
		EMAIL_HOST:           os.Getenv("EMAIL_HOST"),
		EMAIL_TYPE:           os.Getenv("EMAIL_TYPE"),
		EMAIL_PORT:           os.Getenv("EMAIL_PORT"),
		JWT_KEY:              os.Getenv("JWT_SECRET_KEY"),
		OPEN_API_KEY:         os.Getenv("OPEN_API_KEY"),
		ALLOW_ORIGINS:        os.Getenv("ALLOW_ORIGINS"),
		GOOGLE_CLIENT_ID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_REDIRECT_URL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		SEND_GRID_API_KEY:    os.Getenv("SEND_GRID_API_KEY"),
		EMAIL_APP_PASSWORD:   os.Getenv("EMAIL_APP_PASSWORD"),
		EMAIL_USERNAME:       os.Getenv("EMAIL_USERNAME"),
		EMAIL_FROM_ADDRESS:   os.Getenv("EMAIL_FROM_ADDRESS"),
		PAYSTACK_BASE_URL:    os.Getenv("PAYSTACK_BASE_URL"),
		PAYSTACK_PUBLIC_KEY:  os.Getenv("PAYSTACK_PUBLIC_KEY"),
		PAYSTACK_SECRET_KEY:  os.Getenv("PAYSTACK_SECRET_KEY"),
		REDIS_URL:            os.Getenv("REDIS_URL"),
		MONGODB_URL:          os.Getenv("MONGODB_URL"),
		TOKEN_EXPIRED:        os.Getenv("TOKEN_EXPIRED"),
	}

	// Validate required fields
	if cfg.PORT == "" {
		return nil, fmt.Errorf("missing required environment variables: PORT")
	}

	// Validate Google OAuth config
	if cfg.GOOGLE_CLIENT_ID == "" || cfg.GOOGLE_CLIENT_SECRET == "" {
		utils.Warn("Google OAuth configuration is incomplete. Google login will be unavailable.")
	}

	return cfg, nil
}
