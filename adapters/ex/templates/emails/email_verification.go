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
<a href="{{.VerificationLink}}" style="
background: var(--primary-color);
	color:  #2563eb;
	text-decoration: none;
	padding: 14px 28px;
	border-radius: 6px;
	display: inline-block;
	font-weight: 600;
	box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);">
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
