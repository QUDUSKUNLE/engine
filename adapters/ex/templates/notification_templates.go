package templates

// Additional fields for notification emails
type NotificationData struct {
	AppointmentEmailData
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
