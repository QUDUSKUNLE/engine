package emails

import (
	"bytes"
	"html/template"
)

const appointmentConfirmationTemplate = `
{{define "appointment_confirmation"}}
<p><strong>Dear {{.PatientName}},</strong></p>
<p>Your appointment has been successfully confirmed. Here are your appointment details:</p>

<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
    </ul>
</div>

{{if .Notes}}
<div class="note">
    <p><strong>Important Notes:</strong></p>
    <p>{{.Notes}}</p>
</div>
{{end}}

<p>You can view or manage your appointment by logging into your account at any time.</p>
{{end}}
`

// GetAppointmentConfirmationTemplate returns the rendered appointment confirmation email
func GetAppointmentConfirmationTemplate(data AppointmentEmailData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(appointmentConfirmationTemplate))

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
