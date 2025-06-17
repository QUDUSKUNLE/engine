package ex

import (
	"fmt"
)

type NotificationAdapter struct{}

func (n *NotificationAdapter) SendEmail(to string, subject string, body string) error {
	// TODO: Add your actual email sending implementation
	// For now, just log
	fmt.Printf("Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}

func (n *NotificationAdapter) SendSMS(phone string, message string) error {
	// TODO: Add your actual SMS sending implementation
	// For now, just log
	fmt.Printf("Sending SMS to %s\nMessage: %s\n", phone, message)
	return nil
}
