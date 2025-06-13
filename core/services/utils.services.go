package services

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func buildCreateMedicalRecordDto(c echo.Context) (*domain.CreateMedicalRecordDTO, error) {
	d, ok := c.Get(utils.ValidatedBodyDTO).(*domain.CreateMedicalRecordDTO)
	if !ok || d == nil {
		return nil, c.JSON(http.StatusBadRequest, "Invalid DTO")
	}

	// If file content already exists in DTO, skip file parsing
	if d.FileUpload.Content != nil {
		return d, nil
	}

	// Parse file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "File upload error: file is required")
	}
	src, err := file.Open()
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "File open error: "+err.Error())
	}
	defer src.Close()
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "File read error")
	}
	ext := filepath.Ext(file.Filename)

	dto := &domain.CreateMedicalRecordDTO{
		UserID:          d.UserID,
		UploaderID:      d.UploaderID,
		UploaderAdminID: d.UploaderAdminID, // Copy UploaderAdminID
		UploaderType:    d.UploaderType,    // Copy UploaderType
		ScheduleID:      d.ScheduleID,
		Title:           d.Title,
		DocumentType:    d.DocumentType,
		FileUpload: domain.File{
			FileName: file.Filename,
			FileSize: file.Size,
			Content:  content,
		},
		FileType:     strings.TrimPrefix(ext, "."),
		UploadedAt:   d.UploadedAt,
		ProviderName: d.ProviderName,
		Specialty:    d.Specialty,
		IsShared:     d.IsShared,
		SharedUntil:  d.SharedUntil,
	}
	return dto, nil
}

// Helper to unmarshal address and contact, and build response
func buildDiagnosticCentreResponseFromRow(row *db.DiagnosticCentre, c echo.Context) (map[string]interface{}, error) {
	var address domain.Address
	if err := utils.UnmarshalJSONField(row.Address, &address, c); err != nil {
		return nil, err
	}
	var contact domain.Contact
	if err := utils.UnmarshalJSONField(row.Contact, &contact, c); err != nil {
		return nil, err
	}
	return buildDiagnosticCentreResponse(row, address, contact), nil
}

func buildCreateDiagnosticCentreParams(context echo.Context, value *domain.CreateDiagnosticDTO) (*db.Create_Diagnostic_CentreParams, error) {
	addressBytes, err := utils.MarshalJSONField(value.Address, context)
	if err != nil {
		return nil, err
	}
	contactBytes, err := utils.MarshalJSONField(value.Contact, context)
	if err != nil {
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
	return params, nil
}

func buildUpdateDiagnosticCentreByOwnerParams(context echo.Context, value *domain.UpdateDiagnosticBodyDTO) (*db.Update_Diagnostic_Centre_ByOwnerParams, error) {

	addressBytes, err := utils.MarshalJSONField(value.Address, context)
	if err != nil {
		return nil, err
	}
	contactBytes, err := utils.MarshalJSONField(value.Contact, context)
	if err != nil {
		return nil, err
	}
	body := &db.Update_Diagnostic_Centre_ByOwnerParams{
		DiagnosticCentreName: value.DiagnosticCentreName,
		Latitude:             pgtype.Float8{Float64: value.Latitude, Valid: true},
		Longitude:            pgtype.Float8{Float64: value.Longitude, Valid: true},
		Address:              addressBytes,
		Contact:              contactBytes,
		Doctors:              value.Doctors,
		AvailableTests:       value.AvailableTests,
		AdminID:              value.ADMINID.String(),
	}
	return body, nil
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

func buildDiagnosticCentre(value db.Get_Nearest_Diagnostic_CentresRow) *db.DiagnosticCentre {
	return &db.DiagnosticCentre{
		ID:                   value.ID,
		DiagnosticCentreName: value.DiagnosticCentreName,
		Latitude:             value.Latitude,
		Longitude:            value.Longitude,
		Address:              value.Address,
		Contact:              value.Contact,
		Doctors:              value.Doctors,
		AvailableTests:       value.AvailableTests,
		CreatedAt:            value.CreatedAt,
		UpdatedAt:            value.UpdatedAt,
	}
}

// PaginationParams interface for any struct that has Limit and Offset fields
type PaginationParams interface {
	GetLimit() int32
	GetOffset() int32
	SetLimit(limit int32)
	SetOffset(offset int32)
}

// SetDefaultPagination sets default values for any pagination parameters
func SetDefaultPagination[T PaginationParams](params T) T {
	if params.GetLimit() <= 0 {
		params.SetLimit(50) // Default limit to 50 if not set or invalid
	}
	if params.GetOffset() < 0 {
		params.SetOffset(0) // Default offset to 0 if negative
	}
	return params
}
