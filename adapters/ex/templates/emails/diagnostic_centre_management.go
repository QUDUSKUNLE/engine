package emails

const diagnosticCentreManagerNotificationTemplate = `
{{define "dc_management_notification"}}
<p><strong>Dear {{.Name}},</strong></p>
<p>You have been assigned a Manager Role in the diagnostic centre below:</p>
<div class="details">
  <ul>
    <li><strong>Diagnostic Centre Name:</strong> {{.CentreName}}</li>
    <li><strong>Address:</strong> {{.CentreAddress}}</li>
  </ul>
</div>
<p style="margin-top: 16px;">Wish you the very best.</p>
{{end}}
`
