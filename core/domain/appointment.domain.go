package domain

import (
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	AppointmentStatusPending     AppointmentStatus = "pending"
	AppointmentStatusConfirmed   AppointmentStatus = "confirmed"
	AppointmentStatusInProgress  AppointmentStatus = "in_progress"
	AppointmentStatusCompleted   AppointmentStatus = "completed"
	AppointmentStatusCancelled   AppointmentStatus = "cancelled"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

// First, let's define TimeSlot as its own type with clear validation
type TimeSlot string

const (
	TimeSlot0800 TimeSlot = "08:00"
	TimeSlot0830 TimeSlot = "08:30"
	TimeSlot0900 TimeSlot = "09:00"
	TimeSlot0930 TimeSlot = "09:30"
	TimeSlot1000 TimeSlot = "10:00"
	TimeSlot1030 TimeSlot = "10:30"
	TimeSlot1100 TimeSlot = "11:00"
	TimeSlot1130 TimeSlot = "11:30"
	TimeSlot1200 TimeSlot = "12:00"
	TimeSlot1230 TimeSlot = "12:30"
	TimeSlot1300 TimeSlot = "13:00"
	TimeSlot1330 TimeSlot = "13:30"
	TimeSlot1400 TimeSlot = "14:00"
	TimeSlot1430 TimeSlot = "14:30"
	TimeSlot1500 TimeSlot = "15:00"
	TimeSlot1530 TimeSlot = "15:30"
	TimeSlot1600 TimeSlot = "16:00"
	TimeSlot1630 TimeSlot = "16:30"
	TimeSlot1700 TimeSlot = "17:00"
)

// CreateAppointmentDTO represents the request body for creating an appointment
type CreateAppointmentDTO struct {
	DiagnosticCentreID uuid.UUID   `json:"diagnostic_centre_id" validate:"required,uuid"`
	TestType           string `json:"test_type" validate:"required,oneof=BLOOD_TEST URINE_TEST X_RAY MRI CT_SCAN ULTRASOUND ECG EEG BIOPSY SKIN_TEST IMMUNOLOGY_TEST HORMONE_TEST VIRAL_TEST BACTERIAL_TEST PARASITIC_TEST FUNGAL_TEST MOLECULAR_TEST TOXICOLOGY_TEST ECHO COVID_19_TEST BLOOD_SUGAR_TEST LIPID_PROFILE HEMOGLOBIN_TEST THYROID_TEST LIVER_FUNCTION_TEST KIDNEY_FUNCTION_TEST URIC_ACID_TEST VITAMIN_D_TEST VITAMIN_B12_TEST HEMOGRAM COMPLETE_BLOOD_COUNT BLOOD_GROUPING HEPATITIS_B_TEST HEPATITIS_C_TEST HIV_TEST MALARIA_TEST DENGUE_TEST TYPHOID_TEST COVID_19_ANTIBODY_TEST COVID_19_RAPID_ANTIGEN_TEST COVID_19_RT_PCR_TEST PREGNANCY_TEST ALLERGY_TEST GENETIC_TEST OTHER"`
	AppointmentDate    time.Time   `json:"appointment_date" validate:"required"`
	// TimeSlot           TimeSlot  `json:"time_slot" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	PreferredDoctor string  `json:"preferred_doctor" validate:"omitempty,oneof=Male Female"`
	PaymentProvider string  `json:"payment_provider" validate:"oneof=PAYSTACK FLUTTERWAVE STRIPE MONNIFY"`
	Notes           string  `json:"notes" validate:"max=500"`
}

// GetAppointmentDTO represents the request parameters for getting an appointment
type GetAppointmentDTO struct {
	AppointmentID string `param:"appointment_id" validate:"required,uuid"`
}

// ListAppointmentsDTO represents the query parameters for listing appointments
type ListAppointmentsDTO struct {
	DiagnosticCentreID string            `query:"diagnostic_centre_id" validate:"omitempty,uuid"`
	PatientID          string            `query:"patient_id" validate:"omitempty,uuid"`
	Status             AppointmentStatus `query:"status" validate:"omitempty,oneof=pending confirmed in_progress completed cancelled rescheduled"`
	FromDate           time.Time         `query:"from_date" validate:"omitempty"`
	ToDate             time.Time         `query:"to_date" validate:"omitempty,gtefield=FromDate"`
	Page               int               `query:"page" validate:"min=1"`
	PageSize           int               `query:"page_size" validate:"min=1,max=100"`
}

// UpdateAppointmentDTO represents the request body for updating an appointment
type UpdateAppointmentDTO struct {
	AppointmentID string            `param:"appointment_id" validate:"required,uuid"`
	Status        AppointmentStatus `json:"status" validate:"required,oneof=pending confirmed in_progress completed cancelled rescheduled"`
	Notes         string            `json:"notes" validate:"max=500"`
}

// CancelAppointmentDTO represents the request body for cancelling an appointment
type CancelAppointmentDTO struct {
	AppointmentID string `param:"appointment_id" validate:"required,uuid"`
	Reason        string `json:"reason" validate:"required,max=500"`
}

// RescheduleAppointmentDTO represents the request body for rescheduling an appointment
type RescheduleAppointmentDTO struct {
	AppointmentID    string    `param:"appointment_id" validate:"required,uuid"`
	NewScheduleID    string    `json:"new_schedule_id" validate:"required,uuid"`
	NewDate          time.Time `json:"new_date" validate:"required"`
	NewTimeSlot      string    `json:"new_time_slot" validate:"required"`
	RescheduleReason string    `json:"reschedule_reason" validate:"required,max=500"`
}

// CreateScheduleDTO represents the request body for creating a schedule
