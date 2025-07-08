package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/medivue/core/domain"
	"github.com/medivue/core/utils"
)

const (
	tokenExpiration = 72 * time.Hour
	authScheme      = "Bearer"
)

// JWTConfig returns an enhanced JWT configuration for Echo middleware
func jWTConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(context echo.Context) jwt.Claims {
			return new(domain.JwtCustomClaimsDTO)
		},
		SigningKey: []byte(secret),
		BeforeFunc: func(c echo.Context) {
			// Log sanitized auth header for debugging (only first few chars of token)
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 {
					tokenPreview := parts[1]
					if len(tokenPreview) > 10 {
						tokenPreview = tokenPreview[:10] + "..."
					}
					utils.Debug("JWT token received",
						utils.LogField{Key: "token_preview", Value: tokenPreview},
						utils.LogField{Key: "scheme", Value: parts[0]})
				}
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			authHeader := c.Request().Header.Get("Authorization")

			utils.Error("JWT authentication failed",
				utils.LogField{Key: "error", Value: err.Error()},
				utils.LogField{Key: "path", Value: c.Request().URL.Path},
				utils.LogField{Key: "method", Value: c.Request().Method},
				utils.LogField{Key: "remote_ip", Value: c.RealIP()},
				utils.LogField{Key: "has_auth_header", Value: authHeader != ""})

			// Check Authorization header format first
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format. Use 'Bearer <token>'")
				}
				if parts[0] != authScheme {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization scheme. Use 'Bearer'")
				}
				// If header exists but jwt.Error is missing, the token format is wrong
				if err == nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "Malformed JWT token")
				}
			}

			switch err {
			case echojwt.ErrJWTMissing:
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
			case echojwt.ErrJWTInvalid:
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token: token signature is invalid")
			default:
				if err != nil {
					errMsg := err.Error()
					switch {
					case errMsg == "token is expired" || errMsg == "Token is expired":
						return echo.NewHTTPError(http.StatusUnauthorized, "Expired JWT token")
					case strings.Contains(errMsg, "signing method"):
						return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token: incorrect signing method")
					case strings.Contains(errMsg, "claims"):
						return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token: malformed claims")
					case strings.Contains(errMsg, "token contains an invalid number of segments"):
						return echo.NewHTTPError(http.StatusUnauthorized, "Malformed JWT token: invalid token format")
					default:
						return echo.NewHTTPError(http.StatusUnauthorized, "Malformed JWT token: "+errMsg)
					}
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or missing JWT token")
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

// ConditionalJWTMiddleware skips JWT for unauthenticated routes
func ConditionalJWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	noAuthRoutes := map[string]bool{
		"POST /v1/login":                                   true,
		"POST /v1/register":                                true,
		"GET /v1/verify_email":                             true,
		"POST /v1/reset_password":                          true,
		"POST /v1/resend_verification":                     true,
		"POST /v1/request_password_reset":                  true,
		"GET /v1/diagnostic_centres":                       true,
		"GET /v1/diagnostic_centres/:diagnostic_centre_id": true,
		"POST /v1/auth/google":                             true,
		"GET /v1/health":                                   true,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := fmt.Sprintf("%s %s", c.Request().Method, c.Path())

			if noAuthRoutes[key] {
				return next(c)
			}

			cfg := jWTConfig(jwtSecret)
			cfg.ErrorHandler = func(c echo.Context, err error) error {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or malformed jwt"})
			}

			return echojwt.WithConfig(cfg)(next)(c)
		}
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
