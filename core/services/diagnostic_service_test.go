package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (m *MockDiagnosticRepo) GetDiagnosticCentreByOwner(ctx context.Context, params db.Get_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) UpdateDiagnosticCentreByOwner(ctx context.Context, params db.Update_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
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

func (m *MockDiagnosticRepo) DeleteDiagnosticCentreByOwner(ctx context.Context, params db.Delete_Diagnostic_Centre_ByOwnerParams) (*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) ListDiagnosticCentresByOwner(ctx context.Context, params db.List_Diagnostic_Centres_ByOwnerParams) ([]*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) RetrieveDiagnosticCentres(ctx context.Context, params db.Retrieve_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.DiagnosticCentre), args.Error(1)
}

func (m *MockDiagnosticRepo) SearchDiagnosticCentres(ctx context.Context, params db.Search_Diagnostic_CentresParams) ([]*db.DiagnosticCentre, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.DiagnosticCentre), args.Error(1)
}

type MockScheduleRepo struct{ mock.Mock }

func (m *MockScheduleRepo) GetDiagnosticSchedulesByCentre(ctx context.Context, params db.Get_Diagnsotic_Schedules_By_CentreParams) ([]*db.DiagnosticSchedule, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*db.DiagnosticSchedule), args.Error(1)
}

// Tests related to DeleteDiagnosticCentre
func TestDeleteDiagnosticCentre_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	ownerID := uuid.New()
	centreID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/diagnostic-centres/"+centreID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims for owner
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   ownerID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Mock GetDiagnosticCentreByOwner
	DiagnosticRepo.On("GetDiagnosticCentreByOwner", mock.Anything, mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByOwnerParams) bool {
		return params.ID == centreID.String() && params.CreatedBy == ownerID.String()
	})).Return(&db.DiagnosticCentre{
		ID:        centreID.String(),
		CreatedBy: ownerID.String(),
	}, nil)

	// Mock DeleteDiagnosticCentreByOwner
	DiagnosticRepo.On("DeleteDiagnosticCentreByOwner", mock.Anything, mock.MatchedBy(func(params db.Delete_Diagnostic_Centre_ByOwnerParams) bool {
		return params.ID == centreID.String() && params.CreatedBy == ownerID.String()
	})).Return(&db.DiagnosticCentre{
		ID:        centreID.String(),
		CreatedBy: ownerID.String(),
	}, nil)

	// Call service method
	err := h.DeleteDiagnosticCentre(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

func TestDeleteDiagnosticCentre_NotFound(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	ownerID := uuid.New()
	centreID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/diagnostic-centres/"+centreID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims for owner
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   ownerID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Mock GetDiagnosticCentreByOwner to return not found
	DiagnosticRepo.On("GetDiagnosticCentreByOwner", mock.Anything, mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByOwnerParams) bool {
		return params.ID == centreID.String() && params.CreatedBy == ownerID.String()
	})).Return(nil, utils.ErrNotFound)

	// Call service method and expect NotFound error
	err := h.DeleteDiagnosticCentre(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
		assert.Contains(t, httpErr.Message, "diagnostic centre not found")
	} else {
		t.Error("Expected HTTP error with NotFound status")
	}

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

func TestDeleteDiagnosticCentre_Unauthorized(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	centreID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/diagnostic-centres/"+centreID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims with wrong user type
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uuid.New(),
		UserType: db.UserEnumDIAGNOSTICCENTREMANAGER, // Not an owner
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Call service method and expect Unauthorized error
	err := h.DeleteDiagnosticCentre(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	} else {
		t.Error("Expected HTTP error with Unauthorized status")
	}
}

