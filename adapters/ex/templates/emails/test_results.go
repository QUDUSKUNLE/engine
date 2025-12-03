package emails

const testResultsTemplate = `
{{define "test_results_available"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>Your test results are now available for review. Here are the details:</p>

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
        <li>Navigate to the "Test Results" section</li>
        <li>Select the appointment date mentioned above</li>
    </ol>
</div>

<p>If you have any questions about your results, please contact your healthcare provider.</p>
<p>Thank you for using DiagnoxixAI.</p>
{{end}}
`
