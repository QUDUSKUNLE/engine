package domain

type (
	LabTest struct {
		PatientID          string            `json:"patient_id" validate:"required"`
		DiagnosticCentreID string            `json:"diagnostic_centre_id" validate:"required"`
		TestName           string            `json:"test_name" validate:"required"`
		Results            map[string]string `json:"results" validate:"required"`
		ReferenceRanges    map[string]string `json:"reference_ranges" validate:"required"`
	}
)