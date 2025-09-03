package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// CreateMedicalRecord handles the creation of a new medical record.
func (service *ServicesHandler) CreateMedicalRecord(cont echo.Context) error {
	// Authentication & Authorization
	uploader, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, utils.ErrUnauthorized, cont)
	}

	// Build and validate DTO
	dto, err := buildCreateMedicalRecordDto(cont)
	if err != nil {
		cont.Logger().Error("Failed to build medical record DTO:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}
	if dto.FileUpload.Content == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File upload is required")
	}
	ctx := cont.Request().Context()

	// Set uploader info
	dto.UploaderType = uploader.UserType
	dto.UploaderAdminID = uploader.UserID

	// Validate uploader_admin_id and uploader_id before uploading data to cloud
	params := db.Get_Diagnostic_Centre_ByManagerParams{
		ID:      dto.DiagnosticCentreID.String(),
		AdminID: pgtype.UUID{Bytes: uploader.UserID, Valid: true},
	}
	_, err = service.diagnosticPort.GetDiagnosticCentreByManager(ctx, params)
	if err != nil {
		return utils.ErrorResponse(http.StatusNotFound, utils.ErrNotFoundDiagnositcCentre, cont)
	}

	// Parse SharedUntil string to time.Time
	var sharedUntilTime pgtype.Timestamp
	if dto.SharedUntil != "" {
		t, err := time.Parse(time.RFC3339, dto.SharedUntil)
		if err != nil {
			t, err = time.Parse("2006-01-02", dto.SharedUntil)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format. Expected RFC3339 or YYYY-MM-DD")
			}
		}
		sharedUntilTime = pgtype.Timestamp{Time: t, Valid: true}
	}

	// Create medical record
	createParams := db.CreateMedicalRecordParams{
		UserID:          dto.UserID.String(),
		UploaderID:      dto.DiagnosticCentreID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: dto.UploaderAdminID, Valid: true},
		UploaderType:    db.UserEnum(dto.UploaderType),
		ScheduleID:      dto.ScheduleID.String(),
		Title:           dto.Title,
		DocumentType:    dto.DocumentType,
		FilePath:        "uplaoding",
		FileType:        pgtype.Text{String: dto.FileType, Valid: true},
		UploadedAt:      pgtype.Timestamp{Time: time.Now(), Valid: true},
		ProviderName:    pgtype.Text{String: dto.ProviderName, Valid: true},
		Specialty:       pgtype.Text{String: dto.Specialty, Valid: true},
		IsShared:        pgtype.Bool{Bool: dto.IsShared, Valid: true},
		SharedUntil:     sharedUntilTime,
	}

	record, err := service.recordPort.CreateMedicalRecord(ctx, createParams)
	if err != nil {
		cont.Logger().Error("Failed to create medical record:", err)
		switch {
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, utils.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid record data")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create medical record")
		}
	}

	go func(recordID string, fileContent []byte) {
		fileUrl, err := service.filePort.UploadFile(context.Background(), fileContent)
		// service.aiPort.
		if err != nil {
			cont.Logger().Errorf("File upload failed for record %s: %v", recordID, err)
			return
		}
		_, err = service.aiPort.OCR.Parse(context.Background(), fileContent)
		if err != nil {
			fmt.Printf("Error %v", err)
			return
		}
		// Update DB with final file URL
		_, err = service.recordPort.UpdateFilePath(context.Background(), db.UpdateFilePathParams{
			ID:       recordID,
			FilePath: fileUrl,
		})
		if err != nil {
			cont.Logger().Errorf("Failed to update file path for record %s: %v", recordID, err)
			return
		}
	}(record.ID, dto.FileUpload.Content)

	return utils.ResponseMessage(http.StatusAccepted, map[string]string{
		"record_id": record.ID,
		"status":    "File is uploading in background",
	}, cont)
}

// GetMedicalRecord retrieves a single medical record.
func (service *ServicesHandler) GetMedicalRecord(cont echo.Context) error {
	// Authentication check
	user, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, utils.AuthenticationRequired)
	}

	// This validated at the middleware level
	param, _ := cont.Get(utils.ValidatedQueryParamDTO).(*domain.GetMedicalRecordParamsDTO)

	// Fetch medical record
	response, err := service.recordPort.GetMedicalRecord(
		cont.Request().Context(),
		db.GetMedicalRecordParams{
			ID:     param.RecordID.String(),
			UserID: user.UserID.String(),
		},
	)
	if err != nil {
		cont.Logger().Error("Failed to get medical record:", err)
		switch {
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Medical record not found")
		case errors.Is(err, utils.ErrPermissionDenied):
			return echo.NewHTTPError(http.StatusForbidden, "Access denied to medical record")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve medical record")
		}
	}

	if response == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Medical record not found")
	}

	return utils.ResponseMessage(http.StatusOK, response, cont)
}

