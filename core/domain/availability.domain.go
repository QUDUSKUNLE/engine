package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/medicue/adapters/db"
)

type AvailabilitySlot struct {
	ID                 uuid.UUID  `json:"id"`
	DiagnosticCentreID uuid.UUID  `json:"diagnostic_centre_id"`
	DayOfWeek          db.Weekday `json:"day_of_week"`
	StartTime          time.Time  `json:"start_time"`
	EndTime            time.Time  `json:"end_time"`
	MaxAppointments    int        `json:"max_appointments"`
	SlotDuration       string     `json:"slot_duration"` // e.g., "30 minutes"
	BreakTime          string     `json:"break_time"`    // e.g., "5 minutes"
}

type Slots struct {
	DayOfWeek       db.Weekday `json:"day_of_week" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime       string     `json:"start_time" validate:"required,time"`                 // Format: "HH:MM"
	EndTime         string     `json:"end_time" validate:"required,time,gtfield=StartTime"` // Format: "HH:MM"
	MaxAppointments int        `json:"max_appointments" validate:"required,min=0"`
	SlotDuration    string     `json:"slot_duration" validate:"required"` // e.g., "30 minutes"
	BreakTime       string     `json:"break_time" validate:"required"`    // e.g., "5 minutes"
}

type CreateAvailabilityDTO struct {
	DiagnosticCentreID string  `json:"diagnostic_centre_id" validate:"required,uuid"`
	Slots              []Slots `json:"slots" validate:"required,min=1,dive"`
}

type UpdateAvailabilityDTO struct {
	StartTime       *string `json:"start_time,omitempty" validate:"omitempty,time"`
	EndTime         *string `json:"end_time,omitempty" validate:"omitempty,time"`
	MaxAppointments *int    `json:"max_appointments,omitempty" validate:"omitempty,min=0"`
	SlotDuration    *string `json:"slot_duration,omitempty"`
	BreakTime       *string `json:"break_time,omitempty"`
}

type GetAvailabilityDTO struct {
	DiagnosticCentreID string `json:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          string `json:"day_of_week,omitempty" validate:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
}

// UpdateManyAvailabilitySlot represents a single slot in the update many request
type UpdateManyAvailabilitySlot struct {
	DiagnosticCentreID string     `json:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          db.Weekday `json:"day_of_week" validate:"required"`
	StartTime          *string    `json:"start_time,omitempty" validate:"omitempty,datetime=15:04:05"`
	EndTime            *string    `json:"end_time,omitempty" validate:"omitempty,datetime=15:04:05"`
	MaxAppointments    *int       `json:"max_appointments,omitempty" validate:"omitempty,min=1"`
	SlotDuration       *string    `json:"slot_duration,omitempty" validate:"omitempty"`
	BreakTime          *string    `json:"break_time,omitempty" validate:"omitempty"`
}

// UpdateManyAvailabilityDTO represents the request body for updating multiple availability slots
type UpdateManyAvailabilityDTO struct {
	Slots []UpdateManyAvailabilitySlot `json:"slots" validate:"required,min=1,dive"`
}
