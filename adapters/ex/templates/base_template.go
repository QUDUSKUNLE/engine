package templates

import (
	"html/template"
)

const baseTemplateLayout = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen-Sans, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: #4A90E2;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
            margin-bottom: 20px;
        }
        .content {
            background: #fff;
            padding: 20px;
            border-radius: 0 0 8px 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 14px;
            margin-top: 20px;
            padding: 20px;
        }
        ul, ol {
            margin: 15px 0;
            padding-left: 30px;
        }
        li {
            margin: 10px 0;
        }
        .highlight {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-radius: 4px;
            border-left: 4px solid #4A90E2;
        }
        .warning {
            background: #fff3cd;
            padding: 15px;
            margin: 10px 0;
            border-radius: 4px;
            border-left: 4px solid #ffc107;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            background: #4A90E2;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            margin: 10px 0;
        }
        @media screen and (max-width: 480px) {
            body {
                padding: 10px;
            }
            .header, .content, .footer {
                padding: 15px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>{{.Title}}</h2>
    </div>
    <div class="content">
        {{.Content}}
    </div>
    <div class="footer">
        <p>Best regards,<br/>{{.AppName}} Team</p>
        {{if .FooterContent}}
        <p>{{.FooterContent}}</p>
        {{end}}
    </div>
</body>
</html>
`

func init() {
	var err error
	baseTemplate, err = template.New("base").Parse(baseTemplateLayout)
	if err != nil {
		panic(err)
	}
}
