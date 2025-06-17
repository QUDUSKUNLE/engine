package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

// CreateMedicalRecord handles the creation of a new medical record.
func (service *ServicesHandler) CreateMedicalRecord(context echo.Context) error {
	ctx := context.Request().Context()

	// Authentication & Authorization
	uploader, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	// Build and validate DTO
	dto, err := buildCreateMedicalRecordDto(context)
	if err != nil {
		context.Logger().Error("Failed to build medical record DTO:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}
	if dto.FileUpload.Content == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File upload is required")
	}

	// Set uploader info
	dto.UploaderType = uploader.UserType
	dto.UploaderAdminID = uploader.UserID

	// Validate uploader_admin_id and uploader_id before uploading data to cloud
	params := db.Get_Diagnostic_Centre_ByManagerParams{
		ID:      dto.UploaderID.String(),
		AdminID: dto.UploaderAdminID.String(),
	}
	_, err = service.DiagnosticRepo.GetDiagnosticCentreByManager(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	// Upload file to cloud
	fileUrl, err := service.FileRepo.UploadFile(ctx, dto.FileUpload.Content)
	if err != nil {
		context.Logger().Error("File upload failed:", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Failed to upload file")
	}

	// Parse SharedUntil string to time.Time
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

	// Create medical record
	createParams := db.CreateMedicalRecordParams{
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
	}

	record, err := service.RecordRepo.CreateMedicalRecord(ctx, createParams)
	if err != nil {
		context.Logger().Error("Failed to create medical record:", err)
		switch {
		case errors.Is(err, utils.ErrDatabaseError):
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error occurred")
		case errors.Is(err, utils.ErrInvalidInput):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid record data")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create medical record")
		}
	}

	return utils.ResponseMessage(http.StatusCreated, record, context)
}

// GetMedicalRecord retrieves a single medical record.
func (service *ServicesHandler) GetMedicalRecord(context echo.Context) error {
	// Authentication check
	user, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, utils.AuthenticationRequired)
	}

	// This validated at the middleware level
	param, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetMedicalRecordParamsDTO)

	// Fetch medical record
	response, err := service.RecordRepo.GetMedicalRecord(
		context.Request().Context(),
		db.GetMedicalRecordParams{
			ID:     param.RecordID.String(),
			UserID: user.UserID.String(),
		},
	)
	if err != nil {
		context.Logger().Error("Failed to get medical record:", err)
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

	return utils.ResponseMessage(http.StatusOK, response, context)
}

// GetMedicalRecords retrieves multiple medical records for a user.
func (service *ServicesHandler) GetMedicalRecords(context echo.Context) error {
	user, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumUSER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// This validated at the middleware level
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetMedicalRecordsParamQueryDTO)

	query = SetDefaultPagination(query).(*domain.GetMedicalRecordsParamQueryDTO)

	response, err := service.RecordRepo.GetMedicalRecords(context.Request().Context(), db.GetMedicalRecordsParams{
		UserID: user.UserID.String(),
		Limit:  query.GetLimit(),
		Offset: query.GetOffset(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, context)
	}

	return utils.ResponseMessage(http.StatusOK, response, context)
}

// GetUploaderMedicalRecord retrieves a single medical record uploaded by a specific uploader.
func (service *ServicesHandler) GetUploaderMedicalRecord(context echo.Context) error {
	uploader, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}
	// This validated at the middleware level
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetUploaderMedicalRecordParamsDTO)

	response, err := service.RecordRepo.GetUploaderMedicalRecord(context.Request().Context(), db.GetUploaderMedicalRecordParams{
		ID:              query.RecordID.String(),
		UploaderID:      query.UploaderID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: uploader.UserID, Valid: true},
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, response, context)
}

// GetUploaderMedicalRecords retrieves multiple medical records uploaded by a specific uploader.
func (service *ServicesHandler) GetUploaderMedicalRecords(context echo.Context) error {
	_, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	// This validated at the middleware level
	query, _ := context.Get(utils.ValidatedQueryParamDTO).(*domain.GetUploaderMedicalRecordsParamQueryDTO)

	query = SetDefaultPagination(query).(*domain.GetUploaderMedicalRecordsParamQueryDTO)

	response, err := service.RecordRepo.GetUploaderMedicalRecords(context.Request().Context(), db.GetUploaderMedicalRecordsParams{
		UploaderID: query.UploaderID.String(),
		Limit:      query.GetLimit(),
		Offset:     query.GetOffset(),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, context)
	}
	if len(response) == 0 {
		return utils.ResponseMessage(http.StatusOK, []interface{}{}, context)
	}

	return utils.ResponseMessage(http.StatusOK, response, context)
}

// UpdateMedicalRecord updates an existing medical record by the uploader.
func (service *ServicesHandler) UpdateMedicalRecord(context echo.Context) error {
	ctx := context.Request().Context()

	// Authentication & Authorization
	uploader, err := utils.PrivateMiddlewareContext(context, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	// Get the validated DTO from context
	dto, ok := context.Get(utils.ValidatedBodyDTO).(*domain.UpdateMedicalRecordDTO)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	// Set uploader admin ID from the authenticated user
	dto.UploaderAdminID = uploader.UserID

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
		fileUrl, err = service.FileRepo.UploadFile(ctx, dto.FileUpload.Content)
		if err != nil {
			context.Logger().Error("File upload failed:", err)
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
		Title:           dto.Title,
		DocumentType:    dto.DocumentType,
		DocumentDate:    documentDate,
		FileType:        pgtype.Text{String: dto.FileType, Valid: dto.FileType != ""},
		ProviderName:    pgtype.Text{String: dto.ProviderName, Valid: dto.ProviderName != ""},
		Specialty:       pgtype.Text{String: dto.Specialty, Valid: dto.Specialty != ""},
		IsShared:        pgtype.Bool{Bool: dto.IsShared, Valid: true},
		SharedUntil:     sharedUntilTime,
	}

	// Only update file path if a new file was uploaded
	if fileUrl != "" {
		updateParams.FilePath = fileUrl
	}

	// Update the record
	record, err := service.RecordRepo.UpdateMedicalRecord(ctx, updateParams)
	if err != nil {
		context.Logger().Error("Failed to update medical record:", err)
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

	return utils.ResponseMessage(http.StatusOK, record, context)
}
