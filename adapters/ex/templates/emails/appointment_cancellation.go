package emails

import (
	"bytes"
	"html/template"
)

const appointmentCancellationTemplate = `
{{define "appointment_cancellation"}}
<p><strong>Dear {{.PatientName}},</strong></p>
<p>Your appointment has been cancelled. Here are the details of the cancelled appointment:</p>

<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
    </ul>
</div>

<p>You can book a new appointment by logging into your account at any time.</p>
{{end}}
`

// GetAppointmentCancellationTemplate returns the rendered appointment cancellation email
func GetAppointmentCancellationTemplate(data AppointmentEmailData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(appointmentCancellationTemplate))

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
