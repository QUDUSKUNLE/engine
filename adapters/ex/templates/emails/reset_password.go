package emails

const passwordResetTmpl = `
<!DOCTYPE html>
<html>
<body>
	<h2>{{.Greeting}}</h2>
	<p>You requested a password reset for your {{.AppName}} account.</p>
	<p>Click the link below to reset your password:</p>
	<p><a href="{{.ResetLink}}">Reset Password</a></p>
	<p>This link will expire in {{.ExpiresIn}}.</p>
</body>
</html>
`
