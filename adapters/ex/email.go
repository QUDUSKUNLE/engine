package ex

type EmailService interface {
	SendEmail(to string, subject string, body string) error
	SendSMS(phone string, message string) error
}

type EmailAdapter struct{}

func (e *EmailAdapter) SendEmail(to string, subject string, body string) error {
	// Add your email sending implementation here
	return nil
}

func (n *EmailAdapter) SendSMS(phone string, message string) error {
	return nil
}
