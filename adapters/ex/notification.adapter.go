package ex

import (
	"github.com/medivue/adapters/ex/templates/emails"
)

type NotificationAdapter struct {
	config        *GmailConfig
	emailTemplate emails.EmailTemplateHandler
}

func NewNotificationAdapter(con *GmailConfig) *NotificationAdapter {
	return &NotificationAdapter{
		config:        con,
		emailTemplate: *emails.NewEmailTemplateHandler(),
	}
}
