package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

var (
	// JWT related errors
	ErrMissingSecretKey = errors.New("missing JWT secret key")
	ErrInvalidToken     = errors.New("invalid or expired token")
	ErrUnauthorized     = errors.New("unauthorized to perform this operation")
	// Minimum password length for security
	MinPasswordLength = 12
	passwordChars     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?/|"
)

// PaginationParams interface for any struct that has Limit and Offset fields
type PaginationParams interface {
	GetLimit() int32
	GetOffset() int32
	SetLimit(limit int32)
	SetOffset(offset int32)
}

func buildCreateMedicalRecordDto(c echo.Context) (*domain.CreateMedicalRecordDTO, error) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error("Failed to get file from form",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		utils.Error("Failed to open uploaded file",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "filename", Value: file.Filename})
		return nil, err
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		utils.Error("Failed to read file content",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "filename", Value: file.Filename})
		return nil, err
	}

	dto := &domain.CreateMedicalRecordDTO{}
	if err := c.Bind(dto); err != nil {
		utils.Error("Failed to bind medical record data",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	// Add file information
	dto.FileUpload = domain.File{
		FileName: file.Filename,
		FileSize: file.Size,
		Content:  content,
	}

	utils.Info("Successfully built medical record DTO",
		utils.LogField{Key: "filename", Value: file.Filename},
		utils.LogField{Key: "file_size", Value: file.Size},
		utils.LogField{Key: "user_id", Value: dto.UserID})

	return dto, nil
}

// Helper to unmarshal address and contact, and build response
func buildDiagnosticCentreResponseFromRow(row *db.Get_Nearest_Diagnostic_CentresRow, c echo.Context) (map[string]interface{}, error) {
	var address domain.Address
	if err := utils.UnmarshalJSONField(row.Address, &address, c); err != nil {
		utils.Error("Failed to unmarshal address",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: row.ID})
		return nil, err
	}

	var contact domain.Contact
	if err := utils.UnmarshalJSONField(row.Contact, &contact, c); err != nil {
		utils.Error("Failed to unmarshal contact",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "diagnostic_centre_id", Value: row.ID})
		return nil, err
	}

	return buildDiagnosticCentreResponse(row, address, contact), nil
}

