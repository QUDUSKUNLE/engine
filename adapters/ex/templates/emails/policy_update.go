package emails

import (
	"bytes"
	"html/template"
)

const policyUpdateTemplate = `
{{define "policy_update"}}
<p><strong>Dear {{.PatientName}},</strong></p>
<p>We are writing to inform you about an important update to our policies:</p>

<div class="details">
    <h3>{{.PolicyTitle}}</h3>
    <p>{{.PolicyDetails}}</p>
    <p><strong>Effective Date:</strong> {{.EffectiveDate | formatDate}}</p>
</div>

<div class="note">
    <p>These changes are designed to improve our service and ensure the best possible care for our patients.</p>
    {{if .ActionRequired}}
    <p><strong>Action Required:</strong> {{.ActionRequired}}</p>
    {{end}}
</div>
{{end}}
`

// GetPolicyUpdateTemplate returns the rendered policy update email
func GetPolicyUpdateTemplate(data PolicyUpdateData) (string, error) {
	baseTemplate := template.Must(template.New("base").Funcs(TemplateFuncs).Parse(BaseLayout))
	contentTemplate := template.Must(baseTemplate.New("content").Parse(policyUpdateTemplate))

	var buf bytes.Buffer
	err := contentTemplate.ExecuteTemplate(&buf, "base", map[string]interface{}{
		"Title":         data.Title,
		"Icon":          data.Icon,
		"Content":       data.Content,
		"FooterContent": data.FooterContent,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
