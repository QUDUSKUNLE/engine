package emails

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
