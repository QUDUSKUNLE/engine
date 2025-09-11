package cache

import (
	"time"
)

// CacheTTLConfig defines TTL for different AI operations
type (
		CacheTTLConfig struct {
		LabInterpretation    time.Duration `yaml:"lab_interpretation" default:"2h"`
		SymptomAnalysis      time.Duration `yaml:"symptom_analysis" default:"1h"`
		ReportSummary        time.Duration `yaml:"report_summary" default:"4h"`
		ImageAnalysis        time.Duration `yaml:"image_analysis" default:"6h"`
		AnomalyDetection     time.Duration `yaml:"anomaly_detection" default:"30m"`
		LabPackageAnalysis   time.Duration `yaml:"lab_package_analysis" default:"3h"`
		AutomatedReport      time.Duration `yaml:"automated_report" default:"8h"`
	}
	// CacheKeyConfig defines key patterns for different operations
	CacheKeyConfig struct {
		LabInterpretation  string `yaml:"lab_interpretation" default:"lab:interp"`
		SymptomAnalysis    string `yaml:"symptom_analysis" default:"symptom:analysis"`
		ReportSummary      string `yaml:"report_summary" default:"report:summary"`
		ImageAnalysis      string `yaml:"image_analysis" default:"image:analysis"`
		AnomalyDetection   string `yaml:"anomaly_detection" default:"anomaly:detect"`
		LabPackageAnalysis string `yaml:"lab_package_analysis" default:"lab:package"`
		AutomatedReport    string `yaml:"automated_report" default:"report:auto"`
	}
)


// DefaultCacheTTLConfig returns default TTL configuration
func DefaultCacheTTLConfig() CacheTTLConfig {
	return CacheTTLConfig{
		LabInterpretation:    2 * time.Hour,
		SymptomAnalysis:      1 * time.Hour,
		ReportSummary:        4 * time.Hour,
		ImageAnalysis:        6 * time.Hour,
		AnomalyDetection:     30 * time.Minute,
		LabPackageAnalysis:   3 * time.Hour,
		AutomatedReport:      8 * time.Hour,
	}
}

// GetTTLForOperation returns the appropriate TTL for an AI operation
func (c CacheTTLConfig) GetTTLForOperation(operation string) time.Duration {
	switch operation {
	case "lab_interpretation":
		return c.LabInterpretation
	case "symptom_analysis":
		return c.SymptomAnalysis
	case "report_summary":
		return c.ReportSummary
	case "image_analysis":
		return c.ImageAnalysis
	case "anomaly_detection":
		return c.AnomalyDetection
	case "lab_package_analysis":
		return c.LabPackageAnalysis
	case "automated_report":
		return c.AutomatedReport
	default:
		return 1 * time.Hour // Default TTL
	}
}

// DefaultCacheKeyConfig returns default key configuration
func DefaultCacheKeyConfig() CacheKeyConfig {
	return CacheKeyConfig{
		LabInterpretation:  "lab:interp",
		SymptomAnalysis:    "symptom:analysis",
		ReportSummary:      "report:summary",
		ImageAnalysis:      "image:analysis",
		AnomalyDetection:   "anomaly:detect",
		LabPackageAnalysis: "lab:package",
		AutomatedReport:    "report:auto",
	}
}

// GetKeyPrefixForOperation returns the appropriate key prefix for an operation
func (c CacheKeyConfig) GetKeyPrefixForOperation(operation string) string {
	switch operation {
	case "lab_interpretation":
		return c.LabInterpretation
	case "symptom_analysis":
		return c.SymptomAnalysis
	case "report_summary":
		return c.ReportSummary
	case "image_analysis":
		return c.ImageAnalysis
	case "anomaly_detection":
		return c.AnomalyDetection
	case "lab_package_analysis":
		return c.LabPackageAnalysis
	case "automated_report":
		return c.AutomatedReport
	default:
		return "ai:unknown"
	}
}
