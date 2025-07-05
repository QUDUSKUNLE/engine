package emails

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// TimeFormat formats time in 12-hour format with AM/PM
func TimeFormat(t time.Time) string {
	return t.Format("3:04 PM")
}

// DateFormat formats date in full format (Monday, January 2, 2006)
func DateFormat(t time.Time) string {
	return t.Format("Monday, January 2, 2006")
}

// FormatNaira formats amount in Naira with appropriate symbol and decimals
func FormatNaira(amount float64) string {
	return fmt.Sprintf("₦%.2f", amount)
}

// FormatTestType formats test type for display (e.g., BLOOD_TEST -> Blood Test)
func FormatTestType(testType string) string {
	words := strings.Split(strings.ToLower(testType), "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

// TruncateString truncates a string to the specified length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// AddTemplateFuncs adds custom template functions to the provided template
func AddTemplateFuncs(t *template.Template) *template.Template {
	t.Funcs(template.FuncMap{
		"formatTime":  TimeFormat,
		"formatDate":  DateFormat,
		"formatNaira": FormatNaira,
		"formatTest":  FormatTestType,
		"truncate":    TruncateString,
		"toLower":     strings.ToLower,
		"toUpper":     strings.ToUpper,
		"title":       strings.Title,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"contains":    strings.Contains,
		"replace":     strings.Replace,
		"split":       strings.Split,
		"join":        strings.Join,
		"timeAdd":     func(t time.Time, d string) time.Time { dur, _ := time.ParseDuration(d); return t.Add(dur) },
		"timeSub":     func(t time.Time, d string) time.Time { dur, _ := time.ParseDuration(d); return t.Add(-dur) },
		"timeUntil":   time.Until,
		"timeSince":   time.Since,
		"defaultString": func(val, defaultVal string) string {
			if val == "" {
				return defaultVal
			}
			return val
		},
	})
	return t
}

// Common template strings
const (
	ContactSupportFooter = "If you need any assistance, please contact our support team."
	UrgentContactFooter  = "For urgent matters, please call our 24/7 support line."
	DefaultFooter        = "Best regards,\nThe Medivue Team"
)
