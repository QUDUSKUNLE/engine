package emails

const appointmentRescheduleTemplate = `
{{define "appointment_reschedule"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>Your appointment has been successfully rescheduled. Here are the new details:</p>

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
    <p>If this new schedule doesn’t work for you, feel free to contact us to find a more suitable time.</p>
</div>

<p>Thank you for using Diagnoxix.</p>
{{end}}
`
