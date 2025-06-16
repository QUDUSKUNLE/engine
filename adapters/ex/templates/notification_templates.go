package templates

// Additional fields for notification emails
type NotificationData struct {
	AppointmentEmailData
	// Fields for test results
	ResultsAvailable bool
	ResultsPortalURL string
	// Fields for payment
	PaymentAmount  float64
	TransactionID  string
	PaymentMethod  string
	LastFourDigits string
	// Fields for policy updates
	PolicyTitle   string
	PolicyDetails string
	// Fields for staff notifications
	StaffName      string
	RequiredAction string
	SpecialNotes   string
}

// const paymentConfirmationTmpl = `
// <!DOCTYPE html>
// <html>
// <body>
//     <h2>Payment Confirmation</h2>
//     <p>Dear {{.PatientName}},</p>
//     <p>Thank you for your payment. Here are your transaction details:</p>
//     <ul>
//         <li>Transaction ID: {{.TransactionID}}</li>
//         <li>Amount Paid: ${{printf "%.2f" .PaymentAmount}}</li>
//         <li>Payment Method: {{.PaymentMethod}} ending in {{.LastFourDigits}}</li>
//         <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
//     </ul>
//     <h3>Appointment Details:</h3>
//     <ul>
//         <li>Appointment ID: {{.AppointmentID}}</li>
//         <li>Test: {{.TestType}}</li>
//         <li>Centre: {{.CentreName}}</li>
//         <li>Time: {{.TimeSlot}}</li>
//     </ul>
//     <p>This email serves as your receipt. Please keep it for your records.</p>
//     <p>Best regards,<br/>{{.AppName}} Team</p>
// </body>
// </html>
// `

// const policyUpdateTmpl = `
// <!DOCTYPE html>
// <html>
// <body>
//     <h2>Important Policy Update</h2>
//     <p>Dear {{.PatientName}},</p>
//     <p>We are writing to inform you about an important update to our policies:</p>
//     <div style="background-color: #f8f9fa; padding: 15px; margin: 10px 0;">
//         <h3>{{.PolicyTitle}}</h3>
//         <p>{{.PolicyDetails}}</p>
//         <p><strong>Effective Date:</strong> {{.AppointmentDate.Format "Monday, January 2, 2006"}}</p>
//     </div>
//     <p>These changes are designed to improve our service and ensure the best possible care for our patients.</p>
//     <p>If you have any questions about these updates, please don't hesitate to contact us.</p>
//     <p>Best regards,<br/>{{.AppName}} Team</p>
// </body>
// </html>
// `

// const staffNotificationTmpl = `
// <!DOCTYPE html>
// <html>
// <body>
//     <h2>Staff Notification - Action Required</h2>
//     <p>Dear {{.StaffName}},</p>
//     <p><strong>Required Action:</strong> {{.RequiredAction}}</p>
//     <h3>Appointment Details:</h3>
//     <ul>
//         <li>Patient: {{.PatientName}}</li>
//         <li>Appointment ID: {{.AppointmentID}}</li>
//         <li>Date: {{.AppointmentDate.Format "Monday, January 2, 2006"}}</li>
//         <li>Time: {{.TimeSlot}}</li>
//         <li>Test: {{.TestType}}</li>
//     </ul>
//     {{if .SpecialNotes}}
//     <div style="background-color: #fff3cd; padding: 15px; margin: 10px 0;">
//         <h4>Special Notes:</h4>
//         <p>{{.SpecialNotes}}</p>
//     </div>
//     {{end}}
//     <p>Please ensure all necessary preparations are made before the appointment time.</p>
//     <p>Best regards,<br/>{{.AppName}} Team</p>
// </body>
// </html>
// `

// GetTestResultsAvailableTemplate returns the rendered test results available email template
// func GetTestResultsAvailableTemplate(data NotificationData) (string, error) {
// 	tmpl := template.Must(template.New("test_results_available").Parse(testResultsAvailableTmpl))
// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

// // GetPaymentConfirmationTemplate returns the rendered payment confirmation email template
// func GetPaymentConfirmationTemplate(data NotificationData) (string, error) {
// 	tmpl := template.Must(template.New("payment_confirmation").Parse(paymentConfirmationTmpl))
// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

// // GetPolicyUpdateTemplate returns the rendered policy update email template
// func GetPolicyUpdateTemplate(data NotificationData) (string, error) {
// 	tmpl := template.Must(template.New("policy_update").Parse(policyUpdateTmpl))
// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

// // GetStaffNotificationTemplate returns the rendered staff notification email template
// func GetStaffNotificationTemplate(data NotificationData) (string, error) {
// 	tmpl := template.Must(template.New("staff_notification").Parse(staffNotificationTmpl))
// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }
