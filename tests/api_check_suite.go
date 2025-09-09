package api_check_suite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

// Test configuration
type TestConfig struct {
	BaseURL string
	Token   string
}

// Test results tracking
type TestResult struct {
	Name     string
	Endpoint string
	Status   string
	Duration time.Duration
	Error    string
}

type TestSuite struct {
	config  TestConfig
	results []TestResult
	client  *http.Client
}

func RunTests() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize test suite
	suite := &TestSuite{
		config: TestConfig{
			BaseURL: getEnv("TEST_BASE_URL", "http://localhost:7556"),
			Token:   getEnv("TEST_JWT_TOKEN", ""), // You'll need to provide a valid JWT token
		},
		client: &http.Client{Timeout: 30 * time.Second},
	}

	fmt.Println("🧪 Starting Diagnoxix API Test Suite")
	fmt.Println("=====================================")

	// Run all tests
	suite.runAllTests()

	// Print results
	suite.printResults()
}

func (suite *TestSuite) runAllTests() {
	// Test 1: Health Check
	suite.testHealthCheck()

	// Test 2: AI Capabilities
	suite.testAICapabilities()

	// Test 3: Lab Result Interpretation
	suite.testLabResultInterpretation()

	// Test 4: Symptom Analysis
	suite.testSymptomAnalysis()

	// Test 5: Report Summarization
	suite.testReportSummarization()

	// Test 6: Medical Image Analysis
	suite.testMedicalImageAnalysis()

	// Test 7: Anomaly Detection
	suite.testAnomalyDetection()

	// Test 8: Lab Package Analysis
	suite.testLabPackageAnalysis()

	// Test 9: Automated Report Generation
	suite.testAutomatedReportGeneration()

	// Test 10: WebSocket Stats
	suite.testWebSocketStats()
}

