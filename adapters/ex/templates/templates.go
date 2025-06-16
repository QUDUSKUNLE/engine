package templates

import (
	"html/template"
	"time"
)

var baseTemplate *template.Template

// EmailData contains the base email template data
type EmailData struct {
	AppName       string
	Title         string
	FooterContent string
}

// PasswordResetData contains data for password reset emails
type PasswordResetData struct {
	EmailData
	ResetLink string
	Token     string
	ExpiresIn string
}

type AppointmentEmailData struct {
	EmailData
	PatientName     string
	AppointmentID   string
	AppointmentDate time.Time
	TimeSlot        string
	CentreName      string
	TestType        string
	Status          string
	Notes           string
	Content         string
	// Fields for test results
	ResultsAvailable bool
	ResultsPortalURL string
	// Fields for payment
	PaymentAmount  float64
	TransactionID  string
	PaymentMethod  string
	LastFourDigits string
	// Fields for policy updates
	PolicyTitle   string
	PolicyDetails string
	EffectiveDate time.Time
	// Fields for staff notifications
	StaffName      string
	RequiredAction string
	SpecialNotes   string
}
