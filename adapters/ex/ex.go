package ex

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
)

// SendGridEmailService implements EmailService using SendGrid
type SendGridEmailService struct {
	client *sendgrid.Client
	from   string
}

// NewSendGridEmailService creates a new SendGrid email service
func NewSendGridEmailService() *SendGridEmailService {
	return &SendGridEmailService{
		client: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
		from:   os.Getenv("EMAIL_FROM_ADDRESS"),
	}
}
