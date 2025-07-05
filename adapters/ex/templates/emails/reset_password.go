package emails

import (
	"bytes"
	"html/template"
)

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

// GetPasswordResetTemplate returns the rendered password reset email template
func GetPasswordResetTemplate(data PasswordResetData) (string, error) {
	tmpl := template.Must(template.New("password_reset").Parse(passwordResetTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
