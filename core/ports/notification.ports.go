package ports

// NotificationService is the interface that wraps notification operations
type NotificationService interface {
	// SendEmail sends an email to a recipient
	SendEmail(to string, subject string, templateName string, data interface{}) error
	// SendSMS sends an SMS message to a phone number
	SendSMS(phone string, message string) error
}
