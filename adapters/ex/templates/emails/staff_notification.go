package emails

const staffNotificationTemplate = `
{{define "staff_notification"}}
<p><strong>Dear {{.StaffName}},</strong></p>

<p>A new appointment has been scheduled that requires your attention. Below are the appointment details:</p>

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
        <li>Prepare necessary equipment and materials</li>
        <li>Update your schedule accordingly</li>
        {{if .SpecialInstructions}}
        <li>{{.SpecialInstructions}}</li>
        {{end}}
    </ul>
</div>

<p>Please log in to the staff portal to view complete appointment details.</p>
<p>Thank you for your continued service.</p>
{{end}}
`
