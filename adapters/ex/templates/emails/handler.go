package emails

import (
	"bytes"
	"html/template"
)

// EmailTemplateHandler provides methods to render all email templates
type EmailTemplateHandler struct {
	baseTemplate *template.Template
	cache        *TemplateCache
}

// NewEmailTemplateHandler creates a new email template handler
func NewEmailTemplateHandler() *EmailTemplateHandler {
	base := template.New("base").Funcs(TemplateFuncs)
	base = template.Must(base.Parse(BaseLayout))

	handler := &EmailTemplateHandler{
		baseTemplate: base,
		cache:        NewTemplateCache(),
	}

	// Pre-compile all templates
	handler.cache.Compile()

	return handler
}

// RenderAppointmentConfirmation renders the appointment confirmation email
func (h *EmailTemplateHandler) RenderAppointmentConfirmation(data *AppointmentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("appointment_confirmation", appointmentConfirmationTemplate, data, EmailData{
		Title:         TitleAppointmentConfirmed,
		Icon:          IconConfirmed,
		FooterContent: FooterChanges,
	})
}

// RenderAppointmentCancellation renders the appointment cancellation email
func (h *EmailTemplateHandler) RenderAppointmentCancellation(data *AppointmentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("appointment_cancellation", appointmentCancellationTemplate, data, EmailData{
		Title:         TitleAppointmentCancelled,
		Icon:          IconCancelled,
		FooterContent: FooterSupport,
	})
}

// RenderAppointmentReminder renders the appointment reminder email
func (h *EmailTemplateHandler) RenderAppointmentReminder(data *AppointmentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("appointment_reminder", appointmentReminderTemplate, data, EmailData{
		Title:         TitleAppointmentReminder,
		Icon:          IconReminder,
		FooterContent: FooterSupport,
	})
}

// RenderPaymentConfirmation renders the payment confirmation email
func (h *EmailTemplateHandler) RenderPaymentConfirmation(data *PaymentData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("payment_confirmation", paymentConfirmationTemplate, data, EmailData{
		Title:         TitlePaymentConfirmed,
		Icon:          IconPayment,
		FooterContent: FooterPayment,
	})
}

// RenderTestResults renders the test results available email
func (h *EmailTemplateHandler) RenderTestResults(data *TestResultsData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("test_results_available", testResultsTemplate, data, EmailData{
		Title:         TitleTestResults,
		Icon:          IconTestResults,
		FooterContent: FooterResults,
	})
}

// RenderStaffNotification renders the staff notification email
func (h *EmailTemplateHandler) RenderStaffNotification(data *StaffNotificationData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("staff_notification", staffNotificationTemplate, data, EmailData{
		Title:         TitleStaffNotification,
		Icon:          IconStaff,
		FooterContent: FooterStaff,
	})
}

// RenderPolicyUpdate renders the policy update email
func (h *EmailTemplateHandler) RenderPolicyUpdate(data *PolicyUpdateData) (string, error) {
	if err := ValidateTemplateData(data); err != nil {
		return "", err
	}
	return h.renderTemplate("policy_update", policyUpdateTemplate, data, EmailData{
		Title:         TitlePolicyUpdate,
		Icon:          IconPolicy,
		FooterContent: FooterPolicy,
	})
}

// renderTemplate is a helper function to render email templates with common data
func (h *EmailTemplateHandler) renderTemplate(name, content string, data interface{}, emailData EmailData) (string, error) {
	// Try to get template from cache
	tmpl, ok := h.cache.Get(name)
	if !ok {
		// If not in cache, create and cache it
		tmpl = h.cache.GetOrSet(name, func() *template.Template {
			clone := template.Must(h.baseTemplate.Clone())
			tmpl := template.Must(clone.New(name).Parse(content))
			AddTemplateFuncs(tmpl)
			return tmpl
		})
	}

	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base", map[string]interface{}{
		"Title":         emailData.Title,
		"Icon":          emailData.Icon,
		"Content":       data,
		"FooterContent": emailData.FooterContent,
	})
	if err != nil {
		return "", NewTemplateError(name, err)
	}

	return buf.String(), nil
}
