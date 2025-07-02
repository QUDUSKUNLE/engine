package templates

var (
	DiagnosticCentreManagerNotificationTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Diagnostic Centre Management Notification - Medivue</title>
	<style>
		:root {
			--primary-color: #2563eb;
			--secondary-color: #1e40af;
			--accent-color: #60a5fa;
			--text-primary: #1f2937;
			--text-secondary: #4b5563;
			--background: #f3f4f6;
			--success-color: #10b981;
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
			background: rgba(96, 165, 250, 0.05);
			border-radius: 8px;
			padding: 16px;
			margin: 20px 0;
			border: 1px solid rgba(37, 99, 235, 0.1);
		}
		.footer {
			color: var(--text-secondary);
			font-size: 14px;
			margin-top: 40px;
			padding-top: 25px;
			border-top: 1px solid rgba(37, 99, 235, 0.1);
		}
		.small-text {
			color: var(--text-secondary);
			font-size: 13px;
			line-height: 1.5;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Management Notification!</h1>
		</div>
		<div class="content">
			<p><strong>Hi there,</strong></p>
			<p>You have been assigned a Manager Role in the diagnostic centre below:</p>
		<p><strong>Diagnostic Centre Name:</strong> %[1]s</p>
		<p><strong>Address:</strong> %[2]s</p>
		<div>
			Wish your the very best
		</div>
		</div>
		<div class="footer">
			<p>Best regards,<br><strong>The Medivue Team</strong></p>
		</div>
	</div>
</body>
</html>
`
)
