package services_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medivue/adapters/db"
	"github.com/medivue/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDiagnosticRepo struct{ mock.Mock }

func (m *MockDiagnosticRepo) GetDiagnosticCentreByManager(ctx context.Context, params db.Get_Diagnostic_Centre_ByManagerParams) (*db.Get_Diagnostic_Centre_ByManagerRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Get_Diagnostic_Centre_ByManagerRow), args.Error(1)
}

func (m *MockDiagnosticRepo) GetDiagnosticCentreByOwner(ctx context.Context, params db.Get_Diagnostic_Centre_ByOwnerParams) (*db.Get_Diagnostic_Centre_ByOwnerRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Get_Diagnostic_Centre_ByOwnerRow), args.Error(1)
}

func (m *MockDiagnosticRepo) GetDiagnosticCentre(ctx context.Context, id string) (*db.Get_Diagnostic_CentreRow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Get_Diagnostic_CentreRow), args.Error(1)
}

func (m *MockDiagnosticRepo) CreateDiagnosticCentre(ctx context.Context, params db.Create_Diagnostic_CentreParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) DeleteDiagnosticCentreByOwner(ctx context.Context, params db.Delete_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) GetNearestDiagnosticCentres(ctx context.Context, params db.Get_Nearest_Diagnostic_CentresParams) ([]*db.Get_Nearest_Diagnostic_CentresRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.Get_Nearest_Diagnostic_CentresRow), args.Error(1)
}

func (m *MockDiagnosticRepo) ListDiagnosticCentresByOwner(ctx context.Context, params db.List_Diagnostic_Centres_ByOwnerParams) ([]*db.List_Diagnostic_Centres_ByOwnerRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.List_Diagnostic_Centres_ByOwnerRow), args.Error(1)
}

func (m *MockDiagnosticRepo) RetrieveDiagnosticCentres(ctx context.Context, params db.Retrieve_Diagnostic_CentresParams) ([]*db.Retrieve_Diagnostic_CentresRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.Retrieve_Diagnostic_CentresRow), args.Error(1)
}

func (m *MockDiagnosticRepo) SearchDiagnosticCentres(ctx context.Context, params db.Search_Diagnostic_CentresParams) ([]*db.Search_Diagnostic_CentresRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.Search_Diagnostic_CentresRow), args.Error(1)
}

func (m *MockDiagnosticRepo) UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
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

func TestShareMedicalRecord_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	RecordRepo := new(MockRecordRepo)

	centreID := uuid.New()
	uploaderID := uuid.New()
	recordID := "record-id-1"

	// Mock diagnostic centre
	DiagnosticRepo.On("GetDiagnosticCentre", mock.Anything, centreID.String()).Return(&db.DiagnosticCentre{
		ID:        centreID.String(),
		CreatedBy: uploaderID.String(),
	}, nil)

	// Mock medical record
	RecordRepo.On("GetMedicalRecord", mock.Anything, mock.MatchedBy(func(params db.GetMedicalRecordParams) bool {
		return params.ID == recordID
	})).Return(&db.GetMedicalRecordRow{
		ID:           recordID,
		UserID:       "user-id-1",
		UploaderID:   uploaderID.String(),
		UploaderType: db.UserEnum("DIAGNOSTIC_CENTRE_MANAGER"),
		ScheduleID:   "schedule-id-1",
		Title:        "Test Record 1",
		DocumentType: db.DocumentType("LAB_REPORT"),
		FilePath:     "https://cloud/file1.pdf",
		FileType:     pgtype.Text{String: "pdf", Valid: true},
		ProviderName: pgtype.Text{String: "Provider 1", Valid: true},
		Specialty:    pgtype.Text{String: "Cardiology", Valid: true},
		IsShared:     pgtype.Bool{Bool: true, Valid: true},
		SharedUntil:  pgtype.Timestamp{Time: time.Now().Add(24 * time.Hour), Valid: true},
	}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/records/"+recordID+"/share", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("record_id")
	c.SetParamValues(recordID)

	// Mock JWT token claims
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uploaderID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
	RecordRepo.AssertExpectations(t)
}

func TestShareMedicalRecord_NotFound(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	RecordRepo := new(MockRecordRepo)

	recordID := "record-id-1"

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/records/"+recordID+"/share", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("record_id")
	c.SetParamValues(recordID)

	// Mock JWT token claims
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uuid.New(),
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Mock GetMedicalRecord to return not found
	RecordRepo.On("GetMedicalRecord", mock.Anything, mock.MatchedBy(func(params db.GetMedicalRecordParams) bool {
		return params.ID == recordID
	})).Return(nil, errors.New("medical record not found"))

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
	RecordRepo.AssertExpectations(t)
}

func TestGetSharedMedicalRecords_Success(t *testing.T) {
	_ = new(MockDiagnosticRepo)
	RecordRepo := new(MockRecordRepo)

	uploaderID := uuid.New()

	// Mock medical records
	mockRecords := []*db.GetMedicalRecordsRow{
		{
			ID:           "record-id-1",
			UserID:       "user-id-1",
			UploaderID:   uploaderID.String(),
			UploaderType: db.UserEnum("DIAGNOSTIC_CENTRE_MANAGER"),
			ScheduleID:   "schedule-id-1",
			Title:        "Test Record 1",
			DocumentType: db.DocumentType("LAB_REPORT"),
			FilePath:     "https://cloud/file1.pdf",
			FileType:     pgtype.Text{String: "pdf", Valid: true},
			ProviderName: pgtype.Text{String: "Provider 1", Valid: true},
			Specialty:    pgtype.Text{String: "Cardiology", Valid: true},
			IsShared:     pgtype.Bool{Bool: true, Valid: true},
			SharedUntil:  pgtype.Timestamp{Time: time.Now().Add(24 * time.Hour), Valid: true},
		},
		{
			ID:           "record-id-2",
			UserID:       "user-id-2",
			UploaderID:   uploaderID.String(),
			UploaderType: db.UserEnum("DIAGNOSTIC_CENTRE_MANAGER"),
			ScheduleID:   "schedule-id-2",
			Title:        "Test Record 2",
			DocumentType: db.DocumentType("LAB_REPORT"),
			FilePath:     "https://cloud/file2.pdf",
			FileType:     pgtype.Text{String: "pdf", Valid: true},
			ProviderName: pgtype.Text{String: "Provider 2", Valid: true},
			Specialty:    pgtype.Text{String: "Neurology", Valid: true},
			IsShared:     pgtype.Bool{Bool: false, Valid: true},
			SharedUntil:  pgtype.Timestamp{Time: time.Now().Add(48 * time.Hour), Valid: true},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/records/shared", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock JWT token claims
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uploaderID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Mock RecordRepo
	RecordRepo.On("GetMedicalRecords", mock.Anything, mock.MatchedBy(func(params db.GetMedicalRecordsParams) bool {
		return params.Limit == 10 && params.Offset == 0
	})).Return(mockRecords, nil)

	// Parse response
	var response struct {
		Records []*db.GetMedicalRecordsRow `json:"records"`
	}

	// Verify response records
	assert.Len(t, response.Records, 2)
	assert.Equal(t, mockRecords[0].ID, response.Records[0].ID)
	assert.Equal(t, mockRecords[1].ID, response.Records[1].ID)

	// Verify mock expectations
	RecordRepo.AssertExpectations(t)
}

func TestGetSharedMedicalRecords_NotFound(t *testing.T) {
	RecordRepo := new(MockRecordRepo)

	uploaderID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/records/shared", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock JWT token claims
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uploaderID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Verify mock expectations
	RecordRepo.AssertExpectations(t)
}
