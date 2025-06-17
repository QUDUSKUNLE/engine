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
