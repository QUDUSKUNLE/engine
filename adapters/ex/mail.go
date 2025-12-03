package ex

import (
	"fmt"
	"net/smtp"
)

func (n *NotificationAdapter) auth() smtp.Auth {
	return smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.Host)
}

func (n *NotificationAdapter) SendEmail(
	to string,
	subject string,
	templateName string,
	data interface{},
) error {
	auth := n.auth()
	// Generate HTML content from template
	htmlContent, err := n.emailTemplate.ExecuteTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	var msg []byte
	switch n.config.EmailType {
	case GMAIL, ZOHO:
		msg = []byte(fmt.Sprintf(
			"From: DiagnoxixAI <%s>\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/html; charset=utf-8\r\n"+
				"X-Priority: 1 (Highest)\r\n"+
				"X-MSMail-Priority: High\r\n"+
				"Importance: High\r\n"+
				"\r\n"+
				"%s\r\n", n.config.From, to, subject, htmlContent),
		)
	}

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
