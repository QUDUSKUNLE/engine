package emails

import (
	"bytes"
	"html/template"
)

const staffNotificationTemplate = `
{{define "staff_notification"}}
<p><strong>Hello {{.StaffName}},</strong></p>
<p>A new appointment has been scheduled that requires your attention:</p>

<div class="details">
    <ul>
        <li><strong>Patient:</strong> {{.PatientName}}</li>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
    </ul>
</div>

<div class="note">
    <p><strong>Actions Required:</strong></p>
    <ul>
        <li>Review patient history if available</li>
        <li>Prepare necessary equipment/materials</li>
        <li>Update your schedule accordingly</li>
        {{if .SpecialInstructions}}
        <li>{{.SpecialInstructions}}</li>
        {{end}}
    </ul>
</div>

<p>Please log in to the staff portal to view complete appointment details.</p>
{{end}}
`

// GetStaffNotificationTemplate returns the rendered staff notification email
func GetStaffNotificationTemplate(data StaffNotificationData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(staffNotificationTemplate))

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