// Tests related to GetDiagnosticCentreStats
func TestGetDiagnosticCentreStats_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	centreID := uuid.New()
	managerID := uuid.New()

	// Create a test diagnostic centre
	testCentre := &db.DiagnosticCentre{
		ID:                   centreID.String(),
		DiagnosticCentreName: "Test Centre",
		Doctors:              []string{"Dr. Smith", "Dr. Jones"},
		AvailableTests:       []string{"Blood Test", "X-Ray", "MRI"},
		CreatedBy:            uuid.New().String(),
		AdminID:              managerID.String(),
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/diagnostic-centres/"+centreID.String()+"/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims for manager
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   managerID,
		UserType: db.UserEnumDIAGNOSTICCENTREMANAGER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Mock GetDiagnosticCentre
	DiagnosticRepo.On("GetDiagnosticCentre", mock.Anything, centreID.String()).Return(testCentre, nil)

	// Call service method
	err := h.GetDiagnosticCentreStats(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response contains expected stats
	assert.Equal(t, testCentre.ID, response["diagnostic_centre_id"])
	assert.Equal(t, testCentre.DiagnosticCentreName, response["diagnostic_centre_name"])
	assert.Equal(t, float64(len(testCentre.Doctors)), response["total_doctors"])
	assert.Equal(t, float64(len(testCentre.AvailableTests)), response["total_tests"])

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

func TestGetDiagnosticCentreStats_NotFound(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	centreID := uuid.New()
	managerID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/diagnostic-centres/"+centreID.String()+"/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims for manager
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   managerID,
		UserType: db.UserEnumDIAGNOSTICCENTREMANAGER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Mock GetDiagnosticCentre to return not found
	DiagnosticRepo.On("GetDiagnosticCentre", mock.Anything, centreID.String()).Return(nil, errors.New("diagnostic centre not found"))

	// Call service method and expect NotFound error
	err := h.GetDiagnosticCentreStats(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	} else {
		t.Error("Expected HTTP error with NotFound status")
	}

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

// Tests related to GetDiagnosticCentresByManager
func TestGetDiagnosticCentresByManager_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	managerID := uuid.New()
	centreID := uuid.New()

	// Create a test diagnostic centre
	testCentre := &db.DiagnosticCentre{
		ID:                   centreID.String(),
		DiagnosticCentreName: "Test Centre",
		AdminID:              managerID.String(),
		Doctors:              []string{"Dr. Smith"},
		AvailableTests:       []string{"Blood Test"},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/diagnostic-centres/manager", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock JWT token claims for manager
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   managerID,
		UserType: db.UserEnumDIAGNOSTICCENTREMANAGER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Set pagination params
	paginationDTO := &domain.PaginationQueryDTO{
		Page:    1,
		PerPage: 10,
	}
	c.Set(utils.ValidatedQueryParamDTO, paginationDTO)

	// Mock GetDiagnosticCentreByManager
	DiagnosticRepo.On("GetDiagnosticCentreByManager", mock.Anything, mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByManagerParams) bool {
		return params.AdminID == managerID.String()
	})).Return(testCentre, nil)

	// Call service method
	err := h.GetDiagnosticCentresByManager(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response []map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response contains the diagnostic centre
	require.Len(t, response, 1)
	assert.Equal(t, testCentre.ID, response[0]["diagnostic_centre_id"])
	assert.Equal(t, testCentre.DiagnosticCentreName, response[0]["diagnostic_centre_name"])

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

func TestGetDiagnosticCentresByManager_Unauthorized(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/diagnostic-centres/manager", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock JWT token claims with wrong user type
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   uuid.New(),
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER, // Not a manager
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Call service method and expect Unauthorized error
	err := h.GetDiagnosticCentresByManager(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	} else {
		t.Error("Expected HTTP error with Unauthorized status")
	}
}

// Tests related to UpdateDiagnosticCentreManager
func TestUpdateDiagnosticCentreManager_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	ownerID := uuid.New()
	centreID := uuid.New()
	newManagerID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/diagnostic-centres/"+centreID.String()+"/manager",
		strings.NewReader(`{"manager_id":"`+newManagerID.String()+`"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims for owner
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   ownerID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Set update DTO
	managerDTO := &domain.UpdateDiagnosticManagerDTO{
		ManagerID: newManagerID.String(),
	}
	c.Set(utils.ValidatedBodyDTO, managerDTO)

	// Mock GetDiagnosticCentreByOwner
	DiagnosticRepo.On("GetDiagnosticCentreByOwner", mock.Anything, mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByOwnerParams) bool {
		return params.ID == centreID.String() && params.CreatedBy == ownerID.String()
	})).Return(&db.DiagnosticCentre{
		ID:        centreID.String(),
		CreatedBy: ownerID.String(),
	}, nil)

	// Mock UpdateDiagnosticCentreByOwner
	DiagnosticRepo.On("UpdateDiagnosticCentreByOwner", mock.Anything, mock.MatchedBy(func(params db.Update_Diagnostic_Centre_ByOwnerParams) bool {
		return params.ID == centreID.String() &&
			params.CreatedBy == ownerID.String() &&
			params.AdminID == newManagerID.String()
	})).Return(&db.DiagnosticCentre{
		ID:        centreID.String(),
		CreatedBy: ownerID.String(),
		AdminID:   newManagerID.String(),
	}, nil)

	// Call service method
	err := h.UpdateDiagnosticCentreManager(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

func TestUpdateDiagnosticCentreManager_NotFound(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	ownerID := uuid.New()
	centreID := uuid.New()
	newManagerID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/diagnostic-centres/"+centreID.String()+"/manager",
		strings.NewReader(`{"manager_id":"`+newManagerID.String()+`"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

	// Mock JWT token claims for owner
	claims := &domain.JwtCustomClaimsDTO{
		UserID:   ownerID,
		UserType: db.UserEnumDIAGNOSTICCENTREOWNER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	// Set update DTO
	managerDTO := &domain.UpdateDiagnosticManagerDTO{
		ManagerID: newManagerID.String(),
	}
	c.Set(utils.ValidatedBodyDTO, managerDTO)

	// Mock GetDiagnosticCentreByOwner to return not found
	DiagnosticRepo.On("GetDiagnosticCentreByOwner", mock.Anything, mock.MatchedBy(func(params db.Get_Diagnostic_Centre_ByOwnerParams) bool {
		return params.ID == centreID.String() && params.CreatedBy == ownerID.String()
	})).Return(nil, utils.ErrNotFound)

	// Call service method and expect NotFound error
	err := h.UpdateDiagnosticCentreManager(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	} else {
		t.Error("Expected HTTP error with NotFound status")
	}

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}

// Tests related to GetDiagnosticCentreRecords
func TestGetDiagnosticCentreRecords_Success(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	RecordRepo := new(MockRecordRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
		RecordRepo:     RecordRepo,
	}

	centreID := uuid.New()
	uploaderID := uuid.New()

	// Create test diagnostic centre
	DiagnosticRepo.On("GetDiagnosticCentre", mock.Anything, centreID.String()).Return(&db.DiagnosticCentre{
		ID:        centreID.String(),
		CreatedBy: uploaderID.String(),
	}, nil)

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

	// Mock pagination params
	paginationDTO := &domain.PaginationQueryDTO{
		Page:    1,
		PerPage: 10,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/diagnostic-centres/"+centreID.String()+"/records", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())
	c.Set(utils.ValidatedQueryParamDTO, paginationDTO)

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

	// Call service method
	err := h.GetDiagnosticCentreRecords(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response struct {
		Records []*db.GetMedicalRecordsRow `json:"records"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response records
	assert.Len(t, response.Records, 2)
	assert.Equal(t, mockRecords[0].ID, response.Records[0].ID)
	assert.Equal(t, mockRecords[1].ID, response.Records[1].ID)

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
	RecordRepo.AssertExpectations(t)
}

func TestGetDiagnosticCentreRecords_NotFound(t *testing.T) {
	DiagnosticRepo := new(MockDiagnosticRepo)
	h := &services.ServicesHandler{
		DiagnosticRepo: DiagnosticRepo,
	}

	centreID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/diagnostic-centres/"+centreID.String()+"/records", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("diagnostic_centre_id")
	c.SetParamValues(centreID.String())

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

	// Mock GetDiagnosticCentre to return not found
	DiagnosticRepo.On("GetDiagnosticCentre", mock.Anything, centreID.String()).Return(nil, errors.New("diagnostic centre not found"))

	// Call service method and expect NotFound error
	err := h.GetDiagnosticCentreRecords(c)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	} else {
		t.Error("Expected HTTP error with NotFound status")
	}

	// Verify mock expectations
	DiagnosticRepo.AssertExpectations(t)
}
