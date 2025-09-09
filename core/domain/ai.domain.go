package domain

type (
	LabTest struct {
		PatientID          string            `json:"patient_id" validate:"required"`
		DiagnosticCentreID string            `json:"diagnostic_centre_id" validate:"required"`
		TestName           string            `json:"test_name" validate:"required"`
		Results            map[string]string `json:"results" validate:"required"`
		ReferenceRanges    map[string]string `json:"reference_ranges" validate:"required"`
	}
	SymptomAnalysisRequest struct {
		Symptoms []string `json:"symptoms" validate:"required,min=1"`
		Age      int      `json:"age" validate:"required,min=1,max=120"`
		Gender   string   `json:"gender" validate:"required,oneof=male female other"`
	}

	ReportSummaryRequest struct {
		MedicalReport   string `json:"medical_report" validate:"required"`
		PatientFriendly bool   `json:"patient_friendly"`
	}
	// MedicalImageAnalysisRequest represents a medical image analysis request
	MedicalImageAnalysisRequest struct {
		ImageURL      string `json:"image_url" validate:"required,url"`
		ImageType     string `json:"image_type" validate:"required,oneof=XRAY CT_SCAN MRI ULTRASOUND MAMMOGRAM"`
		BodyPart      string `json:"body_part" validate:"required"`
		PatientAge    int    `json:"patient_age" validate:"required,min=1,max=120"`
		PatientGender string `json:"patient_gender" validate:"required,oneof=male female"`
	}
)
