package ports

import (
	"context"

	"github.com/diagnoxix/adapters/db"
)
// NotificationService is the interface that wraps notification operations
type NotificationService interface {
	// SendEmail sends an email to a recipient
	SendEmail(
		to,
		subject,
		templateName string,
		data interface{},
	) error
	// SendSMS sends an SMS message to a phone number
	SendSMS(
		phone,
		message string,
	) error
}

type NotificationRepository interface {
	CreateNotification(
		ctx context.Context,
		arg db.CreateNotificationParams,
	) (*db.Notification, error)

	GetNotifications(
		ctx context.Context,
		arg db.GetUserNotificationsParams,
	) ([]*db.Notification, error)

	MarkNotificationRead(
		ctx context.Context,
		arg db.MarkAsReadParams,
	) (*db.Notification, error)

	MarkNotificationReadAll(
		ctx context.Context,
		arg string,
	) error
}
