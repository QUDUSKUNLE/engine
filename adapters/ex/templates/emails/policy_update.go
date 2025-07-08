package emails

const policyUpdateTemplate = `
{{define "policy_update"}}
<p><strong>Dear {{.PatientName}},</strong></p>

<p>We are writing to inform you of an important update to our policies:</p>

<div class="details">
    <h3 style="margin-top: 0; color: var(--primary-color);">{{.PolicyTitle}}</h3>
    <p>{{.PolicyDetails}}</p>
    <p><strong>Effective Date:</strong> {{.EffectiveDate | formatDate}}</p>
</div>

<div class="note">
    <p>These updates are designed to enhance our service and ensure better care for all patients.</p>
    {{if .ActionRequired}}
    <p><strong>Action Required:</strong> {{.ActionRequired}}</p>
    {{end}}
</div>

<p>Thank you for being a valued part of the Medivue community.</p>
{{end}}
`
