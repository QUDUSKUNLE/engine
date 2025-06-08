package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/medicue/core/domain"
)

func JWTConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(context echo.Context) jwt.Claims {
			return new(domain.JwtCustomClaimsDTO)
		},
		SigningKey: []byte(secret),
	}
}
