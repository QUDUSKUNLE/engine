package emails

import (
	"html/template"
	"sync"
)

// TemplateCache provides thread-safe caching for parsed templates
type TemplateCache struct {
	cache sync.Map
}

// NewTemplateCache creates a new template cache
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{}
}

// Get retrieves a template from the cache
func (c *TemplateCache) Get(name string) (*template.Template, bool) {
	if tmpl, ok := c.cache.Load(name); ok {
		return tmpl.(*template.Template), true
	}
	return nil, false
}

// Set stores a template in the cache
func (c *TemplateCache) Set(name string, tmpl *template.Template) {
	c.cache.Store(name, tmpl)
}

// Delete removes a template from the cache
func (c *TemplateCache) Delete(name string) {
	c.cache.Delete(name)
}

// Clear removes all templates from the cache
func (c *TemplateCache) Clear() {
	c.cache = sync.Map{}
}

// GetOrSet retrieves a template from the cache or sets it if not found
func (c *TemplateCache) GetOrSet(name string, creator func() *template.Template) *template.Template {
	if tmpl, ok := c.Get(name); ok {
		return tmpl
	}

	tmpl := creator()
	c.Set(name, tmpl)
	return tmpl
}

func (c *TemplateCache) Compile() error {
	// Combine all templates in one go
	allTemplates := BaseLayout +
		appointmentConfirmationTemplate +
		appointmentCancellationTemplate +
		appointmentReminderTemplate +
		appointmentRescheduleTemplate +
		paymentConfirmationTemplate +
		testResultsTemplate +
		staffNotificationTemplate +
		policyUpdateTemplate +
		emailVerificationTemplate +
		passwordResetTemplate +
		diagnosticCentreManagerEmailVerificationTemplate +
		diagnosticCentreManagerNotificationTemplate

	base := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(allTemplates))
	AddTemplateFuncs(base)

	// Store each by its name
	c.Set(TemplateEmailVerification, base.Lookup(TemplateEmailVerification))
	c.Set(TemplateResetPassword, base.Lookup(TemplateResetPassword))
	c.Set(TemplateAppointmentConfirmed, base.Lookup(TemplateAppointmentConfirmed))
	c.Set(TemplateAppointmentCancelled, base.Lookup(TemplateAppointmentCancelled))
	c.Set(TemplateAppointmentReminder, base.Lookup(TemplateAppointmentReminder))
	c.Set(TemplateAppointmentReschedule, base.Lookup(TemplateAppointmentReschedule))
	c.Set(TemplatePaymentConfirmation, base.Lookup(TemplatePaymentConfirmation))
	c.Set(TemplateTestResults, base.Lookup(TemplateTestResults))
	c.Set(TemplatePolicyUpdate, base.Lookup(TemplatePolicyUpdate))
	c.Set(TemplateStaffNotification, base.Lookup(TemplateStaffNotification))
	c.Set(TemplateDiagnosticCentreManager, base.Lookup(TemplateDiagnosticCentreManager))
	c.Set(TemplateDiagnosticCentreManagement, base.Lookup(TemplateDiagnosticCentreManagement))

	return nil
}

