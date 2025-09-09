package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	AppointmentCreated     NotificationType = "APPOINTMENT_CREATED"
	AppointmentConfirmed   NotificationType = "APPOINTMENT_CONFIRMED"
	AppointmentRescheduled NotificationType = "APPOINTMENT_RESCHEDULED"
	AppointmentCancelled   NotificationType = "APPOINTMENT_CANCELLED"
	AppointmentReminder    NotificationType = "APPOINTMENT_REMINDER"
)

type (
	Notification struct {
		ID        uuid.UUID              `json:"id"`
		UserID    uuid.UUID              `json:"user_id"`
		Type      NotificationType       `json:"type"`
		Title     string                 `json:"title"`
		Message   string                 `json:"message"`
		Read      bool                   `json:"read"`
		ReadAt    *time.Time             `json:"read_at,omitempty"`
		CreatedAt time.Time              `json:"created_at"`
		Metadata  map[string]interface{} `json:"metadata,omitempty"`
	}
	MarkNotificationReadDTO struct {
		NotificationID uuid.UUID `param:"notification_id" json:"notification_id" validate:"uuid,required"`
	}
)
