package emails

const emailVerificationTemplate = `
{{define "email_verification"}}
<p>
	<strong>
		Dear {{.Name}},
	</strong>
</p>
<p>
	Thank you for registering with Medivue. Please verify your email address to complete your registration.
</p>
<p>
	Click the button below to verify your email address:
</p>
<a href="{{.VerificationLink}}" class="button" style="background: var(--primary-color); color: blue; text-decoration: none; padding: 12px 24px; border-radius: 4px; display: inline-block; margin: 20px 0;">
		Verify Email Address
</a>
<p>
	If the button doesn't work, you can copy and paste this link into your browser:
</p>
<p style="word-break: break-all;">
	{{.VerificationLink}}
</p>
<p>
	<strong>
		Note:
	</strong>
		This verification link will expire in {{.ExpiryDuration}}.
</p>
<p>
	If you did not create an account with us, please ignore this email.
</p>
{{end}}
`
