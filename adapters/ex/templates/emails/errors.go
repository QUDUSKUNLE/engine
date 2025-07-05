package emails

import (
	"fmt"
)

// TemplateError represents an error that occurred during template rendering
type TemplateError struct {
	Template string
	Err      error
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("template error in %s: %v", e.Template, e.Err)
}

// ValidationError represents an error in template data validation
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// NewTemplateError creates a new template error
func NewTemplateError(template string, err error) error {
	return &TemplateError{
		Template: template,
		Err:      err,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ValidateTemplateData validates the required fields in template data
func ValidateTemplateData(data interface{}) error {
	switch v := data.(type) {
	case *AppointmentData:
		if v == nil {
			return NewValidationError("AppointmentData", "data cannot be nil")
		}
		if v.PatientName == "" {
			return NewValidationError("PatientName", "cannot be empty")
		}
		if v.AppointmentID == "" {
			return NewValidationError("AppointmentID", "cannot be empty")
		}
		if v.AppointmentDate.IsZero() {
			return NewValidationError("AppointmentDate", "cannot be zero")
		}
	case *PaymentData:
		if v == nil {
			return NewValidationError("PaymentData", "data cannot be nil")
		}
		if v.TransactionID == "" {
			return NewValidationError("TransactionID", "cannot be empty")
		}
		if v.PaymentAmount <= 0 {
			return NewValidationError("PaymentAmount", "must be greater than zero")
		}
	case *TestResultsData:
		if v == nil {
			return NewValidationError("TestResultsData", "data cannot be nil")
		}
		if v.PatientName == "" {
			return NewValidationError("PatientName", "cannot be empty")
		}
		if v.ResultsPortalURL == "" {
			return NewValidationError("ResultsPortalURL", "cannot be empty")
		}
	case *StaffNotificationData:
		if v == nil {
			return NewValidationError("StaffNotificationData", "data cannot be nil")
		}
		if v.StaffName == "" {
			return NewValidationError("StaffName", "cannot be empty")
		}
		if v.PatientName == "" {
			return NewValidationError("PatientName", "cannot be empty")
		}
	}
	return nil
}
