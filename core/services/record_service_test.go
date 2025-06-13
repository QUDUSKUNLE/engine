package services_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/services"
	"github.com/medicue/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDiagnosticRepo struct{ mock.Mock }

func (m *MockDiagnosticRepo) GetDiagnosticCentre(ctx context.Context, id string) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

// Add stub for UpdateDiagnosticCentreByOwner to satisfy the interface
func (m *MockDiagnosticRepo) UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

// Add stub for GetNearestDiagnosticCentres to satisfy the interface
func (m *MockDiagnosticRepo) GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.Get_Nearest_Diagnostic_CentresRow), args.Error(1)
}

func (m *MockDiagnosticRepo) GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) CreateDiagnosticCentre(ctx context.Context, params db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

type MockFileRepo struct{ mock.Mock }

func (m *MockFileRepo) UploadFile(ctx context.Context, filePath []byte) (string, error) {
	args := m.Called(ctx, filePath)
	return args.String(0), args.Error(1)
}

func (m *MockFileRepo) SaveFile(filePath string, data []byte) error {
	args := m.Called(filePath, data)
	return args.Error(0)
}

func (m *MockFileRepo) ExtractTextFromImage(filePath string) (map[string]interface{}, error) {
	args := m.Called(filePath)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockFileRepo) ExtractTextWithMonai(filePath string, someParam string) (string, error) {
	args := m.Called(filePath, someParam)
	return args.String(0), args.Error(1)
}

// Add DeleteFile method to satisfy ports.FileService interface
func (m *MockFileRepo) DeleteFile(filePath string) error {
	args := m.Called(filePath)
	return args.Error(0)
}

type MockRecordRepo struct{ mock.Mock }

func (m *MockRecordRepo) CreateMedicalRecord(ctx context.Context, params db.CreateMedicalRecordParams) (*db.MedicalRecord, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.MedicalRecord), args.Error(1)
}

func (m *MockRecordRepo) UpdateMedicalRecord(ctx context.Context, params db.UpdateMedicalRecordByUploaderParams) (*db.MedicalRecord, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.MedicalRecord), args.Error(1)
}

func (m *MockRecordRepo) GetMedicalRecord(ctx context.Context, params db.GetMedicalRecordParams) (*db.GetMedicalRecordRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.GetMedicalRecordRow), args.Error(1)
}

func (m *MockRecordRepo) GetMedicalRecords(ctx context.Context, params db.GetMedicalRecordsParams) ([]*db.GetMedicalRecordsRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.GetMedicalRecordsRow), args.Error(1)
}

func (m *MockRecordRepo) GetUploaderMedicalRecord(ctx context.Context, params db.GetUploaderMedicalRecordParams) (*db.GetUploaderMedicalRecordRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.GetUploaderMedicalRecordRow), args.Error(1)
}

func (m *MockRecordRepo) GetUploaderMedicalRecords(ctx context.Context, params db.GetUploaderMedicalRecordsParams) ([]*db.GetUploaderMedicalRecordsRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.GetUploaderMedicalRecordsRow), args.Error(1)
}

