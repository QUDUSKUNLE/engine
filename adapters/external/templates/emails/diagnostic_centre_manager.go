package emails

const DiagnosticCentreManagerEmailVerificationTemplate = `
{{define "dc_manager_notification"}}
    <p><strong>Dear {{.ManagerName}},</strong></p>

    <p>Thank you for registering with DiagnoxixAI. We’re excited to have you onboard as a Diagnostic Centre Manager!</p>

    <p>Below are your login credentials:</p>

    <div class="details">
        <ul>
            <li><strong>Email:</strong> {{.Email}}</li>
            <li><strong>Password:</strong> {{.Password}}</li>
        </ul>
    </div>

    <p>🔒 <strong>Security Tip:</strong> Please change your password after your first login for added security.</p>
{{end}}
`
