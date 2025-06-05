package domain

import (
	"github.com/google/uuid"
	"github.com/medicue/adapters/db"
)

type (
	Address struct {
		Street  string `json:"street" validate:"max=50,required"`
		City    string `json:"city" validate:"max=50,required"`
		State   string `json:"state" validate:"max=50,required"`
		Country string `json:"country" validate:"max=50,required"`
	}
	Contact struct {
		Phone []string `json:"phone_numbers" validate:"gt=0,dive,required"`
		Email string   `json:"email" validate:"required,email"`
	}
	CreateDiagnosticDTO struct {
		DiagnosticCentreName string              `json:"diagnostic_centre_name" validate:"gte=10,lte=250,required"`
		Latitude             float64             `json:"latitude" validate:"min=-90.00,max=90.00,required"`
		Longitude            float64             `json:"longitude" validate:"min=-180.00,max=180.00,required"`
		Address              Address             `json:"address"`
		Contact              Contact             `json:"contact"`
		Doctors              []db.Doctor         `json:"doctors"`
		AvailableTests       []db.AvailableTests `json:"available_tests"`
		CreatedBy            uuid.UUID           `json:"created_by"`
		AdminId              uuid.UUID           `json:"admin_id" validate:"uuid,required"`
	}
	GetDiagnosticParamDTO struct {
		DiagnosticCentreID uuid.UUID `json:"diagnostic_centre_id" validate:"uuid,required"`
	}
	UpdateDiagnosticDTO struct {
		DiagnosticCentreName string              `json:"diagnostic_centre_name"`
		Latitude             float64             `json:"latitude"`
		Longitude            float64             `json:"longitude"`
		Address              Address             `json:"address"`
		Contact              Contact             `json:"contact"`
		Doctors              []db.Doctor         `json:"doctors"`
		AvailableTests       []db.AvailableTests `json:"available_tests"`
		CreatedBy            uuid.UUID           `json:"created_by"`
		ADMINID              uuid.UUID           `json:"admin_id" validate:"uuid,required"`
	}
)
