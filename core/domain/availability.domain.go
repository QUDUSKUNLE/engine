package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/medicue/adapters/db"
)

type AvailabilitySlot struct {
	ID                 uuid.UUID `json:"id"`
	DiagnosticCentreID uuid.UUID `json:"diagnostic_centre_id"`
	DayOfWeek          string    `json:"day_of_week"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	MaxAppointments    int       `json:"max_appointments"`
	SlotDuration       string    `json:"slot_duration"` // e.g., "30 minutes"
	BreakTime          string    `json:"break_time"`    // e.g., "5 minutes"
}

type Slots struct {
	DayOfWeek       string    `json:"day_of_week" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime       time.Time `json:"start_time" validate:"required"`
	EndTime         time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	MaxAppointments int32       `json:"max_appointments" validate:"required,min=0"`
	SlotDuration    int32    `json:"slot_duration" validate:"required"` // e.g., "30 minutes"
	BreakTime       int32    `json:"break_time" validate:"required"`    // e.g., "5 minutes"
}

type CreateSingleAvailabilityDTO struct {
	DiagnosticCentreID uuid.UUID `json:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          string    `json:"day_of_week" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime          time.Time `json:"start_time" validate:"required"`
	EndTime            time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	MaxAppointments    int32     `json:"max_appointments" validate:"required,min=1"`
	SlotDuration       int32   `json:"slot_duration" validate:"required"` // e.g., "30 minutes"
	BreakTime          int32   `json:"break_time" validate:"required"`    // e.g., "5 minutes"
}

// MarshalJSON implements custom JSON marshaling for Slots
func (s Slots) MarshalJSON() ([]byte, error) {
	type Alias Slots
	return json.Marshal(&struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Alias
	}{
		StartTime: s.StartTime.Format(time.RFC3339),
		EndTime:   s.EndTime.Format(time.RFC3339),
		Alias:     Alias(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Slots
func (s *Slots) UnmarshalJSON(data []byte) error {
	type Alias Slots
	aux := &struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	s.StartTime, err = time.Parse(time.RFC3339, aux.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start time format: %w", err)
	}

	s.EndTime, err = time.Parse(time.RFC3339, aux.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end time format: %w", err)
	}

	return nil
}

type CreateAvailabilityDTO struct {
	DiagnosticCentreID string  `json:"diagnostic_centre_id" validate:"required,uuid"`
	Slots              []Slots `json:"slots" validate:"required,min=1,dive"`
}

type UpdateAvailabilityDTO struct {
	StartTime       *time.Time `json:"start_time,omitempty" validate:"omitempty,time"`
	EndTime         *time.Time `json:"end_time,omitempty" validate:"omitempty,time"`
	MaxAppointments *int32     `json:"max_appointments,omitempty" validate:"omitempty,min=0"`
	SlotDuration    *int32     `json:"slot_duration,omitempty"`
	BreakTime       *int32     `json:"break_time,omitempty"`
}

type GetAvailabilityDTO struct {
	DiagnosticCentreID string `param:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          string `query:"day_of_week,omitempty" validate:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
}

// UpdateManyAvailabilitySlot represents a single slot in the update many request
type UpdateManyAvailabilitySlot struct {
	DiagnosticCentreID string     `json:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          db.Weekday `json:"day_of_week" validate:"required"`
	StartTime          *time.Time `json:"start_time,omitempty" validate:"omitempty,time"`
	EndTime            *time.Time `json:"end_time,omitempty" validate:"omitempty,time"`
	MaxAppointments    *int32     `json:"max_appointments,omitempty" validate:"omitempty,min=1"`
	SlotDuration       *int32     `json:"slot_duration,omitempty" validate:"omitempty"`
	BreakTime          *int32     `json:"break_time,omitempty" validate:"omitempty"`
}

// UpdateManyAvailabilityDTO represents the request body for updating multiple availability slots
type UpdateManyAvailabilityDTO struct {
	Slots []UpdateManyAvailabilitySlot `json:"slots" validate:"required,min=1,dive"`
}
