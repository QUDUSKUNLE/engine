package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

const (
	tokenExpiration = 72 * time.Hour
	authScheme      = "Bearer"
)

// JWTConfig returns an enhanced JWT configuration for Echo middleware
func JWTConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(context echo.Context) jwt.Claims {
			return new(domain.JwtCustomClaimsDTO)
		},
		SigningKey:    []byte(secret),
		SigningMethod: jwt.SigningMethodHS256.Name,
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c echo.Context, err error) error {
			utils.Error("JWT authentication failed",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "path", Value: c.Request().URL.Path},
				utils.LogField{Key: "method", Value: c.Request().Method},
				utils.LogField{Key: "remote_ip", Value: c.RealIP()})

			switch err {
			case echojwt.ErrJWTMissing:
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
			case echojwt.ErrJWTInvalid:
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
			// echojwt.ErrJWTExpired does not exist; handle expiration by inspecting the error message
			default:
				if err != nil && (err.Error() == "token is expired" || err.Error() == "Token is expired") {
					return echo.NewHTTPError(http.StatusUnauthorized, "Expired JWT token")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or malformed JWT")
			}
		},
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*domain.JwtCustomClaimsDTO)

			utils.Info("JWT authentication successful",
				utils.LogField{Key: "user_id", Value: claims.UserID},
				utils.LogField{Key: "user_type", Value: claims.UserType},
				utils.LogField{Key: "path", Value: c.Request().URL.Path})
		},
	}
}

// ValidateToken validates a JWT token string and returns the claims
func ValidateToken(tokenString string, secret string) (*domain.JwtCustomClaimsDTO, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JwtCustomClaimsDTO{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			utils.Error("Unexpected signing method",
				utils.LogField{Key: "method", Value: token.Method.Alg()})
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		utils.Error("Token validation failed",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JwtCustomClaimsDTO); ok && token.Valid {
		return claims, nil
	}

	utils.Error("Invalid token claims")
	return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
}

// GenerateToken creates a new JWT token for the given claims
func GenerateToken(claims domain.JwtCustomClaimsDTO, secret string) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(tokenExpiration))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.Error("Failed to generate token",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "user_id", Value: claims.UserID})
		return "", err
	}

	utils.Info("Token generated successfully",
		utils.LogField{Key: "user_id", Value: claims.UserID},
		utils.LogField{Key: "expires_at", Value: claims.ExpiresAt})

	return tokenString, nil
}