func (suite *TestSuite) testHealthCheck() {
	start := time.Now()

	resp, err := suite.client.Get(suite.config.BaseURL + "/health")
	duration := time.Since(start)

	result := TestResult{
		Name:     "Health Check",
		Endpoint: "/health",
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAIL"
		result.Error = err.Error()
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			result.Status = "PASS"
		} else {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Status: %d", resp.StatusCode)
		}
	}

	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testAICapabilities() {
	start := time.Now()

	resp, err := suite.client.Get(suite.config.BaseURL + "/v1/ai/capabilities")
	duration := time.Since(start)

	result := TestResult{
		Name:     "AI Capabilities",
		Endpoint: "/v1/ai/capabilities",
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAIL"
		result.Error = err.Error()
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var response map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				result.Status = "FAIL"
				result.Error = "Invalid JSON response"
			} else {
				features, ok := response["features"].([]interface{})
				if ok && len(features) >= 7 {
					result.Status = "PASS"
				} else {
					result.Status = "FAIL"
					result.Error = "Expected at least 7 AI features"
				}
			}
		} else {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Status: %d", resp.StatusCode)
		}
	}

	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testLabResultInterpretation() {
	start := time.Now()

	payload := map[string]interface{}{
		"patient_id":           "test-patient-123",
		"diagnostic_centre_id": "test-center-456",
		"test_name":            "Complete Blood Count",
		"results": map[string]string{
			"Hemoglobin":        "8.5 g/dL",
			"White Blood Cells": "12,000 /μL",
			"Platelets":         "150,000 /μL",
		},
		"reference_ranges": map[string]string{
			"Hemoglobin":        "12.0-15.5 g/dL",
			"White Blood Cells": "4,000-11,000 /μL",
			"Platelets":         "150,000-450,000 /μL",
		},
	}

	result := suite.makePostRequest("Lab Result Interpretation", "/v1/ai/interpret-lab", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testSymptomAnalysis() {
	start := time.Now()

	payload := map[string]interface{}{
		"symptoms": []string{"persistent cough", "shortness of breath", "chest pain"},
		"age":      45,
		"gender":   "male",
	}

	result := suite.makePostRequest("Symptom Analysis", "/v1/ai/analyze-symptoms", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testReportSummarization() {
	start := time.Now()

	payload := map[string]interface{}{
		"medical_report": `CHEST X-RAY REPORT
		
Patient: John Doe, Age: 45, Male
Date: 2024-01-15

FINDINGS:
- Bilateral lower lobe infiltrates consistent with pneumonia
- No pleural effusion
- Heart size within normal limits

IMPRESSION:
Bilateral pneumonia. Clinical correlation recommended.`,
		"patient_friendly": true,
	}

	result := suite.makePostRequest("Report Summarization", "/v1/ai/summarize-report", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testMedicalImageAnalysis() {
	start := time.Now()

	payload := map[string]interface{}{
		"image_url":      "https://example.com/chest-xray.jpg",
		"image_type":     "XRAY",
		"body_part":      "chest",
		"patient_age":    45,
		"patient_gender": "male",
	}

	result := suite.makePostRequest("Medical Image Analysis", "/v1/ai/analyze-medical-image", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testAnomalyDetection() {
	start := time.Now()

	payload := map[string]interface{}{
		"data":      []float64{120, 80, 98.6, 72, 16, 98},
		"data_type": "vital_signs",
		"patient_profile": map[string]interface{}{
			"age":             45,
			"gender":          "male",
			"medical_history": []string{"hypertension", "diabetes"},
			"medications":     []string{"metformin", "lisinopril"},
		},
	}

	result := suite.makePostRequest("Anomaly Detection", "/v1/ai/detect-anomalies", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testLabPackageAnalysis() {
	start := time.Now()

	payload := map[string]interface{}{
		"package_data": map[string]interface{}{
			"package_type": "Comprehensive Metabolic Panel",
			"patient_profile": map[string]interface{}{
				"age":             35,
				"gender":          "female",
				"medical_history": []string{"hypothyroidism"},
				"medications":     []string{"levothyroxine"},
			},
			"test_results": map[string]interface{}{
				"Glucose":    "95 mg/dL",
				"BUN":        "18 mg/dL",
				"Creatinine": "0.9 mg/dL",
				"Sodium":     "140 mEq/L",
				"Potassium":  "4.2 mEq/L",
			},
			"reference_ranges": map[string]string{
				"Glucose":    "70-100 mg/dL",
				"BUN":        "7-20 mg/dL",
				"Creatinine": "0.6-1.2 mg/dL",
				"Sodium":     "136-145 mEq/L",
				"Potassium":  "3.5-5.0 mEq/L",
			},
			"test_date": "2025-01-15",
		},
	}

	result := suite.makePostRequest("Lab Package Analysis", "/v1/ai/analyze-lab-package", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testAutomatedReportGeneration() {
	start := time.Now()

	payload := map[string]interface{}{
		"report_data": map[string]interface{}{
			"report_type":     "Comprehensive Health Assessment",
			"report_purpose":  "Annual physical examination",
			"target_audience": "patient",
			"patient_info": map[string]interface{}{
				"age":             42,
				"gender":          "male",
				"medical_history": []string{"hypertension", "high cholesterol"},
				"medications":     []string{"atorvastatin", "amlodipine"},
			},
			"test_results": []map[string]interface{}{
				{
					"test_name":       "Total Cholesterol",
					"value":           "195",
					"unit":            "mg/dL",
					"reference_range": "<200 mg/dL",
					"status":          "normal",
				},
			},
			"clinical_data": map[string]interface{}{
				"weight":         "185 lbs",
				"height":         "5'10\"",
				"bmi":            "26.5",
				"smoking_status": "former smoker",
			},
		},
	}

	result := suite.makePostRequest("Automated Report Generation", "/v1/ai/generate-report", payload, start)
	suite.results = append(suite.results, result)
}

func (suite *TestSuite) testWebSocketStats() {
	start := time.Now()

	resp, err := suite.client.Get(suite.config.BaseURL + "/v1/ws/stats")
	duration := time.Since(start)

	result := TestResult{
		Name:     "WebSocket Stats",
		Endpoint: "/v1/ws/stats",
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAIL"
		result.Error = err.Error()
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			result.Status = "PASS"
		} else {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Status: %d", resp.StatusCode)
		}
	}

	suite.results = append(suite.results, result)
}

func (suite *TestSuite) makePostRequest(name, endpoint string, payload interface{}, start time.Time) TestResult {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return TestResult{
			Name:     name,
			Endpoint: endpoint,
			Status:   "FAIL",
			Duration: time.Since(start),
			Error:    "Failed to marshal JSON: " + err.Error(),
		}
	}

	req, err := http.NewRequest("POST", suite.config.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return TestResult{
			Name:     name,
			Endpoint: endpoint,
			Status:   "FAIL",
			Duration: time.Since(start),
			Error:    "Failed to create request: " + err.Error(),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	if suite.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+suite.config.Token)
	}

	resp, err := suite.client.Do(req)
	duration := time.Since(start)

	result := TestResult{
		Name:     name,
		Endpoint: endpoint,
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAIL"
		result.Error = err.Error()
		return result
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		result.Status = "PASS"
	} else if resp.StatusCode == http.StatusUnauthorized {
		result.Status = "SKIP"
		result.Error = "Authentication required (provide TEST_JWT_TOKEN)"
	} else {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, string(body))
	}

	return result
}

func (suite *TestSuite) printResults() {
	fmt.Println("\n🧪 Test Results Summary")
	fmt.Println("=======================")

	passed := 0
	failed := 0
	skipped := 0

	for _, result := range suite.results {
		status := result.Status
		icon := "❌"
		if status == "PASS" {
			icon = "✅"
			passed++
		} else if status == "SKIP" {
			icon = "⏭️"
			skipped++
		} else {
			failed++
		}

		fmt.Printf("%s %s (%s) - %v\n", icon, result.Name, result.Endpoint, result.Duration)
		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}
	}

	fmt.Println("\n📊 Summary:")
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("⏭️ Skipped: %d\n", skipped)
	fmt.Printf("📈 Total: %d\n", len(suite.results))

	if failed == 0 {
		fmt.Println("\n🎉 All tests passed! Your API is working correctly.")
	} else {
		fmt.Printf("\n⚠️ %d tests failed. Please check the errors above.\n", failed)
	}
}

