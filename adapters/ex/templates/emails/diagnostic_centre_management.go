package emails

const diagnosticCentreManagerNotificationTemplate = `
{{define "dc_management_notification"}}
<div style="background: linear-gradient(165deg, #ffffff 0%, #f0f9ff 100%); border-radius: 20px; padding: 48px; box-shadow: 0 10px 25px -5px rgba(0,0,0,0.05), 0 8px 10px -6px rgba(0,0,0,0.03); max-width: 85%; margin: 20px auto; border: 1px solid rgba(37, 99, 235, 0.1); backdrop-filter: blur(10px);">

    <div style="text-align: left; margin-bottom: 40px; padding-bottom: 30px; border-bottom: 2px solid #60a5fa; position: relative;">
      <h1 style="margin: 0; font-size: 36px; font-weight: 800; letter-spacing: -0.5px; background: linear-gradient(135deg, #2563eb, #1e40af); -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent;">
        🎉 Management Notification!
      </h1>
    </div>

    <div style="margin: 32px 0; background: rgba(96, 165, 250, 0.04); padding: 35px; border-radius: 16px; border-left: 4px solid #2563eb;">
      <p><strong>Dear {{.Name}},</strong></p>
      <p>You have been assigned a Manager Role in the diagnostic centre below:</p>
      <p><strong>Diagnostic Centre Name:</strong> {{.CentreName}}</p>
      <p><strong>Address:</strong> {{.CentreAddress}}</p>
      <p style="margin-top: 16px;">Wish you the very best.</p>
    </div>
  </div>
{{end}}
`
