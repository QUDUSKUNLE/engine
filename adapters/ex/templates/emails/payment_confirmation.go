package emails

import (
	"bytes"
	"html/template"
)

const paymentConfirmationTemplate = `
{{define "payment_confirmation"}}
<p><strong>Dear {{.PatientName}},</strong></p>
<p>Thank you for your payment. Here are your transaction details:</p>

<div class="details">
    <ul>
        <li><strong>Transaction ID:</strong> {{.TransactionID}}</li>
        <li><strong>Amount Paid:</strong> {{.PaymentAmount | formatCurrency}}</li>
        <li><strong>Payment Method:</strong> {{.PaymentMethod}} {{if .LastFourDigits}}ending in {{.LastFourDigits}}{{end}}</li>
        <li><strong>Date:</strong> {{.PaymentDate | formatDate}}</li>
    </ul>
</div>

<div class="note">
    <p><strong>Appointment Details:</strong></p>
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
    </ul>
</div>

<p>This email serves as your receipt. Please keep it for your records.</p>
{{end}}
`

// GetPaymentConfirmationTemplate returns the rendered payment confirmation email
func GetPaymentConfirmationTemplate(data PaymentData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(paymentConfirmationTemplate))

	var buf bytes.Buffer
	err := contentTemplate.ExecuteTemplate(&buf, "base", map[string]interface{}{
		"Title":         data.Title,
		"Icon":          data.Icon,
		"Content":       data.Content,
		"FooterContent": data.FooterContent,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
