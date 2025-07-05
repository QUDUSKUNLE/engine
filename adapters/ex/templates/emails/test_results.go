package emails

import (
	"bytes"
	"html/template"
)

const testResultsTemplate = `
{{define "test_results_available"}}
<p><strong>Dear {{.PatientName}},</strong></p>
<p>Your test results are now available for review. Here are the details of your test:</p>

<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
    </ul>
</div>

<div class="note">
    <p><strong>To view your results securely:</strong></p>
    <ol>
        <li>Log in to your patient portal at: <a href="{{.ResultsPortalURL}}">{{.ResultsPortalURL}}</a></li>
        <li>Navigate to "Test Results" section</li>
        <li>Select the appointment date mentioned above</li>
    </ol>
</div>

<p>If you have any questions about your results, please contact your healthcare provider.</p>
{{end}}
`

// GetTestResultsTemplate returns the rendered test results available email
func GetTestResultsTemplate(data interface{}) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(testResultsTemplate))

	var buf bytes.Buffer
	err := contentTemplate.ExecuteTemplate(&buf, "base", map[string]interface{}{
		"Title":         "Test Results Available",
		"Icon":          "🔬",
		"Content":       data,
		"FooterContent": "For any technical issues accessing your results, please contact our support team.",
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
