package ports

// EmailService is the interface that wraps the basic email operations
type EmailService interface {
	Send(to string, subject string, body string) error
}
