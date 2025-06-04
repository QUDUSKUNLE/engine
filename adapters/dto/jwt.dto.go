package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtCustomClaims struct {
	UserID       uuid.UUID `json:"user_id"`
	DiagnosticID uuid.UUID `json:"diagnostic_id"`
	UserType     string    `json:"user_type"`
	jwt.RegisteredClaims
}
