package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AvailabilitySlot struct {
	ID                 string `json:"id"`
	DiagnosticCentreID string `json:"diagnostic_centre_id"`
	DayOfWeek          string    `json:"day_of_week"`
	StartTime          string    `json:"start_time"`
	EndTime            string    `json:"end_time"`
	MaxAppointments    int32       `json:"max_appointments"`
	SlotDuration       string     `json:"slot_duration"` // minutes
	BreakTime          string     `json:"break_time"`    // minutes
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

type Slots struct {
	DayOfWeek       string    `json:"day_of_week" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime       time.Time `json:"start_time" validate:"required"`
	EndTime         time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	MaxAppointments int32     `json:"max_appointments" validate:"required,min=0"`
	SlotDuration    int32     `json:"slot_duration" validate:"required,min=1"` // minutes
	BreakTime       int32     `json:"break_time" validate:"required,min=0"`    // minutes
}

// UnmarshalJSON implements custom JSON unmarshaling for Slots
func (s *Slots) UnmarshalJSON(data []byte) error {
	aux := &struct {
		DayOfWeek       string `json:"day_of_week"`
		StartTime       string `json:"start_time"`
		EndTime         string `json:"end_time"`
		MaxAppointments int32  `json:"max_appointments"`
		SlotDuration    int32  `json:"slot_duration"`
		BreakTime       int32  `json:"break_time"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.DayOfWeek = aux.DayOfWeek
	s.MaxAppointments = aux.MaxAppointments
	s.SlotDuration = aux.SlotDuration
	s.BreakTime = aux.BreakTime

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

type CreateSingleAvailabilityDTO struct {
	DiagnosticCentreID uuid.UUID `json:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          string    `json:"day_of_week" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime          time.Time `json:"start_time" validate:"required"`
	EndTime            time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	MaxAppointments    int32     `json:"max_appointments" validate:"required,min=1"`
	SlotDuration       int32     `json:"slot_duration" validate:"required"` // e.g., "30 minutes"
	BreakTime          int32     `json:"break_time" validate:"required"`    // e.g., "5 minutes"
}

// MarshalJSON implements custom JSON marshaling for Slots
func (s Slots) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		DayOfWeek       string `json:"day_of_week"`
		StartTime       string `json:"start_time"`
		EndTime         string `json:"end_time"`
		MaxAppointments int32  `json:"max_appointments"`
		SlotDuration    int32  `json:"slot_duration"`
		BreakTime       int32  `json:"break_time"`
	}{
		DayOfWeek:       s.DayOfWeek,
		StartTime:       s.StartTime.Format(time.RFC3339),
		EndTime:         s.EndTime.Format(time.RFC3339),
		MaxAppointments: s.MaxAppointments,
		SlotDuration:    s.SlotDuration,
		BreakTime:       s.BreakTime,
	})
}

type CreateAvailabilityDTO struct {
	DiagnosticCentreID string  `json:"diagnostic_centre_id" validate:"required,uuid"`
	Slots              []Slots `json:"slots" validate:"required,min=1,dive"`
}

// UnmarshalJSON implements custom JSON unmarshaling for CreateAvailabilityDTO
func (dto *CreateAvailabilityDTO) UnmarshalJSON(data []byte) error {
	type Alias CreateAvailabilityDTO
	aux := &struct {
		*Alias
		Slots []json.RawMessage `json:"slots"`
	}{
		Alias: (*Alias)(dto),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	dto.Slots = make([]Slots, len(aux.Slots))
	for i, slotData := range aux.Slots {
		var slot Slots
		if err := json.Unmarshal(slotData, &slot); err != nil {
			return fmt.Errorf("invalid slot data at index %d: %w", i, err)
		}
		dto.Slots[i] = slot
	}

	return nil
}

type UpdateAvailabilityDTO struct {
	StartTime       *time.Time `json:"start_time,omitempty" validate:"omitempty,time"`
	EndTime         *time.Time `json:"end_time,omitempty" validate:"omitempty,time"`
	MaxAppointments *int32     `json:"max_appointments,omitempty" validate:"omitempty,min=0"`
	SlotDuration    *int32     `json:"slot_duration,omitempty" validate:"omitempty,min=1"` // minutes
	BreakTime       *int32     `json:"break_time,omitempty" validate:"omitempty,min=0"`    // minutes
}

type GetAvailabilityDTO struct {
	DiagnosticCentreID string `param:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          string `query:"day_of_week,omitempty" validate:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
}

// UpdateManyAvailabilitySlot represents a single slot in the update many request
type UpdateManyAvailabilitySlot struct {
	DiagnosticCentreID string     `json:"diagnostic_centre_id" validate:"required,uuid"`
	DayOfWeek          string     `json:"day_of_week" validate:"required"`
	StartTime          *time.Time `json:"start_time,omitempty" validate:"omitempty,time"`
	EndTime            *time.Time `json:"end_time,omitempty" validate:"omitempty,time"`
	MaxAppointments    *int32     `json:"max_appointments,omitempty" validate:"omitempty,min=1"`
	SlotDuration       *int32     `json:"slot_duration,omitempty" validate:"omitempty,min=1"` // minutes
	BreakTime          *int32     `json:"break_time,omitempty" validate:"omitempty,min=0"`    // minutes
}

// UpdateManyAvailabilityDTO represents the request body for updating multiple availability slots
type UpdateManyAvailabilityDTO struct {
	Slots []UpdateManyAvailabilitySlot `json:"slots" validate:"required,min=1,dive"`
}
