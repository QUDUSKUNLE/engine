package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	AppointmentStatusPending     AppointmentStatus = "pending"
	AppointmentStatusConfirmed   AppointmentStatus = "confirmed"
	AppointmentStatusInProgress  AppointmentStatus = "in_progress"
	AppointmentStatusCompleted   AppointmentStatus = "completed"
	AppointmentStatusCancelled   AppointmentStatus = "cancelled"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

type (
	// AppointmentStatus represents the status of an appointment
	AppointmentStatus    string
	CreateAppointmentDTO struct {
		DiagnosticCentreID uuid.UUID `json:"diagnostic_centre_id" validate:"required,uuid"`
		TestType           string    `json:"test_type" validate:"required,oneof=BLOOD_TEST URINE_TEST X_RAY MRI CT_SCAN ULTRASOUND ECG EEG BIOPSY SKIN_TEST IMMUNOLOGY_TEST HORMONE_TEST VIRAL_TEST BACTERIAL_TEST PARASITIC_TEST FUNGAL_TEST MOLECULAR_TEST TOXICOLOGY_TEST ECHO COVID_19_TEST BLOOD_SUGAR_TEST LIPID_PROFILE HEMOGLOBIN_TEST THYROID_TEST LIVER_FUNCTION_TEST KIDNEY_FUNCTION_TEST URIC_ACID_TEST VITAMIN_D_TEST VITAMIN_B12_TEST HEMOGRAM COMPLETE_BLOOD_COUNT BLOOD_GROUPING HEPATITIS_B_TEST HEPATITIS_C_TEST HIV_TEST MALARIA_TEST DENGUE_TEST TYPHOID_TEST COVID_19_ANTIBODY_TEST COVID_19_RAPID_ANTIGEN_TEST COVID_19_RT_PCR_TEST PREGNANCY_TEST ALLERGY_TEST GENETIC_TEST OTHER"`
		AppointmentDate    time.Time `json:"appointment_date" validate:"required"`
		Amount             float64   `json:"amount" validate:"required"`
		PreferredDoctor    string    `json:"preferred_doctor" validate:"omitempty,oneof=Male Female"`
		PaymentProvider    string    `json:"payment_provider" validate:"oneof=PAYSTACK FLUTTERWAVE STRIPE MONNIFY"`
		Notes              string    `json:"notes" validate:"max=500"`
	}
	// CreatePaymentDTO represents the request body for creating a payment
	ConfirmAppointmentDTO struct {
		AppointmentID     string      `json:"appointment_id" validate:"required,uuid"`
		Amount            float64     `json:"amount" validate:"required,gt=0"`
		Currency          string      `json:"currency" validate:"required,len=3"`
		PaymentMethod     string      `json:"payment_method" validate:"required,oneof=card transfer cash wallet"`
		PaymentProvider   string      `json:"payment_provider" validate:"required,oneof=PAYSTACK FLUTTERWAVE STRIPE MONNIFY"`
		PaymentMetadata   interface{} `json:"payment_metadata,omitempty"`
		ProviderReference string      `json:"provider_reference" validate:"required"`
	}
	// GetAppointmentDTO represents the request parameters for getting an appointment
	GetAppointmentDTO struct {
		AppointmentID string `param:"appointment_id" validate:"required,uuid"`
	}
	// ListAppointmentsDTO represents the query parameters for listing appointments
	ListAppointmentsDTO struct {
		DiagnosticCentreID string            `query:"diagnostic_centre_id" validate:"omitempty,uuid"`
		PatientID          string            `query:"patient_id" validate:"omitempty,uuid"`
		Status             AppointmentStatus `query:"status" validate:"omitempty,oneof=pending confirmed in_progress completed cancelled rescheduled"`
		FromDate           time.Time         `query:"from_date" validate:"omitempty"`
		ToDate             time.Time         `query:"to_date" validate:"omitempty,gtefield=FromDate"`
		PaginationQueryDTO
	}

	// UpdateAppointmentDTO represents the request body for updating an appointment
	UpdateAppointmentDTO struct {
		AppointmentID string            `param:"appointment_id" validate:"required,uuid"`
		Status        AppointmentStatus `json:"status" validate:"required,oneof=pending confirmed in_progress completed cancelled rescheduled"`
		Notes         string            `json:"notes" validate:"max=500"`
	}

	// CancelAppointmentDTO represents the request body for cancelling an appointment
	CancelAppointmentDTO struct {
		AppointmentID string `param:"appointment_id" validate:"required,uuid"`
		Reason        string `json:"reason" validate:"required,max=500"`
	}

	// RescheduleAppointmentDTO represents the request body for rescheduling an appointment
	RescheduleAppointmentDTO struct {
		AppointmentID    string    `param:"appointment_id" validate:"required,uuid"`
		NewScheduleID    string    `json:"new_schedule_id" validate:"required,uuid"`
		NewDate          time.Time `json:"new_date" validate:"required"`
		NewTimeSlot      string    `json:"new_time_slot" validate:"required"`
		RescheduleReason string    `json:"reschedule_reason" validate:"required,max=500"`
	}
	GetNotificationsDTO struct {
		PaginationQueryDTO
	}
)
