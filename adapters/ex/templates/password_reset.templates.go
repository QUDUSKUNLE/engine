package templates

var (
	PasswordResetTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Reset Your Password - Medivue</title>
	<style>
		:root {
			--primary-color: #2563eb;
			--secondary-color: #1e40af;
			--accent-color: #60a5fa;
			--text-primary: #1f2937;
			--text-secondary: #4b5563;
			--background: #f3f4f6;
			--danger-color: #ef4444;
		}
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
			line-height: 1.6;
			color: var(--text-primary);
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: var(--background);
		}
		.container {
			background: linear-gradient(165deg, #ffffff 0%, #f0f9ff 100%);
			border-radius: 20px;
			padding: 48px;
			box-shadow: 
				0 10px 25px -5px rgba(0, 0, 0, 0.05),
				0 8px 10px -6px rgba(0, 0, 0, 0.03);
			max-width: 85%;
			box-sizing: border-box;
			margin: 20px auto;
			position: relative;
			border: 1px solid rgba(37, 99, 235, 0.1);
			backdrop-filter: blur(10px);
		}
		.header {
			text-align: left;
			margin-bottom: 40px;
			padding-bottom: 30px;
			border-bottom: 2px solid var(--accent-color);
			position: relative;
		}
		.header::after {
			content: '';
			position: absolute;
			bottom: -2px;
			left: 0;
			width: 50%;
			height: 2px;
			background: linear-gradient(90deg, var(--accent-color) 0%, transparent 100%);
		}
		.header h1 {
			color: var(--primary-color);
			margin: 0;
			font-size: 36px;
			font-weight: 800;
			letter-spacing: -0.5px;
			background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
			-webkit-background-clip: text;
			background-clip: text;
			-webkit-text-fill-color: transparent;
		}
		.content {
			margin: 32px 0;
			background: rgba(96, 165, 250, 0.04);
			padding: 35px;
			border-radius: 16px;
			border-left: 4px solid var(--primary-color);
			position: relative;
			overflow: hidden;
		}
		.content::before {
			content: '';
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			height: 2px;
			background: linear-gradient(90deg, var(--primary-color), transparent);
		}
		.button {
			display: inline-block;
			padding: 18px 36px;
			background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
			color: white;
			text-decoration: none;
			border-radius: 14px;
			margin: 30px 0;
			font-weight: 600;
			font-size: 16px;
			letter-spacing: 0.3px;
			transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
			box-shadow: 
				0 4px 6px -1px rgba(37, 99, 235, 0.25),
				0 2px 4px -1px rgba(37, 99, 235, 0.15),
				inset 0 1px 0 rgba(255, 255, 255, 0.1);
		}
		.button:hover {
			transform: translateY(-2px);
			box-shadow: 
				0 12px 20px -5px rgba(37, 99, 235, 0.4),
				0 8px 10px -6px rgba(37, 99, 235, 0.2),
				inset 0 1px 0 rgba(255, 255, 255, 0.2);
		}
		.icon {
			display: inline-block;
			vertical-align: middle;
			margin-right: 8px;
			font-size: 20px;
		}
		.note {
			background: rgba(239, 68, 68, 0.05);
			border-radius: 12px;
			padding: 20px;
			margin: 24px 0;
			border: 1px solid rgba(239, 68, 68, 0.1);
			position: relative;
		}
		.note::before {
			content: '';
			position: absolute;
			top: 0;
			left: 0;
			width: 4px;
			height: 100%;
			background: linear-gradient(180deg, var(--danger-color), transparent);
			border-radius: 4px 0 0 4px;
		}
		.footer {
			color: var(--text-secondary);
			font-size: 14px;
			margin-top: 45px;
			padding-top: 30px;
			border-top: 1px solid rgba(37, 99, 235, 0.08);
			position: relative;
		}
		.footer::before {
			content: '';
			position: absolute;
			top: -1px;
			left: 0;
			width: 30%;
			height: 1px;
			background: linear-gradient(90deg, var(--accent-color), transparent);
		}
		.small-text {
			color: var(--text-secondary);
			font-size: 13px;
			line-height: 1.6;
			padding: 12px;
			background: rgba(96, 165, 250, 0.02);
			border-radius: 8px;
			margin-top: 16px;
		}
		.link-text {
			color: var(--accent-color);
			word-break: break-all;
		}
		@media (max-width: 640px) {
			.container {
				padding: 32px 24px;
				max-width: 95%;
			}
			.header h1 {
				font-size: 28px;
			}
			.button {
				width: 100%;
				text-align: center;
				padding: 16px 24px;
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Reset Your Password</h1>
		</div>
		<div class="content">
			<p><strong>Hi there,</strong></p>
			<p>We received a request to reset your password for your Medivue account. Click the button below to set a new password:</p>
			<button style="text-align: left;">
				<a href="%[1]s/v1/reset_password?token=%[2]s&email=%[3]s" class="button">Reset Password</a>
			</button>
			<div class="note">
				<p><strong>Important:</strong> This password reset link will expire in 15 minutes for security reasons.</p>
				<p style="margin-bottom: 0;">If you didn't request a password reset, please ignore this email and make sure you can still log into your account.</p>
			</div>
		</div>
		<div class="footer">
			<p>Best regards,<br><strong>The Medivue Team</strong></p>
			<div class="small-text">
				<p style="margin: 0;">If you're having trouble with the button, copy and paste this URL into your web browser:</p>
				<p class="link-text" style="margin: 8px 0;">%[1]s/v1/reset_password?token=%[2]s&email=%[3]s</p>
			</div>
		</div>
	</div>
</body>
</html>
`
)