// GetMedicalRecords retrieves multiple medical records for a user.
func (service *ServicesHandler) GetMedicalRecords(cont echo.Context) error {
	user, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, cont)
	}

	// This validated at the middleware level
	query, _ := cont.Get(utils.ValidatedQueryParamDTO).(*domain.PaginationQueryDTO)

	query = SetDefaultPagination(query).(*domain.PaginationQueryDTO)

	response, err := service.recordPort.GetMedicalRecords(cont.Request().Context(), db.GetMedicalRecordsParams{
		UserID: user.UserID.String(),
		Limit:  query.GetLimit(),
		Offset: query.GetOffset(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, cont)
	}
	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, cont)
	}

	return utils.ResponseMessage(http.StatusOK, response, cont)
}

// GetUploaderMedicalRecord retrieves a single medical record uploaded by a specific uploader.
func (service *ServicesHandler) GetUploaderMedicalRecord(cont echo.Context) error {
	manager, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, cont)
	}
	// This validated at the middleware level
	query, _ := cont.Get(utils.ValidatedQueryParamDTO).(*domain.GetUploaderMedicalRecordParamsDTO)

	response, err := service.recordPort.GetUploaderMedicalRecord(cont.Request().Context(), db.GetUploaderMedicalRecordParams{
		ID:              query.RecordID.String(),
		UploaderID:      query.DiagnosticCentreID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: manager.UserID, Valid: true},
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, cont)
	}
	return utils.ResponseMessage(http.StatusOK, response, cont)
}

// GetUploaderMedicalRecords retrieves multiple medical records uploaded by a specific uploader.
func (service *ServicesHandler) GetUploaderMedicalRecords(cont echo.Context) error {
	admin, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, cont)
	}

	// This validated at the middleware level
	query, _ := cont.Get(utils.ValidatedQueryParamDTO).(*domain.GetUploaderMedicalRecordsParamQueryDTO)

	query = SetDefaultPagination(query).(*domain.GetUploaderMedicalRecordsParamQueryDTO)

	response, err := service.recordPort.GetUploaderMedicalRecords(cont.Request().Context(), db.GetUploaderMedicalRecordsParams{
		UploaderID:      query.DiagnosticCentreID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: admin.UserID, Valid: true},
		Limit:           query.GetLimit(),
		Offset:          query.GetOffset(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, cont)
	}
	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, cont)
	}

	return utils.ResponseMessage(http.StatusOK, response, cont)
}

// UpdateMedicalRecord updates an existing medical record by the uploader.
func (service *ServicesHandler) UpdateMedicalRecord(cont echo.Context) error {
	ctx := cont.Request().Context()

	// Authentication & Authorization
	manager, err := PrivateMiddlewareContext(cont, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	// Get the validated DTO from cont
	dto, ok := cont.Get(utils.ValidatedBodyDTO).(*domain.UpdateMedicalRecordDTO)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	// Set uploader admin ID from the authenticated user
	dto.UploaderAdminID = manager.UserID

	// Parse SharedUntil string to time.Time if provided
	var sharedUntilTime pgtype.Timestamp
	if dto.SharedUntil != "" {
		// Try RFC3339 format first (includes time)
		t, err := time.Parse(time.RFC3339, dto.SharedUntil)
		if err != nil {
			// Fallback to date-only format
			t, err = time.Parse("2006-01-02", dto.SharedUntil)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format. Expected RFC3339 or YYYY-MM-DD")
			}
		}
		sharedUntilTime = pgtype.Timestamp{Time: t, Valid: true}
	}

	// Handle file upload if provided
	var fileUrl string
	if dto.FileUpload.Content != nil {
		fileUrl, err = service.filePort.UploadFile(ctx, dto.FileUpload.Content)
		if err != nil {
			cont.Logger().Error("File upload failed:", err)
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Failed to upload file")
		}
	}

	// Parse document date if provided
	var documentDate pgtype.Date
	if dto.DocumentDate != "" {
		t, err := time.Parse("2006-01-02", dto.DocumentDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid document date format. Expected YYYY-MM-DD")
		}
		documentDate = pgtype.Date{Time: t, Valid: true}
	}

	// Create update parameters
	updateParams := db.UpdateMedicalRecordByUploaderParams{
		ID:              dto.RecordID.String(),
		UploaderID:      dto.UploaderID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: dto.UploaderAdminID, Valid: true},
		// FilePath:        fileUrl,
		Title:        dto.Title,
		DocumentType: dto.DocumentType,
		DocumentDate: documentDate,
		FileType:     pgtype.Text{String: dto.FileType, Valid: dto.FileType != ""},
		ProviderName: pgtype.Text{String: dto.ProviderName, Valid: dto.ProviderName != ""},
		Specialty:    pgtype.Text{String: dto.Specialty, Valid: dto.Specialty != ""},
		IsShared:     pgtype.Bool{Bool: dto.IsShared, Valid: true},
		SharedUntil:  sharedUntilTime,
	}

	// Only update file path if a new file was uploaded
	if fileUrl != "" {
		updateParams.FilePath = fileUrl
	}

	// Update the record
	record, err := service.recordPort.UpdateMedicalRecord(ctx, updateParams)
	if err != nil {
		cont.Logger().Error("Failed to update medical record:", err)
		switch {
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, utils.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Medical record not found")
		case errors.Is(err, utils.ErrPermissionDenied):
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to update this record")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update medical record")
		}
	}

	return utils.ResponseMessage(http.StatusOK, record, cont)
}
