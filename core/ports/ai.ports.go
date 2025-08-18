package ports

import "context"

type (
	AIService interface {
		InterpretLabResult(prompt string) (interface{}, error)
	}
	// OCR means Optical Character Recognition.
	OCR interface {
		Parse(ctx context.Context, imgURL []byte) ([]OCRWord, error)
	}
	// Anomaly detection model
	AnomalyDetector interface {
		Detect(ctx context.Context, data []float64) ([]string, error)
	}

	// Automated report generator
	ReportGenerator interface {
		Generate(ctx context.Context, input map[string]interface{}) (string, error)
	}

	// Decision support system
	DecisionSupporter interface {
		Recommend(ctx context.Context, patientID string) ([]string, error)
	}

	// Image analysis (DICOM, X-ray, etc.)
	ImageAnalyzer interface {
		Analyze(ctx context.Context, imagePath string) (map[string]interface{}, error)
	}

	// Package analysis (lab panels, CBC, CMP, etc.)
	PackageAnalyzer interface {
		AnalyzePackage(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
	}

	NLPResult struct {
		// normalized to your v1 schema
		CBC   map[string]float64 `json:"cbc,omitempty"`
		CMP   map[string]float64 `json:"cmp,omitempty"`
		Lipid map[string]float64 `json:"lipid,omitempty"`
		// extraction metadata
		Units      map[string]string `json:"units,omitempty"`
		Entities   []Entity          `json:"entities,omitempty"`
		Confidence float32           `json:"confidence"`
		Notes      string            `json:"notes,omitempty"`
	}
	Entity struct {
		Label, Value, Unit, Raw string
		Score                   float32
	}
	// OCRWord represents one recognized word with metadata
	OCRWord struct {
		Text       string `json:"text"`
		BBox       [4]int `json:"bbox"`       // x1, y1, x2, y2
		Confidence int    `json:"confidence"` // 0–100
	}
)
