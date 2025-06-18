package ex

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"fmt"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridEmailService implements EmailService using SendGrid
type SendGridEmailService struct {
	client *sendgrid.Client
	from   string
}

type GmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string // Gmail App Password
	From     string
}


// NewSendGridEmailService creates a new SendGrid email service
func NewSendGridEmailService() *SendGridEmailService {
	return &SendGridEmailService{
		client: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
		from:   os.Getenv("SENDGRID_FROM_ADDRESS"),
	}
}

func NewGmailConfig(c GmailConfig) *GmailConfig {
	return &GmailConfig{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
		From:     c.From,
	}
}

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