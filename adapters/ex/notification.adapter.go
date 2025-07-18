package ex

import (
	"github.com/medivue/adapters/ex/templates/emails"
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
