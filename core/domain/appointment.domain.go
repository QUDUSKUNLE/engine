package domain

import "time"

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

// CreateAppointmentDTO represents the request body for creating an appointment
type CreateAppointmentDTO struct {
	DiagnosticCentreID string    `json:"diagnostic_centre_id" validate:"required,uuid"`
	PatientID          string    `json:"patient_id" validate:"required,uuid"`
	ScheduleID         string    `json:"schedule_id" validate:"required,uuid"`
	AppointmentDate    time.Time `json:"appointment_date" validate:"required"`
	TimeSlot           string    `json:"time_slot" validate:"required"`
	Notes              string    `json:"notes" validate:"max=500"`
}

// GetAppointmentDTO represents the request parameters for getting an appointment
type GetAppointmentDTO struct {
	AppointmentID string `param:"appointment_id" validate:"required,uuid"`
}

// ListAppointmentsDTO represents the query parameters for listing appointments
type ListAppointmentsDTO struct {
	DiagnosticCentreID string           `query:"diagnostic_centre_id" validate:"omitempty,uuid"`
	PatientID          string           `query:"patient_id" validate:"omitempty,uuid"`
	Status             AppointmentStatus `query:"status" validate:"omitempty,oneof=pending confirmed in_progress completed cancelled rescheduled"`
	FromDate           time.Time        `query:"from_date" validate:"omitempty"`
	ToDate             time.Time        `query:"to_date" validate:"omitempty,gtefield=FromDate"`
	Page               int              `query:"page" validate:"min=1"`
	PageSize           int              `query:"page_size" validate:"min=1,max=100"`
}

// UpdateAppointmentDTO represents the request body for updating an appointment
type UpdateAppointmentDTO struct {
	AppointmentID string           `param:"appointment_id" validate:"required,uuid"`
	Status        AppointmentStatus `json:"status" validate:"required,oneof=pending confirmed in_progress completed cancelled rescheduled"`
	Notes         string           `json:"notes" validate:"max=500"`
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
