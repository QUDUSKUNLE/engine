package domain

import (
	"time"

	"github.com/google/uuid"
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
		Phone []string `json:"phone" validate:"gt=0,dive,required"`
		Email string   `json:"email" validate:"required,email"`
	}
	CreateDiagnosticDTO struct {
		DiagnosticCentreName string    `json:"diagnostic_centre_name" validate:"gte=10,lte=250,required"`
		Latitude             float64   `json:"latitude" validate:"min=-90.00,max=90.00,required"`
		Longitude            float64   `json:"longitude" validate:"min=-180.00,max=180.00,required"`
		Address              Address   `json:"address"`
		Contact              Contact   `json:"contact"`
		Doctors              []string  `json:"doctors"`
		AvailableTests       []string  `json:"available_tests"`
		CreatedBy            uuid.UUID `json:"created_by"`
		AdminId              uuid.UUID `json:"admin_id" validate:"uuid,required"`
	}
	GetDiagnosticParamDTO struct {
		DiagnosticCentreID string `param:"diagnostic_centre_id" validate:"uuid,required"`
	}
	SearchDiagnosticCentreQueryDTO struct {
		Latitude  float64 `query:"latitude"`
		Longitude float64 `query:"longitude"`
		Doctor    string  `query:"doctor"`
		Test      string  `query:"available_tests"`
		Limit     int32   `query:"limit" validate:"omitempty,gte=0"`
		Offset    int32   `query:"offset" validate:"omitempty,gte=0"`
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
	// UpdateDiagnosticManagerDTO represents the payload for updating a diagnostic centre manager
	UpdateDiagnosticManagerDTO struct {
		ManagerID string `json:"manager_id" validate:"required,uuid"`
	}
	// GetDiagnosticRecordsParamDTO represents query parameters for fetching diagnostic records
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
