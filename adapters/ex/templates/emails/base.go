package emails

import (
	"html/template"
	"time"
)

// Common email template styles and layout
const BaseLayout = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - DiagnoxixAI</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #1f2937;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f3f4f6;
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
        }
        .header {
            text-align: left;
            margin-bottom: 40px;
            padding-bottom: 30px;
            border-bottom: 2px solid #60a5fa;
            position: relative;
        }
        .header::after {
            content: '';
            position: absolute;
            bottom: -2px;
            left: 0;
            width: 50%;
            height: 2px;
            background: linear-gradient(90deg, #60a5fa 0%, transparent 100%);
        }
        .header h1, .header h2, .header h3 {
            color: #ffffff;
            margin: 0;
            font-size: 36px;
            font-weight: 800;
            letter-spacing: -0.5px;
            background: linear-gradient(135deg, #2563eb, #1e40af);
            -webkit-background-clip: text;
            background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .content {
            margin: 32px 0;
            background: rgba(96, 165, 250, 0.04);
            padding: 35px;
            border-radius: 16px;
            border-left: 4px solid #2563eb;
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
            background: linear-gradient(90deg, #2563eb, transparent);
        }
        .details {
            background: rgba(96, 165, 250, 0.05);
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            border: 1px solid rgba(37, 99, 235, 0.1);
        }
        .details ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
        .details li {
            padding: 10px 0;
            border-bottom: 1px solid rgba(37, 99, 235, 0.1);
        }
        .details li:last-child {
            border-bottom: none;
        }
        .note {
            background: rgba(96, 165, 250, 0.05);
            border-radius: 8px;
            padding: 16px;
            margin: 20px 0;
            border: 1px solid rgba(37, 99, 235, 0.1);
        }
        .footer {
            color: #4b5563;
            font-size: 14px;
            margin-top: 40px;
            padding-top: 25px;
            border-top: 1px solid rgba(37, 99, 235, 0.1);
        }
        .small-text {
            color: #4b5563;
            font-size: 13px;
            line-height: 1.5;
        }
        .button {
            background: #2563eb;
            color: white;
            text-decoration: none;
            padding: 14px 28px;
            border-radius: 6px;
            display: inline-block;
            font-weight: 600;
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h3>{{.Icon}} {{.Title}}</h3>
        </div>
        <div class="content">
            {{if eq .Type "email_verification" }}
                {{template "email_verification" .Data}}
            {{else if eq .Type "dc_management_notification" }}
                {{template "dc_management_notification" .Data}}
            {{else if eq .Type "appointment_cancellation"}}
                {{template "appointment_cancellation" .Data}}
            {{else if eq .Type "dc_manager_notification"}} 
                {{template "dc_manager_notification" .Data}}
            {{else if eq .Type "password_reset"}}
                {{template "password_reset" .Data}}
            {{end}}
        </div>
        <div class="footer">
            <p>Best regards,<br><strong>The DiagnoxixAI Team</strong></p>
            <div class="small-text">
                {{.FooterContent}}
            </div>
        </div>
    </div>
</body>
</html>
`

// Template helper functions
var TemplateFuncs = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("Monday, January 2, 2006")
	},
	"formatTime": func(t time.Time) string {
		return t.Format("3:04 PM")
	},
	"formatCurrency": func(amount float64) string {
		return "₦%.2f" // Using Naira symbol for Nigerian currency
	},
}
