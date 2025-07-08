package emails

const passwordResetTemplate = `
{{define "password_reset"}}
<p><strong>Dear {{.Name}},</strong></p>

<p>We received a request to reset your Medivue account password. Click the button below to set a new password:</p>

<p style="text-align: center; margin: 30px 0;">
    <a href="{{.ResetLink}}" style="display: inline-block; padding: 12px 24px; background-color: var(--primary-color); color: #2563eb; text-decoration: none; border-radius: 6px; font-weight: bold;">
        Reset Password
    </a>
</p>

<div class="note">
    <p><strong>Important:</strong> This password reset link will expire in {{.ExpiresIn}} for security reasons.</p>
    <p>If you didn’t request a password reset, please ignore this email and ensure your account security.</p>
</div>
{{end}}
`
