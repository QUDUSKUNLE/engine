package emails

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
