package emails

const passwordResetTemplate = `
{{define "password_reset"}}
<p>
	<strong>
		Dear {{.Name}},
	</strong>
</p>
<p>We received a request to reset your password for your Medivue account. Click the button below to set a new password:
</p>
<p style="text-align: center; margin: 30px 0;">
	<a href="{{.ResetLink}}" style="display: inline-block; padding: 12px 24px; background-color: #2563eb; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: bold;">Reset Password</a>
</p>
<div style="background-color: #fef2f2; border: 1px solid #fecaca; padding: 16px; border-left: 4px solid #ef4444; border-radius: 4px; margin-top: 20px;">
			<p><strong>Important:</strong> This password reset link will expire in {{.ExpiresIn}} for security reasons.</p>
			<p>If you didn't request a password reset, please ignore this email and make sure you can still log into your account.</p>
</div>
{{end}}
`
