package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/medivue/core/utils"
)

func LoadEnvironmentVariable() error {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
		return nil
	}
	return nil
}

type EnvConfiguration struct {
	Port         string
	DB_URL       string
	JwtKey       string
	AllowOrigins string
	AppUrl       string
	// Google OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// Generated token
	TokenExpired string

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
}

func LoadEnvironmentVariables(serviceName string) (*EnvConfiguration, error) {
	if err := LoadEnvironmentVariable(); err != nil {
		return nil, err
	}

	config := &EnvConfiguration{}

	switch serviceName {
	case "MEDIVUE":
		config.Port = os.Getenv("PORT")
		config.DB_URL = os.Getenv("DB_URL")
		config.JwtKey = os.Getenv("JWT_SECRET_KEY")
		config.AllowOrigins = os.Getenv("ALLOW_ORIGINS")
		config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
		config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
		config.GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
		config.AppUrl = os.Getenv("APP_URL")
		config.SEND_GRID_API_KEY = os.Getenv("SEND_GRID_API_KEY")
		config.EMAIL_APP_PASSWORD = os.Getenv("EMAIL_APP_PASSWORD")
		config.EMAIL_USERNAME = os.Getenv("EMAIL_USERNAME")
		config.EMAIL_FROM_ADDRESS = os.Getenv("EMAIL_FROM_ADDRESS")
		config.EMAIL_HOST = os.Getenv("EMAIL_HOST")
		config.EMAIL_TYPE = os.Getenv("EMAIL_TYPE")
		config.EMAIL_PORT = os.Getenv("EMAIL_PORT")
		config.OPEN_API_KEY = os.Getenv("OPEN_API_KEY")
		// PAYSTACK
		config.PAYSTACK_BASE_URL = os.Getenv("PAYSTACK_BASE_URL")
		config.PAYSTACK_PUBLIC_KEY = os.Getenv("PAYSTACK_PUBLIC_KEY")
		config.PAYSTACK_SECRET_KEY = os.Getenv("PAYSTACK_SECRET_KEY")

		// Validate required fields
		if config.Port == "" {
			return nil, fmt.Errorf("missing required environment variables: PORT")
		}

		// Validate Google OAuth config
		if config.GoogleClientID == "" || config.GoogleClientSecret == "" {
			utils.Warn("Google OAuth configuration is incomplete. Google login will be unavailable.")
		}
		return config, nil
	default:
		// Validate required fields for other services
		if config.Port == "" {
			return nil, fmt.Errorf("missing required environment variables: PORT")
		}
		return config, nil
	}
}
