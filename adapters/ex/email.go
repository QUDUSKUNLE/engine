package ex

type EmailService interface {
	// Send sends an email to the specified recipient
	Send(to string, subject string, body string) error
}

type EmailAdapter struct{}

func (e *EmailAdapter) Send(to string, subject string, body string) error {
	// Add your email sending implementation here
	return nil
}
