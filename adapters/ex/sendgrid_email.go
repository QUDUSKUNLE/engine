package ex

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends an email using SendGrid
func (s *SendGridEmailService) SendEmail(to string, subject string, body string) error {
	from := mail.NewEmail("Medicue", s.from)
	recipient := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, recipient, "", body)

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("email service returned error status: %d", response.StatusCode)
	}

	return nil
}
