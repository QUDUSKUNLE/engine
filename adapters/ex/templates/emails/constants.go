package emails

// Email template titles
const (
	TitleAppointmentConfirmed       = "Appointment Confirmed"
	TitleAppointmentCancelled       = "Appointment Cancelled"
	TitleAppointmentReminder        = "Appointment Reminder"
	TitleAppointmentReschedule      = "Appointment Rescheduled"
	TitlePaymentConfirmed           = "Payment Confirmed"
	TitleTestResults                = "Test Results Available"
	TitleStaffNotification          = "New Appointment"
	TitlePolicyUpdate               = "Policy Update"
	TitleResetPassword              = "Reset Your Password"
	TitleDiagnosticCentreManager    = "Diagnostic Centre Manager Notification"
	TitleDiagnosticCentreManagement = "Diagnostic Centre Management Notification"
	TitleEmailVerification          = "Email Verification"
)

const (
	TemplateTestResults         = "test_results"
	TemplatePolicyUpdate        = "policy_update"
	TemplateStaffNotification   = "staff_notification"
	TemplateResetPassword       = "password_reset"
	TemplateEmailVerification   = "email_verification"
	TemplatePaymentConfirmation = "payment_confirmation"

	TemplateDiagnosticCentreManagement = "dc_management_notification"
	TemplateDiagnosticCentreManager    = "dc_manager_notification"

	TemplateAppointmentConfirmed  = "appointment_confirmation"
	TemplateAppointmentCancelled  = "appointment_cancellation"
	TemplateAppointmentReminder   = "appointment_reminder"
	TemplateAppointmentReschedule = "appointment_reschedule"
)

const (
	SubjectEmailVerification          = "Sign up for DiagnoxixAI - Email Verification"
	SubjectResetPassword              = "Reset Your Password - DiagnoxixAI"
	SubjectDiagnosticCentreManager    = "Diagnostic Centre Manager Notification"
	SubjectDiagnosticCentreManagement = "Diagnostic Centre Management Notification"
)

// Email template icons (emojis)
const (
	IconConfirmed         = "✅"
	IconCancelled         = "❌"
	IconReminder          = "🔔"
	IconReschedule        = "🔄"
	IconPayment           = "💳"
	IconTestResults       = "🔬"
	IconStaff             = "👥"
	IconPolicy            = "📋"
	IconEmailVerification = "✉️"
	IconResetPassword     = "🔄"
)

// Footer messages
const (
	FooterChanges           = "If you need to make any changes to your appointment, please contact us as soon as possible."
	FooterSupport           = "For any questions or concerns, please contact our support team."
	FooterPayment           = "If you have any questions about your payment, please contact our billing department."
	FooterResults           = "For any technical issues accessing your results, please contact our support team."
	FooterStaff             = "If you have any conflicts or concerns, please notify the scheduling department immediately."
	FooterPolicy            = "If you have any questions about these updates, please don't hesitate to contact us."
	FooterEmailVerification = "If you need any assistance, please contact our support team."
	FooterResetPassord      = "If you need any assistance, please contact our support team."
)

// Common template messages
const (
	MsgArriveEarly  = "Please arrive 15 minutes before your scheduled time"
	MsgBringRecords = "Bring any relevant medical records or referrals"
	MsgBringID      = "Don't forget to bring a valid ID"
	MsgCancelPolicy = "If you need to cancel, please do so at least 24 hours in advance"
	MsgViewPortal   = "You can view or manage your appointment by logging into your account"
	MsgKeepReceipt  = ""
)

// Template partial content
const (
	PartialAppointmentDetails = `
<div class="details">
    <ul>
        <li><strong>Appointment ID:</strong> {{.AppointmentID}}</li>
        <li><strong>Date:</strong> {{.AppointmentDate | formatDate}}</li>
        <li><strong>Time:</strong> {{.TimeSlot}}</li>
        <li><strong>Centre:</strong> {{.CentreName}}</li>
        <li><strong>Test:</strong> {{.TestType | formatTest}}</li>
    </ul>
</div>`

	PartialPaymentDetails = `
<div class="details">
    <ul>
        <li><strong>Transaction ID:</strong> {{.TransactionID}}</li>
        <li><strong>Amount Paid:</strong> {{.PaymentAmount | formatNaira}}</li>
        <li><strong>Payment Method:</strong> {{.PaymentMethod}}</li>
        <li><strong>Date:</strong> {{.PaymentDate | formatDate}}</li>
    </ul>
</div>`
)