func TestCreateMedicalRecord_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	FileRepo := new(MockFileRepo)
	RecordRepo := new(MockRecordRepo)

	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
		FileRepo:       FileRepo,
		RecordRepo:     RecordRepo,
	}

	uploaderID := "123e4567-e89b-12d3-a456-426614174000"
	uploaderAdminID := "123e4567-e89b-12d3-a456-426614174001"
	scheduleID := "123e4567-e89b-12d3-a456-426614174002"

	// Convert string IDs to uuid.UUID
	userUUID, _ := uuid.Parse(uploaderID)
	uploaderUUID, _ := uuid.Parse(uploaderID)
	uploaderAdminUUID, _ := uuid.Parse(uploaderAdminID)
	scheduleUUID, _ := uuid.Parse(scheduleID)

	dto := &domain.CreateMedicalRecordDTO{
		UserID:          userUUID,
		UploaderID:      uploaderUUID,
		UploaderAdminID: uploaderAdminUUID,
		UploaderType:    db.UserEnum("DIAGNOSTIC_CENTRE_MANAGER"),
		ScheduleID:      scheduleUUID,
		Title:           "Test Medical Record Title",
		DocumentType:    db.DocumentType("LAB_REPORT"),
		FileUpload: domain.File{
			FileName: "test.pdf",
			FileSize: 1234,
			Content:  []byte("filecontent"),
		},
		FileType:     "pdf",
		UploadedAt:   time.Now().Format(time.RFC3339),
		ProviderName: "Test Provider",
		Specialty:    "Cardiology",
		IsShared:     true,
		SharedUntil:  "2025-06-12T23:17:08.61742+01:00",
	}

	// Create test context with proper multipart form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("user_id", dto.UserID.String())
	_ = writer.WriteField("uploader_id", dto.UploaderID.String())
	_ = writer.WriteField("uploader_admin_id", dto.UploaderAdminID.String())
	_ = writer.WriteField("schedule_id", dto.ScheduleID.String())
	_ = writer.WriteField("title", dto.Title)
	_ = writer.WriteField("document_type", string(dto.DocumentType))
	_ = writer.WriteField("provider_name", dto.ProviderName)
	_ = writer.WriteField("specialty", dto.Specialty)
	_ = writer.WriteField("is_shared", "true")
	_ = writer.WriteField("shared_until", dto.SharedUntil)

	// Add file with the correct field name 'file'
	part, _ := writer.CreateFormFile("file", dto.FileUpload.FileName)
	_, _ = part.Write(dto.FileUpload.Content)
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/records", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock JWT token claims
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   dto.UploaderAdminID,
		UserType: db.UserEnumDIAGNOSTICCENTREMANAGER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)
	c.Set(utils.ValidatedBodyDTO, dto)

	// Mock DiagnosticRepo
	DiagnosticRepo.On("GetDiagnosticCentreByManager",
		mock.Anything,
		mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByManagerParams) bool {
			return params.ID == dto.UploaderID.String() &&
				params.AdminID == dto.UploaderAdminID.String()
		}),
	).Return(&db.DiagnosticCentre{}, nil)

	// Mock FileRepo
	FileRepo.On("UploadFile", mock.Anything, []byte("filecontent")).Return("https://cloud/file.pdf", nil)

	// Parse the SharedUntil time
	sharedUntilTime, _ := time.Parse(time.RFC3339, dto.SharedUntil)
	uploadedAtTime, _ := time.Parse(time.RFC3339, dto.UploadedAt)

	// Set up expectation for CreateMedicalRecord
	expectedParams := db.CreateMedicalRecordParams{
		UserID:          dto.UserID.String(),
		UploaderID:      dto.UploaderID.String(),
		UploaderAdminID: pgtype.UUID{Bytes: uploaderAdminUUID, Valid: true},
		UploaderType:    dto.UploaderType,
		ScheduleID:      dto.ScheduleID.String(),
		Title:           dto.Title,
		DocumentType:    dto.DocumentType,
		FilePath:        "https://cloud/file.pdf",
		FileType:        pgtype.Text{String: dto.FileType, Valid: true},
		UploadedAt:      pgtype.Timestamp{Time: uploadedAtTime, Valid: true},
		ProviderName:    pgtype.Text{String: dto.ProviderName, Valid: true},
		Specialty:       pgtype.Text{String: dto.Specialty, Valid: true},
		IsShared:        pgtype.Bool{Bool: dto.IsShared, Valid: true},
		SharedUntil:     pgtype.Timestamp{Time: sharedUntilTime, Valid: true},
	}

	mockMedicalRecord := &db.MedicalRecord{
		ID:           "test-id",
		UserID:       dto.UserID.String(),
		UploaderID:   dto.UploaderID.String(),
		UploaderType: dto.UploaderType,
		ScheduleID:   dto.ScheduleID.String(),
		Title:        dto.Title,
		DocumentType: dto.DocumentType,
		FilePath:     "https://cloud/file.pdf",
		FileType:     pgtype.Text{String: dto.FileType, Valid: true},
		ProviderName: pgtype.Text{String: dto.ProviderName, Valid: true},
		Specialty:    pgtype.Text{String: dto.Specialty, Valid: true},
		IsShared:     pgtype.Bool{Bool: dto.IsShared, Valid: true},
	}

	// Mock RecordRepo with better match validation
	RecordRepo.On("CreateMedicalRecord", mock.Anything, mock.MatchedBy(func(params db.CreateMedicalRecordParams) bool {
		// Full validation of all fields
		return params.UserID == expectedParams.UserID &&
			params.UploaderID == expectedParams.UploaderID &&
			params.UploaderAdminID.Bytes == expectedParams.UploaderAdminID.Bytes &&
			params.UploaderType == expectedParams.UploaderType &&
			params.ScheduleID == expectedParams.ScheduleID &&
			params.Title == expectedParams.Title &&
			params.DocumentType == expectedParams.DocumentType &&
			params.FilePath == expectedParams.FilePath &&
			params.FileType.String == expectedParams.FileType.String &&
			params.FileType.Valid == expectedParams.FileType.Valid &&
			params.ProviderName.String == expectedParams.ProviderName.String &&
			params.ProviderName.Valid == expectedParams.ProviderName.Valid &&
			params.Specialty.String == expectedParams.Specialty.String &&
			params.Specialty.Valid == expectedParams.Specialty.Valid &&
			params.IsShared.Bool == expectedParams.IsShared.Bool &&
			params.IsShared.Valid == expectedParams.IsShared.Valid &&
			params.SharedUntil.Time.Equal(expectedParams.SharedUntil.Time) &&
			params.SharedUntil.Valid == expectedParams.SharedUntil.Valid
	})).Return(mockMedicalRecord, nil)

	// Call the service method
	err := h.CreateMedicalRecord(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify all mocks
	DiagnosticRepo.AssertExpectations(t)
	FileRepo.AssertExpectations(t)
	RecordRepo.AssertExpectations(t)
}

