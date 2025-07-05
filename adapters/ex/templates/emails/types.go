package emails

import "time"

// EmailData represents the base email template data
type EmailData struct {
	AppName       string
	Title         string
	Header        string
	Icon          string
	Content       interface{}
	FooterContent string
}

// EmailVerificationData contains fields for email verification emails
type EmailVerificationData struct {
	// EmailData
	Name             string
	VerificationLink string
	ExpiryDuration   string
}

// AppointmentData contains common fields for appointment-related emails
type AppointmentData struct {
	EmailData
	PatientName     string
	AppointmentID   string
	AppointmentDate time.Time
	TimeSlot        string
	CentreName      string
	TestType        string
	Notes           string
}

// PaymentData contains fields for payment-related emails
type PaymentData struct {
	EmailData
	PatientName     string
	TransactionID   string
	PaymentAmount   float64
	PaymentMethod   string
	LastFourDigits  string
	PaymentDate     time.Time
	AppointmentID   string
	AppointmentDate time.Time
	TimeSlot        string
	CentreName      string
	TestType        string
}

// TestResultsData contains fields for test results emails
type TestResultsData struct {
	EmailData
	PatientName      string
	AppointmentID    string
	AppointmentDate  time.Time
	TestType         string
	CentreName       string
	ResultsPortalURL string
}

// StaffNotificationData contains fields for staff notification emails
type StaffNotificationData struct {
	EmailData
	StaffName           string
	PatientName         string
	AppointmentID       string
	AppointmentDate     time.Time
	TimeSlot            string
	TestType            string
	SpecialInstructions string
	CentreName          string
	SpecialNotes        string
	RequiredAction      string
}

// PolicyUpdateData contains fields for policy update emails
type PolicyUpdateData struct {
	EmailData
	PatientName    string
	PolicyTitle    string
	PolicyDetails  string
	EffectiveDate  time.Time
	ActionRequired string
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

// Additional fields for notification emails
type NotificationData struct {
	EmailData
	// AppointmentEmailData
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
	// Fields for staff notifications
	StaffName      string
	RequiredAction string
	SpecialNotes   string
}
