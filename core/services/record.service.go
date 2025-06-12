package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/utils"
)

// CreateMedicalRecord handles the creation of a new medical record.
func (service *ServicesHandler) CreateMedicalRecord(context echo.Context) error {
	ctx := context.Request().Context()
	uploader, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	dto, err := buildCreateMedicalRecordDto(context)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	// Set uploader info
	dto.UploaderType = uploader.UserType
	dto.UploaderAdminID = uploader.UserID

	// Validate uploader_admin_id and uploader_id before uploading data to cloud
	_, err = service.diagnosticRepo.GetDiagnosticCentreByManager(ctx, db.Get_Diagnostic_Centre_ByManagerParams{
		ID:      dto.UploaderID.String(),
		AdminID: dto.UploaderAdminID.String(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusForbidden, err, context)
	}

	// Upload file to cloud
	fileUrl, err := service.fileRepo.UploadFile(ctx, dto.FileUpload.Content)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnprocessableEntity, err, context)
	}

	// Parse SharedUntil string to time.Time
	var sharedUntilTime pgtype.Timestamp
	if dto.SharedUntil != "" {
		t, err := time.Parse("2006-01-02", dto.SharedUntil)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid SharedUntil format: %w", err), context)
		}
		sharedUntilTime = pgtype.Timestamp{Time: t, Valid: true}
	}

	record, err := service.recordRepo.CreateMedicalRecord(ctx, db.CreateMedicalRecordParams{
		UserID:          dto.UserID.String(),
		UploaderID:      dto.UploaderID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: dto.UploaderAdminID, Valid: true},
		UploaderType:    dto.UploaderType,
		ScheduleID:      dto.ScheduleID.String(),
		Title:           dto.Title,
		DocumentType:    dto.DocumentType,
		FilePath:        fileUrl,
		FileType:        pgtype.Text{String: dto.FileType, Valid: true},
		UploadedAt:      pgtype.Timestamp{Time: time.Now(), Valid: true},
		ProviderName:    pgtype.Text{String: dto.ProviderName, Valid: true},
		Specialty:       pgtype.Text{String: dto.Specialty, Valid: true},
		IsShared:        pgtype.Bool{Bool: dto.IsShared, Valid: true},
		SharedUntil:     sharedUntilTime,
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, record, context)
}

// GetMedicalRecord retrieves a single medical record.
func (service *ServicesHandler) GetMedicalRecord(context echo.Context) error {
	return nil
}

// GetMedicalRecords retrieves multiple medical records for a user.
func (service *ServicesHandler) GetMedicalRecords(context echo.Context) error {
	return nil
}

// GetUploaderMedicalRecord retrieves a single medical record uploaded by a specific uploader.
func (service *ServicesHandler) GetUploaderMedicalRecord(context echo.Context) error {
	return nil
}

// GetUploaderMedicalRecords retrieves multiple medical records uploaded by a specific uploader.
func (service *ServicesHandler) GetUploaderMedicalRecords(context echo.Context) error {
	return nil
}

// UpdateMedicalRecord updates an existing medical record by the uploader.
func (service *ServicesHandler) UpdateMedicalRecord(context echo.Context) error {
	return nil
}
