package ex

import (
	"fmt"
	"net/smtp"
)

type NotificationAdapter struct{
	config *GmailConfig
}

func (n *NotificationAdapter) SendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", n.config.Host, n.config.Port),
		auth,
		n.config.From,
		[]string{to},
		msg,
	)
}

func (n *NotificationAdapter) SendSMS(phone string, message string) error {
	// TODO: Add your actual SMS sending implementation
	// For now, just log
	fmt.Printf("Sending SMS to %s\nMessage: %s\n", phone, message)
	return nil
}

func NewNotificationAdapter(con *GmailConfig) *NotificationAdapter {
	return &NotificationAdapter{config: con}
}
