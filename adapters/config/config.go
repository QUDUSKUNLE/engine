package config

import (
	"fmt"
	"os"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/utils"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariable() error {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
		return nil
	}
	return nil
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

func LoadEnvironmentVariables(serviceName string) (*EnvConfiguration, error) {
	if err := LoadEnvironmentVariable(); err != nil {
		return nil, err
	}

	config := &EnvConfiguration{}

	switch serviceName {
	case "MEDIVUE":
		config.PORT = os.Getenv("PORT")
		config.DATABASE_URL = os.Getenv("DATABASE_URL")
		config.APP_URL = os.Getenv("APP_URL")
		config.EMAIL_HOST = os.Getenv("EMAIL_HOST")
		config.EMAIL_TYPE = os.Getenv("EMAIL_TYPE")
		config.EMAIL_PORT = os.Getenv("EMAIL_PORT")
		config.JWT_KEY = os.Getenv("JWT_SECRET_KEY")
		config.OPEN_API_KEY = os.Getenv("OPEN_API_KEY")
		config.ALLOW_ORIGINS = os.Getenv("ALLOW_ORIGINS")
		config.GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
		config.GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
		config.GOOGLE_REDIRECT_URL = os.Getenv("GOOGLE_REDIRECT_URL")
		config.SEND_GRID_API_KEY = os.Getenv("SEND_GRID_API_KEY")
		config.EMAIL_APP_PASSWORD = os.Getenv("EMAIL_APP_PASSWORD")
		config.EMAIL_USERNAME = os.Getenv("EMAIL_USERNAME")
		config.EMAIL_FROM_ADDRESS = os.Getenv("EMAIL_FROM_ADDRESS")
		// PAYSTACK
		config.PAYSTACK_BASE_URL = os.Getenv("PAYSTACK_BASE_URL")
		config.PAYSTACK_PUBLIC_KEY = os.Getenv("PAYSTACK_PUBLIC_KEY")
		config.PAYSTACK_SECRET_KEY = os.Getenv("PAYSTACK_SECRET_KEY")
		config.REDIS_URL = os.Getenv("REDIS_URL")

		// Validate required fields
		if config.PORT == "" {
			return nil, fmt.Errorf("missing required environment variables: PORT")
		}

		// Validate Google OAuth config
		if config.GOOGLE_CLIENT_ID == "" || config.GOOGLE_CLIENT_SECRET == "" {
			utils.Warn("Google OAuth configuration is incomplete. Google login will be unavailable.")
		}
		return config, nil
	default:
		// Validate required fields for other services
		if config.PORT == "" {
			return nil, fmt.Errorf("missing required environment variables: PORT")
		}
		return config, nil
	}
}