func TestCreateMedicalRecord_DiagnosticCentreNotFound(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	FileRepo := new(MockFileRepo)
	RecordRepo := new(MockRecordRepo)

	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
		FileRepo:       FileRepo,
		RecordRepo:     RecordRepo,
	}

	// Set up test data
	userID := "123e4567-e89b-12d3-a456-426614174000"
	uploaderID := "123e4567-e89b-12d3-a456-426614174001"
	scheduleID := "123e4567-e89b-12d3-a456-426614174002"

	userUUID, _ := uuid.Parse(userID)
	uploaderUUID, _ := uuid.Parse(uploaderID)
	scheduleUUID, _ := uuid.Parse(scheduleID)

	// Create test context with proper multipart form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("user_id", userUUID.String())
	_ = writer.WriteField("uploader_id", uploaderUUID.String())
	_ = writer.WriteField("schedule_id", scheduleUUID.String())
	_ = writer.WriteField("title", "Test Medical Record Title")
	_ = writer.WriteField("document_type", string(db.DocumentType("LAB_REPORT")))

	// Add file with correct field name 'file'
	part, _ := writer.CreateFormFile("file", "test.pdf")
	_, _ = part.Write([]byte("filecontent"))
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/records", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create DTO without uploader fields - these should be set from JWT token
	dto := &domain.CreateMedicalRecordDTO{
		UserID:       userUUID,
		UploaderID:   uploaderUUID,
		ScheduleID:   scheduleUUID,
		Title:        "Test Medical Record Title",
		DocumentType: db.DocumentType("LAB_REPORT"),
		FileUpload: domain.File{
			FileName: "test.pdf",
			FileSize: 1234,
			Content:  []byte("filecontent"),
		},
	}
	c.Set(utils.ValidatedBodyDTO, dto)

	// Mock JWT token claims
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uploaderUUID,
		UserType: db.UserEnumDIAGNOSTICCENTREMANAGER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Setup mock to return error for diagnostic centre not found
	DiagnosticRepo.On("GetDiagnosticCentreByManager",
		mock.Anything,
		mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByManagerParams) bool {
			return params.ID == uploaderUUID.String() && params.AdminID == uploaderUUID.String()
		}),
	).Return(nil, errors.New("diagnostic centre not found"))

	// Call service method and expect a Forbidden status
	err := h.CreateMedicalRecord(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		t.Logf("Got HTTP error with code: %d and message: %v", httpErr.Code, httpErr.Message)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	} else {
		t.Errorf("Expected HTTP error with Forbidden status, got: %v", err)
	}

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}
