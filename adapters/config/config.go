package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/medicue/core/utils"
)

func LoadEnvironmentVariable() error {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
		return nil
	}
	return nil
}

type Config struct {
	Port         string
	DB_URL       string
	JwtKey       string
	AllowOrigins string
	AppUrl       string
	// Google OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// Email Service
	SEND_GRID_API_KEY string
	EMAIL_FROM_ADDRESS string
}

func LoadConfig(serviceName string) (*Config, error) {
	if err := LoadEnvironmentVariable(); err != nil {
		return nil, err
	}

	config := &Config{}

	switch serviceName {
	case "MEDICUE":
		config.Port = os.Getenv("PORT")
		config.DB_URL = os.Getenv("DB_URL")
		config.JwtKey = os.Getenv("JWT_SECRET_KEY")
		config.AllowOrigins = os.Getenv("ALLOW_ORIGINS")
		config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
		config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
		config.GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
		config.AppUrl = os.Getenv("APP_URL")
		config.EMAIL_FROM_ADDRESS = os.Getenv("EMAIL_FROM_ADDRESS")
		config.SEND_GRID_API_KEY = os.Getenv("SEND_GRID_API_KEY")

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
