package emails

const paymentConfirmationTemplate = `
{{define "payment_confirmation"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>Thank you for your payment. Below are your transaction details:</p>

<div class="details">
    <ul>
        <li><strong>Transaction ID:</strong> {{.TransactionID}}</li>
        <li><strong>Amount Paid:</strong> {{.PaymentAmount}}</li>
        <li><strong>Payment Method:</strong> {{.PaymentMethod}} {{if .LastFourDigits}}ending in {{.LastFourDigits}}{{end}}</li>
        <li><strong>Date:</strong> {{.PaymentDate | formatDate}}</li>
    </ul>
</div>

<div class="note">
    <p><strong>Appointment Details:</strong></p>
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Test:</strong> {{.TestType}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
    </ul>
</div>

<p>This email serves as your receipt. Please keep it for your records.</p>
<p>If you have any questions or concerns, feel free to contact us.</p>
{{end}}
`
