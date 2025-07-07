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

// Compile pre-compiles and caches all templates
func (c *TemplateCache) Compile() error {
	templates := map[string]string{
		TemplateAppointmentConfirmed:       appointmentConfirmationTemplate,
		TemplateAppointmentCancelled:       appointmentCancellationTemplate,
		TemplateAppointmentReminder:        appointmentReminderTemplate,
		TemplateAppointmentReschedule:      appointmentRescheduleTemplate,
		TemplatePaymentConfirmation:        paymentConfirmationTemplate,
		TemplateTestResults:                testResultsTemplate,
		TemplateStaffNotification:          staffNotificationTemplate,
		TemplatePolicyUpdate:               policyUpdateTemplate,
		TemplateEmailVerification:          emailVerificationTemplate,
		TemplateResetPassword:              passwordResetTemplate,
		TemplateDiagnosticCentreManager:    diagnosticCentreManagerEmailVerificationTemplate,
		TemplateDiagnosticCentreManagement: diagnosticCentreManagerNotificationTemplate,
	}

	base := template.New("base").Funcs(TemplateFuncs)
	base = template.Must(base.Parse(BaseLayout))

	for name, content := range templates {
		tmpl := template.Must(base.Clone())
		tmpl = template.Must(tmpl.New(name).Parse(content))
		AddTemplateFuncs(tmpl)
		c.Set(name, tmpl)
	}

	return nil
}
