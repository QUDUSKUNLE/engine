package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// AppointmentEmailData contains common data for appointment efunc GetAppointmentC// GetAppointmentConfirmationTemplate returns the rendered appointment confirmation email template

const emailBaseLayout = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen-Sans, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: #4A90E2;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
            margin-bottom: 20px;
        }
        .content {
            background: #fff;
            padding: 20px;
            border-radius: 0 0 8px 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 14px;
            margin-top: 20px;
            padding: 20px;
        }
        ul, ol {
            margin: 15px 0;
            padding-left: 30px;
        }
        li {
            margin: 10px 0;
        }
        .highlight {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-radius: 4px;
            border-left: 4px solid #4A90E2;
        }
        .warning {
            background: #fff3cd;
            padding: 15px;
            margin: 10px 0;
            border-radius: 4px;
            border-left: 4px solid #ffc107;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            background: #4A90E2;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            margin: 10px 0;
        }
        @media screen and (max-width: 480px) {
            body {
                padding: 10px;
            }
            .header, .content, .footer {
                padding: 15px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>{{.Title}}</h2>
    </div>
    <div class="content">
        {{.Content}}
    </div>
    <div class="footer">
        <p>Best regards,<br/>{{.AppName}} Team</p>
        {{if .FooterContent}}
        <p>{{.FooterContent}}</p>
        {{end}}
    </div>
</body>
</html>
`

const appointmentConfirmationTmpl = `
{{define "appointment_confirmation"}}
<h2>Appointment Confirmation</h2>
<p>Dear {{.PatientName}},</p>
<p>Your appointment has been confirmed with the following details:</p>
<ul>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
    <li>Time: {{.TimeSlot}}</li>
    <li>Centre: {{.CentreName}}</li>
    <li>Test: {{.TestType}}</li>
</ul>
<p>Notes: {{.Notes}}</p>
<p>You can view or manage your appointment by logging into your account.</p>
{{end}}
`

const appointmentCancellationTmpl = `
{{define "appointment_cancellation"}}
<h2>Appointment Cancelled</h2>
<p>Dear {{.PatientName}},</p>
<p>Your appointment with the following details has been cancelled:</p>
<ul>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
    <li>Time: {{.TimeSlot}}</li>
    <li>Centre: {{.CentreName}}</li>
</ul>
<p>You can book a new appointment by logging into your account.</p>
{{end}}
`

const appointmentRescheduleTmpl = `
{{define "appointment_reschedule"}}
<h2>Appointment Rescheduled</h2>
<p>Dear {{.PatientName}},</p>
<p>Your appointment has been rescheduled to:</p>
<ul>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>New Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
    <li>New Time: {{.TimeSlot}}</li>
    <li>Centre: {{.CentreName}}</li>
    <li>Test: {{.TestType}}</li>
</ul>
<p>Please let us know if this new schedule does not work for you.</p>
{{end}}
`

const appointmentReminderTmpl = `
{{define "appointment_reminder"}}
<h2>Appointment Reminder</h2>
<p>Dear {{.PatientName}},</p>
<p>This is a reminder for your upcoming appointment:</p>
<ul>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
    <li>Time: {{.TimeSlot}}</li>
    <li>Centre: {{.CentreName}}</li>
    <li>Test: {{.TestType}}</li>
</ul>
<p>Please arrive 10 minutes before your scheduled time.</p>
{{end}}
`

const testResultsAvailableTmpl = `
{{define "test_results_available"}}
<h2>Test Results Available</h2>
<p>Dear {{.PatientName}},</p>
<p>Your test results are now available for the following appointment:</p>
<ul>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
    <li>Test: {{.TestType}}</li>
    <li>Centre: {{.CentreName}}</li>
</ul>
<p>To view your results securely, please:</p>
<ol>
    <li>Log in to your patient portal at: <a href="{{.ResultsPortalURL}}">{{.ResultsPortalURL}}</a></li>
    <li>Navigate to "Test Results" section</li>
    <li>Select the appointment date mentioned above</li>
</ol>
<p>If you have any questions about your results, please contact your healthcare provider.</p>
{{end}}
`

const paymentConfirmationTmpl = `
{{define "payment_confirmation"}}
<h2>Payment Confirmation</h2>
<p>Dear {{.PatientName}},</p>
<p>Thank you for your payment. Here are your transaction details:</p>
<ul>
    <li>Transaction ID: {{.TransactionID}}</li>
    <li>Amount Paid: ${{printf "%.2f" .PaymentAmount}}</li>
    <li>Payment Method: {{.PaymentMethod}} ending in {{.LastFourDigits}}</li>
    <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
</ul>
<h3>Appointment Details:</h3>
<ul>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>Test: {{.TestType}}</li>
    <li>Centre: {{.CentreName}}</li>
    <li>Time: {{.TimeSlot}}</li>
</ul>
<p>This email serves as your receipt. Please keep it for your records.</p>
{{end}}
`

const policyUpdateTmpl = `
{{define "policy_update"}}
<h2>Important Policy Update</h2>
<p>Dear {{.PatientName}},</p>
<p>We are writing to inform you about an important update to our policies:</p>
<div style="background-color: #f8f9fa; padding: 15px; margin: 10px 0;">
    <h3>{{.PolicyTitle}}</h3>
    <p>{{.PolicyDetails}}</p>
    <p><strong>Effective Date:</strong> {{.EffectiveDate.Format "Monday, January 2, 2006"}}</p>
</div>
<p>These changes are designed to improve our service and ensure the best possible care for our patients.</p>
<p>If you have any questions about these updates, please don't hesitate to contact us.</p>
{{end}}
`

const staffNotificationTmpl = `
{{define "staff_notification"}}
<h2>Staff Notification - Action Required</h2>
<p>Dear {{.StaffName}},</p>
<p><strong>Required Action:</strong> {{.RequiredAction}}</p>
<h3>Appointment Details:</h3>
<ul>
    <li>Patient: {{.PatientName}}</li>
    <li>Appointment ID: {{.AppointmentID}}</li>
    <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
    <li>Time: {{.TimeSlot}}</li>
    <li>Test: {{.TestType}}</li>
</ul>
{{if .SpecialNotes}}
<div style="background-color: #fff3cd; padding: 15px; margin: 10px 0;">
    <h4>Special Notes:</h4>
    <p>{{.SpecialNotes}}</p>
</div>
{{end}}
<p>Please ensure all necessary preparations are made before the appointment time.</p>
{{end}}
`

// Template helper functions
var templateFuncs = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("Monday, January 2, 2006")
	},
	"formatCurrency": func(amount float64) string {
		return fmt.Sprintf("$%.2f", amount)
	},
}

// Initialize templates with common layout and functions
func init() {
	baseTemplate = template.Must(template.New("base").Funcs(templateFuncs).Parse(emailBaseLayout))
}

// GetAppointmentRequestTemplate generates email template for appointment requests
func GetAppointmentRequestTemplate(data AppointmentEmailData) (string, error) {
    tmpl := template.Must(template.New("appointment_request").Parse(`
    Dear {{.PatientName}},

    We have received your appointment request at {{.CentreName}}. 
    
    Appointment Details:
    - Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}
    - Time: {{.TimeSlot}}
    - Appointment ID: {{.AppointmentID}}
    
    We will keep you updated on the status of your appointment.

    Best regards,
    {{.AppName}} Team
    `))
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}

// GetAppointmentConfirmationTemplate returns the rendered appointment confirmation email template
func GetAppointmentConfirmationTemplate(data AppointmentEmailData) (string, error) {
	contentTmpl := template.Must(template.Must(baseTemplate.Clone()).New("content").Parse(appointmentConfirmationTmpl))
	var buf bytes.Buffer
	data.Title = "Appointment Confirmation"
	if err := contentTmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetAppointmentCancellationTemplate returns the rendered appointment cancellation email template
func GetAppointmentCancellationTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("appointment_cancellation").Parse(appointmentCancellationTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetAppointmentRescheduleTemplate returns the rendered appointment reschedule email template
func GetAppointmentRescheduleTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("appointment_reschedule").Parse(appointmentRescheduleTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetAppointmentReminderTemplate returns the rendered appointment reminder email template
func GetAppointmentReminderTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("appointment_reminder").Parse(appointmentReminderTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetTestResultsAvailableTemplate returns the rendered test results available email template
func GetTestResultsAvailableTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("test_results_available").Parse(testResultsAvailableTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetPaymentConfirmationTemplate returns the rendered payment confirmation email template
func GetPaymentConfirmationTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("payment_confirmation").Parse(paymentConfirmationTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetPolicyUpdateTemplate returns the rendered policy update email template
func GetPolicyUpdateTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("policy_update").Parse(policyUpdateTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetStaffNotificationTemplate returns the rendered staff notification email template
func GetStaffNotificationTemplate(data AppointmentEmailData) (string, error) {
	tmpl := template.Must(template.New("staff_notification").Parse(staffNotificationTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
