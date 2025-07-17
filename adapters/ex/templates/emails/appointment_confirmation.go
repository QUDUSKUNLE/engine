package emails

const appointmentConfirmationTemplate = `
{{define "appointment_confirmation"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>Your appointment has been successfully confirmed. Here are the appointment details:</p>

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

<p>You can view or manage your appointment by logging into your Diagnoxix account at any time.</p>
<p>Thank you for choosing Diagnoxix!</p>
{{end}}
`
