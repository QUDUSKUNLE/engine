package ex

import (
	"github.com/diagnoxix/adapters/external/templates/emails"
)

type NotificationAdapter struct {
	config        *EmailConfig
	emailTemplate emails.EmailTemplateHandler
}

func NewNotificationAdapter(con *EmailConfig) *NotificationAdapter {
	return &NotificationAdapter{
		config:        con,
		emailTemplate: *emails.NewEmailTemplateHandler(),
	}
}
