package emails

const appointmentCancellationTemplate = `
{{define "appointment_cancellation"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>Your appointment has been cancelled. Here are the details:</p>

<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Diagnostic Centre:</strong> {{.CentreName}}</li>
    </ul>
</div>

<p>You can book a new appointment by logging into your Medivue account anytime.</p>
<p>If you have any questions, feel free to contact support.</p>
{{end}}
`
