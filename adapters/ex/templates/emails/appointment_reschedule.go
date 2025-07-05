package emails

import (
	"bytes"
	"html/template"
)

const appointmentRescheduleTemplate = `
{{define "appointment_reschedule"}}
<p><strong>Dear {{.PatientName}},</strong></p>
<p>Your appointment has been successfully rescheduled. Here are your new appointment details:</p>

<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>New Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>New Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
    </ul>
</div>

<div class="note">
    <p>Please let us know if this new schedule does not work for you. We'll be happy to find another suitable time.</p>
</div>
{{end}}
`

// GetAppointmentRescheduleTemplate returns the rendered appointment reschedule email
func GetAppointmentRescheduleTemplate(data AppointmentEmailData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(appointmentRescheduleTemplate))

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
