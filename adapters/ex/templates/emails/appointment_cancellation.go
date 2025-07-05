package emails

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
