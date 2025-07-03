package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/medivue/adapters/db"
)

type (
	// PaginationParams interface defines pagination behavior
	PaginationParams interface {
		GetLimit() int32
		GetOffset() int32
		SetLimit(limit int32)
		SetPage(page int32)
	}

	Address struct {
		Street  string `json:"street" validate:"max=250,required"`
		City    string `json:"city" validate:"max=50,required"`
		State   string `json:"state" validate:"max=50,required"`
		Country string `json:"country" validate:"max=50,required"`
	}
	Contact struct {
		Phone []string `json:"phone" validate:"required,min_one,dive"`
		Email string   `json:"email" validate:"required,email"`
	}
	AvailableTest struct {
		AvailableTest string  `json:"available_test" validate:"required,oneof=BLOOD_TEST URINE_TEST X_RAY MRI CT_SCAN ULTRASOUND ECG EEG BIOPSY SKIN_TEST ALLERGY_TEST GENETIC_TEST IMMUNOLOGY_TEST HORMONE_TEST VIRAL_TEST BACTERIAL_TEST PARASITIC_TEST FUNGAL_TEST MOLECULAR_TEST TOXICOLOGY_TEST ECHO COVID_19_TEST OTHER BLOOD_SUGAR_TEST LIPID_PROFILE HEMOGLOBIN_TEST THYROID_TEST LIVER_FUNCTION_TEST KIDNEY_FUNCTION_TEST URIC_ACID_TEST VITAMIN_D_TEST VITAMIN_B12_TEST HEMOGRAM COMPLETE_BLOOD_COUNT BLOOD_GROUPING HEPATITIS_B_TEST HEPATITIS_C_TEST HIV_TEST MALARIA_TEST DENGUE_TEST TYPHOID_TEST COVID_19_ANTIBODY_TEST COVID_19_RAPID_ANTIGEN_TEST COVID_19_RT_PCR_TEST PREGNANCY_TEST"`
		TestPrice     float64 `json:"test_price" validate:"required,min=500.00"`
	}
	CreateDiagnosticDTO struct {
		DiagnosticCentreName string          `json:"diagnostic_centre_name" validate:"required,gte=10,lte=250"`
		Latitude             float64         `json:"latitude" validate:"required,min=-90.00,max=90.00"`
		Longitude            float64         `json:"longitude" validate:"required,min=-180.00,max=180.00"`
		Address              Address         `json:"address" validate:"required"`
		Contact              Contact         `json:"contact" validate:"required"`
		Doctors              []db.Doctor     `json:"doctors" validate:"required,min_one,dive,oneof=Male Female"`
		AvailableTests       []AvailableTest `json:"available_tests" validate:"required,min_one,dive"`
		CreatedBy            uuid.UUID       `json:"created_by"`
		AdminId              uuid.UUID       `json:"admin_id" validate:"required,uuid"`
	}
	TestPrices struct {
		TestType string  `json:"test_type"`
		Price    float64 `json:"price"`
	}
	GetDiagnosticParamDTO struct {
		DiagnosticCentreID string `param:"diagnostic_centre_id" validate:"required,uuid"`
	}
	SearchDiagnosticCentreQueryDTO struct {
		Latitude  float64 `query:"latitude"`
		Longitude float64 `query:"longitude"`
		Doctor    string  `query:"doctor"`
		Test      string  `query:"available_tests"`
		Limit     int32   `query:"limit" validate:"omitempty,gte=0"`
		Offset    int32   `query:"offset" validate:"omitempty,gte=0"`
		// Get Availability
		DayOfWeek string `query:"day_of_week" validate:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	}
	UpdateDiagnosticParamDTO struct {
		DiagnosticCentreID uuid.UUID `param:"diagnostic_centre_id"`
	}
	UpdateDiagnosticBodyDTO struct {
		DiagnosticCentreName string    `json:"diagnostic_centre_name"`
		Latitude             float64   `json:"latitude"`
		Longitude            float64   `json:"longitude"`
		Address              Address   `json:"address"`
		Contact              Contact   `json:"contact"`
		Doctors              []string  `json:"doctors"`
		AvailableTests       []string  `json:"available_tests"`
		CreatedBy            uuid.UUID `json:"created_by"`
		ADMINID              uuid.UUID `json:"admin_id" validate:"uuid,required"`
	}
	PaginationQueryDTO struct {
		Page    int32 `query:"page" validate:"omitempty,min=1" json:"page"`
		PerPage int32 `query:"per_page" validate:"omitempty,min=1,max=100" json:"per_page"`
	}
	UpdateDiagnosticManagerDTO struct {
		ManagerID string `json:"manager_id" validate:"required,uuid"`
	}
	GetDiagnosticRecordsParamDTO struct {
		DiagnosticCentreID string    `param:"diagnostic_centre_id" validate:"required,uuid"`
		StartDate          time.Time `query:"start_date" validate:"omitempty" time_format:"2006-01-02"`
		EndDate            time.Time `query:"end_date" validate:"omitempty" time_format:"2006-01-02"`
		DocumentType       string    `query:"document_type" validate:"omitempty,oneof=LAB_REPORT PRESCRIPTION IMAGING DISCHARGE_SUMMARY OTHER"`
		PaginationQueryDTO
	}
)

// SetLimit sets the per page limit
func (p *PaginationQueryDTO) SetLimit(limit int32) {
	p.PerPage = limit
}

// SetPage sets the page number
func (p *PaginationQueryDTO) SetPage(page int32) {
	p.Page = page
}

// SetOffset sets the offset directly
func (p *PaginationQueryDTO) SetOffset(offset int32) {
	p.Page = (offset / p.GetLimit()) + 1
}

// GetLimit returns the limit for pagination
func (p *PaginationQueryDTO) GetLimit() int32 {
	if p.PerPage == 0 {
		return 10 // default limit
	}
	return p.PerPage
}

// GetOffset returns the offset for pagination
func (p *PaginationQueryDTO) GetOffset() int32 {
	if p.Page == 0 {
		return 0
	}
	return (p.Page - 1) * p.GetLimit()
}
