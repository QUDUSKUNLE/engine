package services

import (
	"io"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
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
func buildDiagnosticCentreResponseFromRow(row *db.DiagnosticCentre, c echo.Context) (map[string]interface{}, error) {
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

	params := &db.Create_Diagnostic_CentreParams{
		DiagnosticCentreName: value.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: value.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: value.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              value.Doctors,
		AvailableTests:       value.AvailableTests,
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

	params := &db.Update_Diagnostic_Centre_ByOwnerParams{
		DiagnosticCentreName: value.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: value.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: value.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              value.Doctors,
		AvailableTests:       value.AvailableTests,
		AdminID:              value.ADMINID.String(),
	}

	utils.Info("Built update diagnostic centre parameters",
		utils.LogField{Key: "centre_name", Value: value.DiagnosticCentreName},
		utils.LogField{Key: "admin_id", Value: value.ADMINID})

	return params, nil
}

// Helper to build diagnostic centre response
func buildDiagnosticCentreResponse(response *db.DiagnosticCentre, address domain.Address, contact domain.Contact) map[string]interface{} {
	return map[string]interface{}{
		"diagnostic_centre_id":   response.ID,
		"diagnostic_centre_name": response.DiagnosticCentreName,
		"latitude":               response.Latitude,
		"longitude":              response.Longitude,
		"address":                address,
		"contact":                contact,
		"doctors":                response.Doctors,
		"available_tests":        response.AvailableTests,
		"created_by":             response.CreatedBy,
		"admin_id":               response.AdminID,
		"created_at":             response.CreatedAt,
		"updated_at":             response.UpdatedAt,
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
