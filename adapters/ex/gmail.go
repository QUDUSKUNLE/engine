package ex

import (
	"fmt"
	"net/smtp"
)

type GmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string // Gmail App Password
	From     string
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

func (n *NotificationAdapter) SendEmail(
	to string,
	subject string,
	templateName string,
	data interface{},
) error {
	auth := smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.Host)

	// Generate HTML content from template
	htmlContent, err := n.emailTemplate.ExecuteTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"X-Priority: 1 (Highest)\r\n"+
		"X-MSMail-Priority: High\r\n"+
		"Importance: High\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, htmlContent))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", n.config.Host, n.config.Port),
		auth,
		n.config.From,
		[]string{to},
		msg,
	)
}

func (n *NotificationAdapter) SendSMS(phone string, message string) error {
	// TODO: Add your actual SMS sending implementation
	// For now, just log
	fmt.Printf("Sending SMS to %s\nMessage: %s\n", phone, message)
	return nil
}
