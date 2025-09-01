package middlewares

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	tokenExpiration   = 72 * time.Hour
	authScheme        = "Bearer"
	maxFileSize       = 32 << 20 // 32MB
	validatedBodyKey  = "validatedBodyDTO"
	validatedQueryKey = "validatedQueryParamDTO"
)

// Custom validator
type CustomValidator struct {
	validator *validator.Validate
}
