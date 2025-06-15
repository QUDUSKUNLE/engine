package ports

// EmailService is the interface that wraps the basic email operations
type NotificationService interface {
	SendEmail(to string, subject string, body string) error
}
