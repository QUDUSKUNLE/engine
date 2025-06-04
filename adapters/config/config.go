package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariable() error {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
		return nil
	}
	return nil
}

type Config struct {
	Port   string
	DB_URL string
	JwtKey string
	AllowOrigins string
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
		// Validate required fields
		if config.Port == "" {
			return nil, fmt.Errorf("missing required environment variables: PORT")
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