func buildCreateDiagnosticCentreParams(context echo.Context, value *domain.CreateDiagnosticDTO) (*db.Create_Diagnostic_CentreParams, error) {
	addressBytes, err := utils.MarshalJSONField(value.Address, context)
	if err != nil {
		utils.Error("Failed to marshal address data",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	contactBytes, err := utils.MarshalJSONField(value.Contact, context)
	if err != nil {
		utils.Error("Failed to marshal contact data",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	availableTests := make([]string, len(value.AvailableTests))
	copy(availableTests, value.AvailableTests)

	doctors := make([]string, len(value.Doctors))
	for i, doctor := range value.Doctors {
		doctors[i] = string(doctor)
	}

	params := &db.Create_Diagnostic_CentreParams{
		DiagnosticCentreName: value.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: value.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: value.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              doctors,
		AvailableTests:       availableTests,
		AdminID:              value.AdminId.String(),
	}

	utils.Info("Built diagnostic centre parameters",
		utils.LogField{Key: "centre_name", Value: value.DiagnosticCentreName},
		utils.LogField{Key: "admin_id", Value: value.AdminId})

	return params, nil
}

func buildUpdateDiagnosticCentreByOwnerParams(context echo.Context, value *domain.UpdateDiagnosticBodyDTO) (*db.Update_Diagnostic_Centre_ByOwnerParams, error) {
	addressBytes, err := utils.MarshalJSONField(value.Address, context)
	if err != nil {
		utils.Error("Failed to marshal address data for update",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	contactBytes, err := utils.MarshalJSONField(value.Contact, context)
	if err != nil {
		utils.Error("Failed to marshal contact data for update",
			utils.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	doctors := make([]string, len(value.Doctors))
	copy(doctors, value.Doctors)

	availableTests := make([]string, len(value.AvailableTests))
	copy(availableTests, value.AvailableTests)

	params := &db.Update_Diagnostic_Centre_ByOwnerParams{
		DiagnosticCentreName: value.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: value.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: value.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              doctors,
		AvailableTests:       availableTests,
		AdminID:              value.ADMINID.String(),
	}

	utils.Info("Built update diagnostic centre parameters",
		utils.LogField{Key: "centre_name", Value: value.DiagnosticCentreName},
		utils.LogField{Key: "admin_id", Value: value.ADMINID})

	return params, nil
}

// Helper to build diagnostic centre response
func buildDiagnosticCentreResponse(response *db.Get_Nearest_Diagnostic_CentresRow, address domain.Address, contact domain.Contact) map[string]interface{} {
	return map[string]interface{}{
		"diagnostic_centre_id":   response.ID,
		"diagnostic_centre_name": response.DiagnosticCentreName,
		"latitude":               response.Latitude,
		"longitude":              response.Longitude,
		"address":                address,
		"contact":                contact,
		"doctors":                response.Doctors,
		"available_tests":        response.AvailableTests,
		"created_at":             response.CreatedAt,
		"updated_at":             response.UpdatedAt,
		"availability":           response.Availability,
	}
}

// isValidUserType checks if the given user type is one of the allowed types
func isValidUserType(allowedTypes []db.UserEnum, userType db.UserEnum) bool {
	for _, t := range allowedTypes {
		if t == userType {
			return true
		}
	}
	return false
}

// SetDefaultPagination sets default values for pagination parameters if not provided
func SetDefaultPagination(params PaginationParams) PaginationParams {
	if params.GetLimit() <= 0 {
		params.SetLimit(10) // Default limit
	}
	if params.GetOffset() < 0 {
		params.SetOffset(0) // Default offset
	}
	return params
}

// isValidLatitude checks if the latitude is within valid range (-90 to 90)
func isValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

// isValidLongitude checks if the longitude is within valid range (-180 to 180)
func isValidLongitude(lon float64) bool {
	return lon >= -180 && lon <= 180
}

func toNumeric(n float64) pgtype.Numeric {
	var num pgtype.Numeric
	err := num.Scan(n)
	if err != nil {
		// If scan fails, try converting through string to maintain decimal precision
		err = num.Scan(fmt.Sprintf("%.2f", n))
		if err != nil {
			num.Valid = false
			return num
		}
	}
	num.Valid = true
	return num
}

func toUUID(id string) pgtype.UUID {
	var uid pgtype.UUID
	if parsed, err := uuid.Parse(id); err == nil {
		uid.Bytes = parsed
		uid.Valid = true
	}
	return uid
}

// Helper functions and type definitions

func toTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func toText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: len(s) > 0}
}

// GenerateToken generates a JWT token for the given user with enhanced security
func GenerateToken(user domain.CurrentUserDTO) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		utils.Error("jwt secret key missing")
		return "", ErrMissingSecretKey
	}

	// Get token expiration from env or use default
	tokenExpiration := 72 * time.Hour
	if expStr := os.Getenv("JWT_EXPIRATION_HOURS"); expStr != "" {
		if exp, err := time.ParseDuration(expStr + "h"); err == nil {
			tokenExpiration = exp
		}
	}

	claims := &domain.JwtCustomClaimsDTO{
		UserID:       user.UserID,
		DiagnosticID: user.DiagnosticID,
		UserType:     db.UserEnum(user.UserType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.Error("failed to sign token", utils.LogField{Key: "error", Value: err.Error()})
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	utils.Info("token generated successfully",
		utils.LogField{Key: "userID", Value: user.UserID.String()},
		utils.LogField{Key: "userType", Value: string(user.UserType)},
	)
	return signedToken, nil
}

// GenerateRandomPassword generates a cryptographically secure random password
func GenerateRandomPassword(length int) (string, error) {
	if length < MinPasswordLength {
		return "", fmt.Errorf("password length must be at least %d characters", MinPasswordLength)
	}

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			utils.Error("failed to generate random number", utils.LogField{Key: "error", Value: err.Error()})
			return "", fmt.Errorf("failed to generate password: %w", err)
		}
		password[i] = passwordChars[num.Int64()]
	}
	return string(password), nil
}

// GenerateRandomToken creates a cryptographically secure random token
func GenerateRandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// CurrentUser extracts the current user from the JWT token in the context
func CurrentUser(c echo.Context) (*domain.CurrentUserDTO, error) {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		utils.Error("missing or invalid user token")
		return nil, ErrInvalidToken
	}

	claims, ok := userToken.Claims.(*domain.JwtCustomClaimsDTO)
	if !ok {
		utils.Error("invalid token claims")
		return nil, ErrInvalidToken
	}

	utils.Debug("user extracted from token",
		utils.LogField{Key: "userID", Value: claims.UserID.String()},
		utils.LogField{Key: "userType", Value: string(claims.UserType)},
	)

	return &domain.CurrentUserDTO{
		UserID:       claims.UserID,
		DiagnosticID: claims.DiagnosticID,
		UserType:     claims.UserType,
	}, nil
}

// PrivateMiddlewareContext validates user type access
func PrivateMiddlewareContext(c echo.Context, userTypes []db.UserEnum) (*domain.CurrentUserDTO, error) {
	user, err := CurrentUser(c)
	if err != nil {
		return nil, err
	}

	for _, allowedType := range userTypes {
		if user.UserType == allowedType {
			return user, nil
		}
	}

	utils.Warn("unauthorized access attempt",
		utils.LogField{Key: "requiredTypes", Value: convertUserTypesToStrings(userTypes)},
		utils.LogField{Key: "actualType", Value: string(user.UserType)},
		utils.LogField{Key: "userID", Value: user.UserID.String()},
	)
	return nil, ErrUnauthorized
}

// Helper function to convert UserEnum slice to string slice for logging
func convertUserTypesToStrings(types []db.UserEnum) []string {
	strings := make([]string, len(types))
	for i, t := range types {
		strings[i] = string(t)
	}
	return strings
}
