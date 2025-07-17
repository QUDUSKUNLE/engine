package ex

type EmailType string

const (
	GMAIL EmailType = "GMAIL"
	ZOHO  EmailType = "ZOHO"
)

type EmailConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string // App Password
	From      string
	EmailType EmailType
}

func NewEmailConfig(c EmailConfig) *EmailConfig {
	return &EmailConfig{
		Host:      c.Host,
		Port:      c.Port,
		Username:  c.Username,
		Password:  c.Password,
		From:      c.From,
		EmailType: c.EmailType,
	}
}
