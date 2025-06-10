package domain

import (
	"github.com/google/uuid"
	"github.com/medicue/adapters/db"
)

type (
	CreateScheduleDTO struct {
		UserID                   string                      `json:"user_id"`
		DiagnosticCentreID       uuid.UUID                   `json:"diagnostic_centre_id" validate:"uuid" required:"true"`
		ScheduleTime             string                      `json:"schedule_time" validate:"datetime=2006-01-02T15:04:05.000Z07:00,required"`
		TestType                 db.TestType                 `json:"test_type" validate:"oneof=BLOOD_TEST URINE_TEST X_RAY MRI CT_SCAN ULTRASOUND ECG COVID_TEST DNA_TEST ALLERGY_TEST GENETIC_TEST OTHER EEG BIOPSY SKIN_TEST IMMUNOLOGY_TEST HORMONE_TEST VIRAL_TEST BACTERIAL_TEST PARASITIC_TEST FUNGAL_TEST MOLECULAR_TEST TOXICOLOGY_TEST ECHO COVID_19_TEST BLOOD_SUGAR_TEST LIPID_PROFILE HEMOGLOBIN_TEST THYROID_TEST LIVER_FUNCTION_TEST KIDNEY_FUNCTION_TEST URIC_ACID_TEST VITAMIN_D_TEST VITAMIN_B12_TEST HEMOGRAM COMPLETE_BLOOD_COUNT BLOOD_GROUPING HEPATITIS_B_TEST HEPATITIS_C_TEST HIV_TEST MALARIA_TEST DENGUE_TEST TYPHOID_TEST COVID_19_ANTIBODY_TEST COVID_19_RAPID_ANTIGEN_TEST COVID_19_RT_PCR_TEST PREGNANCY_TEST,required"`
		Doctor                   db.Doctor                   `json:"doctor" validate:"oneof=Male Female,required"`
		Notes                    string                      `json:"notes"`
		ScheduleAcceptanceStatus db.ScheduleAcceptanceStatus `json:"acceptance_status" validate:"omitempty,oneof=PENDING ACCEPTED REJECTED"`
	}
	UpdateScheduleDTO struct {
		UserID             string            `json:"user_id"`
		DiagnosticCentreID uuid.UUID         `json:"diagnostic_centre_id" validate:"uuid"`
		ScheduleTime       string            `json:"schedule_time" validate:"datetime=2006-01-02T15:04:05Z07:00" required:"true"`
		TestType           db.TestType       `json:"test_type" validate:"oneof=BLOOD_TEST URINE_TEST X_RAY MRI CT_SCAN ULTRASOUND ECG COVID_TEST DNA_TEST ALLERGY_TEST GENETIC_TEST OTHER EEG BIOPSY SKIN_TEST IMMUNOLOGY_TEST HORMONE_TEST VIRAL_TEST BACTERIAL_TEST PARASITIC_TEST FUNGAL_TEST MOLECULAR_TEST TOXICOLOGY_TEST ECHO COVID_19_TEST BLOOD_SUGAR_TEST LIPID_PROFILE HEMOGLOBIN_TEST THYROID_TEST LIVER_FUNCTION_TEST KIDNEY_FUNCTION_TEST URIC_ACID_TEST VITAMIN_D_TEST VITAMIN_B12_TEST HEMOGRAM COMPLETE_BLOOD_COUNT BLOOD_GROUPING HEPATITIS_B_TEST HEPATITIS_C_TEST HIV_TEST MALARIA_TEST DENGUE_TEST TYPHOID_TEST COVID_19_ANTIBODY_TEST COVID_19_RAPID_ANTIGEN_TEST COVID_19_RT_PCR_TEST PREGNANCY_TEST"`
		ScheduleStatus     db.ScheduleStatus `json:"schedule_status" validate:"oneof=SCHEDULED CANCELED"`
		Doctor             db.Doctor         `json:"doctor" validate:"oneof=Male Female"`
		Notes              string            `json:"notes"`
	}
	GetDiagnosticScheduleParamDTO struct {
		ScheduleID uuid.UUID `param:"schedule_id" validate:"uuid,required"`
	}
	GetDiagnosticSchedulesByCentreParamDTO struct {
		DiagnosticCentreID uuid.UUID `param:"diagnostic_centre_id" validate:"uuid,required"`
	}
	GetDiagnosticScheduleByCentreParamDTO struct {
		ScheduleID         uuid.UUID `param:"schedule_id" validate:"uuid,required"`
		DiagnosticCentreID uuid.UUID `param:"diagnostic_centre_id" validate:"uuid,required"`
	}
	UpdateDiagnosticScheduleByCentreDTO struct {
		AcceptanceStatus db.ScheduleAcceptanceStatus `json:"acceptance_status" validate:"oneof=PENDING ACCEPTED REJECTED" required:"true"`
	}
)
