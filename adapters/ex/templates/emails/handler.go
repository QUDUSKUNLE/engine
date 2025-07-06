package emails

import (
	"bytes"
	"fmt"
	"html/template"
)

// EmailTemplateHandler provides methods to render all email templates
type EmailTemplateHandler struct {
	baseTemplate *template.Template
	cache        *TemplateCache
}

// NewEmailTemplateHandler creates a new email template handler
func NewEmailTemplateHandler() *EmailTemplateHandler {
	base := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))

	handler := &EmailTemplateHandler{
		baseTemplate: base,
		cache:        NewTemplateCache(),
	}
	// Pre-compile all templates
	handler.cache.Compile()
	return handler
}

func (h *EmailTemplateHandler) ExecuteTemplate(templateName string, data interface{}) (string, error) {
	switch templateName {
	case TemplateAppointmentConfirmed:
		return h.renderAppointmentConfirmation(data.(*AppointmentData))
	case TemplateAppointmentCancelled:
		return h.renderAppointmentCancellation(data.(*AppointmentData))
	case TemplateAppointmentReminder:
		return h.renderAppointmentReminder(data.(*AppointmentData))
	case TemplatePaymentConfirmation:
		return h.renderPaymentConfirmation(data.(*PaymentData))
	case TemplateEmailVerification:
		return h.renderEmailVerification(data.(*EmailVerificationData))
	case TemplateStaffNotification:
		return h.renderStaffNotification(data.(*StaffNotificationData))
	case TemplateTestResults:
		return h.renderTestResults(data.(*TestResultsData))
	case TemplatePolicyUpdate:
		return h.renderPolicyUpdate(data.(*PolicyUpdateData))
	case TemplateResetPassword:
		return h.renderResetPassword(data.(*PasswordResetData))
	default:
		return "", fmt.Errorf("unknown template: %s", templateName)
	}
}

// RenderAppointmentConfirmation renders the appointment confirmation email
func (h *EmailTemplateHandler) renderAppointmentConfirmation(data *AppointmentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateAppointmentConfirmed,
		data,
		EmailData{
			Title:         TitleAppointmentConfirmed,
			Icon:          IconConfirmed,
			FooterContent: FooterChanges,
		})
}

// RenderAppointmentCancellation renders the appointment cancellation email
func (h *EmailTemplateHandler) renderAppointmentCancellation(data *AppointmentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateAppointmentCancelled,
		data,
		EmailData{
			Title:         TitleAppointmentCancelled,
			Icon:          IconCancelled,
			FooterContent: FooterSupport,
		})
}

// RenderAppointmentReminder renders the appointment reminder email
func (h *EmailTemplateHandler) renderAppointmentReminder(data *AppointmentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateAppointmentReminder,
		data,
		EmailData{
			Title:         TitleAppointmentReminder,
			Icon:          IconReminder,
			FooterContent: FooterSupport,
		})
}

// RenderPaymentConfirmation renders the payment confirmation email
func (h *EmailTemplateHandler) renderPaymentConfirmation(data *PaymentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplatePaymentConfirmation,
		data,
		EmailData{
			Title:         TitlePaymentConfirmed,
			Icon:          IconPayment,
			FooterContent: FooterPayment,
		})
}

// RenderTestResults renders the test results available email
func (h *EmailTemplateHandler) renderTestResults(data *TestResultsData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateTestResults,
		data,
		EmailData{
			Title:         TitleTestResults,
			Icon:          IconTestResults,
			FooterContent: FooterResults,
		})
}

// RenderStaffNotification renders the staff notification email
func (h *EmailTemplateHandler) renderStaffNotification(data *StaffNotificationData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateStaffNotification,
		data, EmailData{
			Title:         TitleStaffNotification,
			Icon:          IconStaff,
			FooterContent: FooterStaff,
		})
}

// RenderPolicyUpdate renders the policy update email
func (h *EmailTemplateHandler) renderPolicyUpdate(data *PolicyUpdateData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplatePolicyUpdate,
		data,
		EmailData{
			Title:         TitlePolicyUpdate,
			Icon:          IconPolicy,
			FooterContent: FooterPolicy,
		})
}

// RenderEmailVerification renders the email verification email
func (h *EmailTemplateHandler) renderEmailVerification(data *EmailVerificationData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateEmailVerification,
		data,
		EmailData{
			Title:         TitleEmailVerification,
			Icon:          IconEmailVerification,
			FooterContent: FooterEmailVerification,
		})
}

// renderResetPassword renders the reset password email
func (h *EmailTemplateHandler) renderResetPassword(data *PasswordResetData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate(
		TemplateResetPassword,
		data,
		EmailData{
			Title:         TitleResetPassword,
			Icon:          IconResetPassword,
			FooterContent: FooterResetPassord,
		})
}

func (h *EmailTemplateHandler) renderTemplate(
	name string,
	data interface{},
	emailData EmailData,
) (string, error) {
	temp, _ := h.cache.Get(name)
	var contentBuf bytes.Buffer
	err := temp.ExecuteTemplate(&contentBuf, name, data)
	if err != nil {
		return "", NewTemplateError(name+" (content rendering)", err)
	}

	var fullBuf bytes.Buffer
	err = temp.ExecuteTemplate(&fullBuf, "base", map[string]interface{}{
		"Title":         emailData.Title,
		"Icon":          emailData.Icon,
		"Content":       template.HTML(contentBuf.String()),
		"FooterContent": emailData.FooterContent,
	})

	if err != nil {
		return "", NewTemplateError(name+"(base rendering)", err)
	}
	return fullBuf.String(), nil
}
