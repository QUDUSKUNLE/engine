package emails

import (
	"bytes"
	"html/template"
)

const appointmentReminderTemplate = `
{{define "appointment_reminder"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>This is a friendly reminder about your upcoming appointment. Here are the details:</p>

<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
    </ul>
</div>

<div class="note">
    <p><strong>Important Reminders:</strong></p>
    <ul>
        <li>Please arrive 15 minutes before your scheduled time.</li>
        <li>Bring any relevant medical records or referrals.</li>
        <li>Don't forget to bring a valid ID.</li>
        {{if .PreTestInstructions}}
        <li>{{.PreTestInstructions}}</li>
        {{end}}
    </ul>
</div>

<p>If you need to reschedule or cancel your appointment, please do so at least 24 hours in advance.</p>
<p>We look forward to seeing you soon!</p>
{{end}}
`

// GetAppointmentReminderTemplate returns the rendered appointment reminder email
func GetAppointmentReminderTemplate(data AppointmentEmailData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(appointmentReminderTemplate))

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
