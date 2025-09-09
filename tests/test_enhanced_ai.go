package api_check_suite

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/diagnoxix/core/services"
	"github.com/joho/godotenv"
)

func TestEnhancedAIFeatures() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	openAIKey := os.Getenv("OPEN_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPEN_API_KEY environment variable is required")
	}

	// Initialize AI service
	aiService := services.NewAIService(openAIKey)

	fmt.Println("🤖 Testing Enhanced AI Features...")

	// Test 1: Medical Image Analysis
	fmt.Println("\n=== Testing Medical Image Analysis ===")
	testMedicalImageAnalysis(aiService)

	// Test 2: Anomaly Detection
	fmt.Println("\n=== Testing Anomaly Detection ===")
	testAnomalyDetection(aiService)

	// Test 3: Lab Package Analysis
	fmt.Println("\n=== Testing Lab Package Analysis ===")
	testLabPackageAnalysis(aiService)

	// Test 4: Automated Report Generation
	fmt.Println("\n=== Testing Automated Report Generation ===")
	testAutomatedReportGeneration(aiService)
}

func testMedicalImageAnalysis(aiService *services.AIService) {
	analysis, err := aiService.AnalyzeMedicalImage(
		context.Background(),
		"https://example.com/chest-xray.jpg",
		"XRAY",
		"chest",
		45,
		"male",
	)
	if err != nil {
		log.Printf("Medical image analysis failed: %v", err)
		return
	}

	result, _ := json.MarshalIndent(analysis, "", "  ")
	fmt.Printf("Medical Image Analysis Result:\n%s\n", string(result))
}

func testAnomalyDetection(aiService *services.AIService) {
	// Sample vital signs data with some anomalies
	vitalSigns := []float64{120, 80, 98.6, 72, 16, 98} // BP sys, BP dia, temp, HR, RR, O2 sat
	
	patientProfile := services.PatientProfile{
		Age:            45,
		Gender:         "male",
		MedicalHistory: []string{"hypertension", "diabetes"},
		Medications:    []string{"metformin", "lisinopril"},
	}

	result, err := aiService.DetectAnomalies(
		context.Background(),
		vitalSigns,
		"vital_signs",
		patientProfile,
	)
	if err != nil {
		log.Printf("Anomaly detection failed: %v", err)
		return
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("Anomaly Detection Result:\n%s\n", string(resultJSON))
}

func testLabPackageAnalysis(aiService *services.AIService) {
	packageData := services.LabPackageData{
		PackageType: "Comprehensive Metabolic Panel",
		PatientProfile: services.PatientProfile{
			Age:            35,
			Gender:         "female",
			MedicalHistory: []string{"hypothyroidism"},
			Medications:    []string{"levothyroxine"},
		},
		TestResults: map[string]interface{}{
			"Glucose":    "95 mg/dL",
			"BUN":        "18 mg/dL",
			"Creatinine": "0.9 mg/dL",
			"Sodium":     "140 mEq/L",
			"Potassium":  "4.2 mEq/L",
			"Chloride":   "102 mEq/L",
			"CO2":        "24 mEq/L",
			"Calcium":    "9.8 mg/dL",
			"Protein":    "7.2 g/dL",
			"Albumin":    "4.1 g/dL",
			"Bilirubin":  "0.8 mg/dL",
			"ALT":        "25 U/L",
			"AST":        "22 U/L",
		},
		ReferenceRanges: map[string]string{
			"Glucose":    "70-100 mg/dL",
			"BUN":        "7-20 mg/dL",
			"Creatinine": "0.6-1.2 mg/dL",
			"Sodium":     "136-145 mEq/L",
			"Potassium":  "3.5-5.0 mEq/L",
			"Chloride":   "98-107 mEq/L",
			"CO2":        "22-28 mEq/L",
			"Calcium":    "8.5-10.5 mg/dL",
			"Protein":    "6.0-8.3 g/dL",
			"Albumin":    "3.5-5.0 g/dL",
			"Bilirubin":  "0.3-1.2 mg/dL",
			"ALT":        "7-56 U/L",
			"AST":        "10-40 U/L",
		},
		TestDate: "2025-01-15",
	}

	analysis, err := aiService.AnalyzeLabPackage(context.Background(), packageData)
	if err != nil {
		log.Printf("Lab package analysis failed: %v", err)
		return
	}

	result, _ := json.MarshalIndent(analysis, "", "  ")
	fmt.Printf("Lab Package Analysis Result:\n%s\n", string(result))
}

func testAutomatedReportGeneration(aiService *services.AIService) {
	reportData := services.ReportGenerationData{
		ReportType:     "Comprehensive Health Assessment",
		ReportPurpose:  "Annual physical examination",
		TargetAudience: "patient",
		PatientInfo: services.PatientProfile{
			Age:            42,
			Gender:         "male",
			MedicalHistory: []string{"hypertension", "high cholesterol"},
			Medications:    []string{"atorvastatin", "amlodipine"},
		},
		TestResults: []services.TestResult{
			{
				TestName:       "Total Cholesterol",
				Value:          "195",
				Unit:           "mg/dL",
				ReferenceRange: "<200 mg/dL",
				Status:         "normal",
			},
			{
				TestName:       "LDL Cholesterol",
				Value:          "115",
				Unit:           "mg/dL",
				ReferenceRange: "<100 mg/dL",
				Status:         "slightly elevated",
			},
			{
				TestName:       "HDL Cholesterol",
				Value:          "45",
				Unit:           "mg/dL",
				ReferenceRange: ">40 mg/dL",
				Status:         "normal",
			},
			{
				TestName:       "Blood Pressure",
				Value:          "135/85",
				Unit:           "mmHg",
				ReferenceRange: "<120/80 mmHg",
				Status:         "elevated",
			},
		},
		ClinicalData: map[string]interface{}{
			"weight":        "185 lbs",
			"height":        "5'10\"",
			"bmi":           "26.5",
			"smoking_status": "former smoker",
			"exercise":      "moderate",
		},
	}

	report, err := aiService.GenerateAutomatedReport(context.Background(), reportData)
	if err != nil {
		log.Printf("Automated report generation failed: %v", err)
		return
	}

	result, _ := json.MarshalIndent(report, "", "  ")
	fmt.Printf("Automated Report Result:\n%s\n", string(result))
}
