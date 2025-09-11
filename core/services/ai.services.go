package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/services/cache"
	"github.com/diagnoxix/core/utils"
	"github.com/google/uuid"
)

type (
	AIService struct {
		openAIKey   string
		client      *http.Client
		cache       *cache.AICache
		cacheConfig cache.CacheTTLConfig
		keyConfig   cache.CacheKeyConfig
		metrics     *cache.CacheMetrics
	}
	OpenAIRequest struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float64   `json:"temperature"`
		MaxTokens   int       `json:"max_tokens"`
	}
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	
	OpenAIResponse struct {
		Choices []Choice `json:"choices"`
		Usage   Usage    `json:"usage"`
	}
	
	Choice struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	}
	
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
	
	LabInterpretation struct {
		Summary           string                   `json:"summary"`
		AbnormalResults   []AbnormalResult         `json:"abnormal_results"`
		Recommendations   []string                 `json:"recommendations"`
		UrgencyLevel      string                   `json:"urgency_level"`
		FollowUpRequired  bool                     `json:"follow_up_required"`
		DetailedAnalysis  map[string]TestAnalysis  `json:"detailed_analysis"`
	}
	
	AbnormalResult struct {
		TestName     string `json:"test_name"`
		Value        string `json:"value"`
		ReferenceRange string `json:"reference_range"`
		Severity     string `json:"severity"`
		Explanation  string `json:"explanation"`
	}
	
	TestAnalysis struct {
		Status      string `json:"status"`
		Explanation string `json:"explanation"`
		Implications string `json:"implications"`
	}
	
	SymptomAnalysis struct {
		PossibleConditions []PossibleCondition `json:"possible_conditions"`
		UrgencyLevel       string              `json:"urgency_level"`
		Recommendations    []string            `json:"recommendations"`
		RedFlags           []string            `json:"red_flags"`
		NextSteps          []string            `json:"next_steps"`
	}
	
	PossibleCondition struct {
		Name        string  `json:"name"`
		Probability string  `json:"probability"`
		Description string  `json:"description"`
		Symptoms    []string `json:"symptoms"`
	}
	// Additional AI service structures for enhanced funct
	MedicalImageAnalysis struct {
		ImageType       string                 `json:"image_type"`
		BodyPart        string                 `json:"body_part"`
		Findings        []ImageFinding         `json:"findings"`
		Abnormalities   []ImageAbnormality     `json:"abnormalities"`
		Measurements    []ImageMeasurement     `json:"measurements"`
		Recommendations []string               `json:"recommendations"`
		UrgencyLevel    string                 `json:"urgency_level"`
		Confidence      float64                `json:"confidence"`
		RequiresReview  bool                   `json:"requires_review"`
	}
	ImageFinding struct {
		Finding     string  `json:"finding"`
		Location    string  `json:"location"`
		Severity    string  `json:"severity"`
		Description string  `json:"description"`
		Confidence  float64 `json:"confidence"`
	}
	
	ImageAbnormality struct {
		Type        string  `json:"type"`
		Location    string  `json:"location"`
		Size        string  `json:"size,omitempty"`
		Description string  `json:"description"`
		Significance string  `json:"significance"`
	}
	
	ImageMeasurement struct {
		Parameter string  `json:"parameter"`
		Value     string  `json:"value"`
		Unit      string  `json:"unit"`
		Normal    bool    `json:"normal"`
	}
	
	AnomalyDetectionResult struct {
		AnomaliesDetected bool              `json:"anomalies_detected"`
		Anomalies         []DetectedAnomaly `json:"anomalies"`
		OverallRisk       string            `json:"overall_risk"`
		Recommendations   []string          `json:"recommendations"`
		Confidence        float64           `json:"confidence"`
		DataQuality       string            `json:"data_quality"`
	}
	DetectedAnomaly struct {
		DataPoint   string  `json:"data_point"`
		Value       float64 `json:"value"`
		ExpectedRange string `json:"expected_range"`
		Severity    string  `json:"severity"`
		Description string  `json:"description"`
	}
	
	PatientProfile struct {
		Age           int      `json:"age"`
		Gender        string   `json:"gender"`
		MedicalHistory []string `json:"medical_history,omitempty"`
		Medications   []string `json:"medications,omitempty"`
		Allergies     []string `json:"allergies,omitempty"`
	}
	
	LabPackageData struct {
		PackageType    string                 `json:"package_type"`
		PatientProfile PatientProfile         `json:"patient_profile"`
		TestResults    map[string]interface{} `json:"test_results"`
		ReferenceRanges map[string]string     `json:"reference_ranges"`
		TestDate       string                 `json:"test_date"`
	}
	
	LabPackageAnalysis struct {
		PackageType     string                    `json:"package_type"`
		OverallAssessment string                  `json:"overall_assessment"`
		KeyFindings     []PackageFinding          `json:"key_findings"`
		SystemAnalysis  map[string]SystemAnalysis `json:"system_analysis"`
		Correlations    []TestCorrelation         `json:"correlations"`
		Recommendations []string                  `json:"recommendations"`
		FollowUpTests   []string                  `json:"follow_up_tests"`
		RiskFactors     []string                  `json:"risk_factors"`
	}
	PackageFinding struct {
		Category    string `json:"category"`
		Finding     string `json:"finding"`
		Significance string `json:"significance"`
		Impact      string `json:"impact"`
	}
	
	SystemAnalysis struct {
		System      string   `json:"system"`
		Status      string   `json:"status"`
		Findings    []string `json:"findings"`
		Concerns    []string `json:"concerns"`
	}
	
	TestCorrelation struct {
		Tests       []string `json:"tests"`
		Relationship string   `json:"relationship"`
		Significance string   `json:"significance"`
	}
	
	ReportGenerationData struct {
		ReportType     string                 `json:"report_type"`
		PatientInfo    PatientProfile         `json:"patient_info"`
		TestResults    []TestResult           `json:"test_results"`
		ClinicalData   map[string]interface{} `json:"clinical_data"`
		ReportPurpose  string                 `json:"report_purpose"`
		TargetAudience string                 `json:"target_audience"`
	}
	
	TestResult struct {
		TestName    string `json:"test_name"`
		Value       string `json:"value"`
		Unit        string `json:"unit,omitempty"`
		ReferenceRange string `json:"reference_range,omitempty"`
		Status      string `json:"status"`
	}
	AutomatedReport struct {
		ReportID       string                 `json:"report_id"`
		ReportType     string                 `json:"report_type"`
		PatientSummary string                 `json:"patient_summary"`
		ExecutiveSummary string               `json:"executive_summary"`
		DetailedFindings []ReportSection      `json:"detailed_findings"`
		Recommendations []string              `json:"recommendations"`
		Conclusions     string                `json:"conclusions"`
		Metadata        map[string]interface{} `json:"metadata"`
		GeneratedAt     string                `json:"generated_at"`
	}
	
	ReportSection struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		Findings    []string `json:"findings"`
		Significance string  `json:"significance"`
	}
)

func NewAIService(openAIKey string) *AIService {
	return &AIService{
		openAIKey: openAIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheConfig: cache.DefaultCacheTTLConfig(),
		keyConfig:   cache.DefaultCacheKeyConfig(),
		metrics:     cache.NewCacheMetrics(),
	}
}

// NewAIServiceWithCache creates a new AI service with caching enabled
func NewAIServiceWithCache(openAIKey string, aiCache *cache.AICache) *AIService {
	return &AIService{
		openAIKey: openAIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache:       aiCache,
		cacheConfig: cache.DefaultCacheTTLConfig(),
		keyConfig:   cache.DefaultCacheKeyConfig(),
		metrics:     cache.NewCacheMetrics(),
	}
}

// InterpretLabResults analyzes lab test results and provides medical interpretation
func (ai *AIService) InterpretLabResults(ctx context.Context, labTest domain.LabTest) (*LabInterpretation, error) {
	operation := "lab_interpretation"
	
	// Try cache first if available
	if ai.cache != nil {
		cacheKey := ai.generateCacheKey(operation, labTest)
		start := time.Now()
		
		if cached, found, err := ai.cache.Get(ctx, cacheKey); err == nil && found {
			ai.metrics.RecordHit(time.Since(start))
			utils.Info("Cache hit for lab interpretation", utils.LogField{Key: "cache_key", Value: cacheKey})
			
			if interpretation, ok := cached.(*LabInterpretation); ok {
				return interpretation, nil
			}
		} else if err == nil {
			ai.metrics.RecordMiss(time.Since(start))
		} else {
			ai.metrics.RecordError(time.Since(start))
			utils.Warn("Cache error for lab interpretation", utils.LogField{Key: "error", Value: err.Error()})
		}
	}

	// Cache miss or no cache - call OpenAI
	prompt := ai.buildLabInterpretationPrompt(labTest)
	
	response, err := ai.callOpenAI(ctx, prompt, "You are a medical AI assistant specializing in lab result interpretation. Provide accurate, helpful analysis while emphasizing the need for professional medical consultation.")
	if err != nil {
		utils.Error("Failed to call OpenAI for lab interpretation", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	interpretation, err := ai.parseLabInterpretation(response)
	if err != nil {
		utils.Error("Failed to parse lab interpretation", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Cache the result if cache is available
	if ai.cache != nil {
		cacheKey := ai.generateCacheKey(operation, labTest)
		ttl := ai.cacheConfig.GetTTLForOperation(operation)
		start := time.Now()
		
		if err := ai.cache.Set(ctx, cacheKey, interpretation, ttl); err != nil {
			ai.metrics.RecordError(time.Since(start))
			utils.Warn("Failed to cache lab interpretation", utils.LogField{Key: "error", Value: err.Error()})
		} else {
			ai.metrics.RecordSet(time.Since(start))
			utils.Info("Cached lab interpretation", 
				utils.LogField{Key: "cache_key", Value: cacheKey},
				utils.LogField{Key: "ttl", Value: ttl.String()})
		}
	}

	return interpretation, nil
}

// AnalyzeSymptoms provides preliminary analysis of patient symptoms
func (ai *AIService) AnalyzeSymptoms(ctx context.Context, symptoms []string, patientAge int, patientGender string) (*SymptomAnalysis, error) {
	operation := "symptom_analysis"
	input := map[string]interface{}{
		"symptoms": symptoms,
		"age":      patientAge,
		"gender":   patientGender,
	}
	
	// Try cache first if available
	if ai.cache != nil {
		cacheKey := ai.generateCacheKey(operation, input)
		start := time.Now()
		
		if cached, found, err := ai.cache.Get(ctx, cacheKey); err == nil && found {
			ai.metrics.RecordHit(time.Since(start))
			utils.Info("Cache hit for symptom analysis", utils.LogField{Key: "cache_key", Value: cacheKey})
			
			if analysis, ok := cached.(*SymptomAnalysis); ok {
				return analysis, nil
			}
		} else if err == nil {
			ai.metrics.RecordMiss(time.Since(start))
		} else {
			ai.metrics.RecordError(time.Since(start))
			utils.Warn("Cache error for symptom analysis", utils.LogField{Key: "error", Value: err.Error()})
		}
	}

	// Cache miss or no cache - call OpenAI
	prompt := ai.buildSymptomAnalysisPrompt(symptoms, patientAge, patientGender)
	
	response, err := ai.callOpenAI(ctx, prompt, "You are a medical AI assistant for preliminary symptom analysis. Always emphasize that this is not a diagnosis and professional medical consultation is required.")
	if err != nil {
		utils.Error("Failed to call OpenAI for symptom analysis", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("symptom analysis failed: %w", err)
	}

	analysis, err := ai.parseSymptomAnalysis(response)
	if err != nil {
		utils.Error("Failed to parse symptom analysis", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Cache the result if cache is available
	if ai.cache != nil {
		cacheKey := ai.generateCacheKey(operation, input)
		ttl := ai.cacheConfig.GetTTLForOperation(operation)
		start := time.Now()
		
		if err := ai.cache.Set(ctx, cacheKey, analysis, ttl); err != nil {
			ai.metrics.RecordError(time.Since(start))
			utils.Warn("Failed to cache symptom analysis", utils.LogField{Key: "error", Value: err.Error()})
		} else {
			ai.metrics.RecordSet(time.Since(start))
			utils.Info("Cached symptom analysis", 
				utils.LogField{Key: "cache_key", Value: cacheKey},
				utils.LogField{Key: "ttl", Value: ttl.String()})
		}
	}

	return analysis, nil
}

// GenerateReportSummary creates patient-friendly summaries of medical reports
func (ai *AIService) GenerateReportSummary(ctx context.Context, medicalReport string, patientFriendly bool) (string, error) {
	var prompt string
	if patientFriendly {
		prompt = fmt.Sprintf(`Please create a patient-friendly summary of this medical report. Use simple language, explain medical terms, and organize the information clearly:

Medical Report:
%s

Please provide:
1. A brief overview in simple terms
2. Key findings explained clearly
3. What this means for the patient
4. Any recommended next steps

Keep the tone reassuring but honest, and always recommend consulting with their healthcare provider for questions.`, medicalReport)
	} else {
		prompt = fmt.Sprintf(`Please create a concise professional summary of this medical report for healthcare providers:

Medical Report:
%s

Focus on:
1. Key clinical findings
2. Significant abnormalities
3. Clinical implications
4. Recommended follow-up`, medicalReport)
	}

	systemMessage := "You are a medical AI assistant specializing in creating clear, accurate medical report summaries."
	
	response, err := ai.callOpenAI(ctx, prompt, systemMessage)
	if err != nil {
		utils.Error("Failed to generate report summary", utils.LogField{Key: "error", Value: err.Error()})
		return "", fmt.Errorf("report summary generation failed: %w", err)
	}

	return response, nil
}

func (ai *AIService) callOpenAI(ctx context.Context, prompt, systemMessage string) (string, error) {
	reqBody := OpenAIRequest{
		Model: "gpt-4o-mini", // Cost-effective model for medical analysis
		Messages: []Message{
			{
				Role:    "system",
				Content: systemMessage,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.3, // Lower temperature for more consistent medical analysis
		MaxTokens:   1500,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.openAIKey)

	resp, err := ai.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (ai *AIService) buildLabInterpretationPrompt(labTest domain.LabTest) string {
	var resultsStr strings.Builder
	for testName, value := range labTest.Results {
		refRange := labTest.ReferenceRanges[testName]
		resultsStr.WriteString(fmt.Sprintf("- %s: %s (Reference: %s)\n", testName, value, refRange))
	}

	return fmt.Sprintf(`Please analyze these lab test results and provide a comprehensive interpretation:

Patient ID: %s
Test Name: %s
Results:
%s

Please provide your analysis in the following JSON format:
{
  "summary": "Brief overall summary of the results",
  "abnormal_results": [
    {
      "test_name": "name of abnormal test",
      "value": "actual value",
      "reference_range": "normal range",
      "severity": "mild/moderate/severe",
      "explanation": "what this means"
    }
  ],
  "recommendations": ["list of recommendations"],
  "urgency_level": "low/medium/high",
  "follow_up_required": true/false,
  "detailed_analysis": {
    "test_name": {
      "status": "normal/abnormal",
      "explanation": "detailed explanation",
      "implications": "clinical implications"
    }
  }
}

Important: This is for informational purposes only and should not replace professional medical consultation.`,
		labTest.PatientID,
		labTest.TestName,
		resultsStr.String())
}

func (ai *AIService) buildSymptomAnalysisPrompt(symptoms []string, age int, gender string) string {
	symptomsStr := strings.Join(symptoms, ", ")
	
	return fmt.Sprintf(`Please analyze these symptoms for a %d-year-old %s patient:

Symptoms: %s

Please provide your analysis in the following JSON format:
{
  "possible_conditions": [
    {
      "name": "condition name",
      "probability": "low/medium/high",
      "description": "brief description",
      "symptoms": ["matching symptoms"]
    }
  ],
  "urgency_level": "low/medium/high/emergency",
  "recommendations": ["general recommendations"],
  "red_flags": ["warning signs to watch for"],
  "next_steps": ["recommended actions"]
}

Important: This is preliminary analysis only. Always recommend professional medical evaluation.`,
		age, gender, symptomsStr)
}

func (ai *AIService) parseLabInterpretation(response string) (*LabInterpretation, error) {
	// Extract JSON from response if it's wrapped in markdown or other text
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var interpretation LabInterpretation
	if err := json.Unmarshal([]byte(jsonStr), &interpretation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lab interpretation: %w", err)
	}
	
	return &interpretation, nil
}

func (ai *AIService) parseSymptomAnalysis(response string) (*SymptomAnalysis, error) {
	// Extract JSON from response if it's wrapped in markdown or other text
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var analysis SymptomAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal symptom analysis: %w", err)
	}
	
	return &analysis, nil
}

// Medical Image Analysis using AI
func (ai *AIService) AnalyzeMedicalImage(ctx context.Context, imageURL, imageType, bodyPart string, patientAge int, patientGender string) (*MedicalImageAnalysis, error) {

	prompt := ai.buildMedicalImagePrompt(imageURL, imageType, bodyPart, patientAge, patientGender)
	
	response, err := ai.callOpenAI(ctx, prompt, "You are a medical AI assistant specializing in medical image analysis. Provide accurate analysis while emphasizing the need for professional radiologist consultation.")
	if err != nil {
		utils.Error("Failed to call OpenAI for medical image analysis", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("medical image analysis failed: %w", err)
	}

	analysis, err := ai.parseMedicalImageAnalysis(response)
	if err != nil {
		utils.Error("Failed to parse medical image analysis", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return analysis, nil
}

// Anomaly Detection in medical data
func (ai *AIService) DetectAnomalies(ctx context.Context, data []float64, dataType string, patientProfile PatientProfile) (*AnomalyDetectionResult, error) {
	prompt := ai.buildAnomalyDetectionPrompt(data, dataType, patientProfile)
	
	response, err := ai.callOpenAI(ctx, prompt, "You are a medical AI assistant specializing in anomaly detection in medical data. Identify unusual patterns that may require medical attention.")
	if err != nil {
		utils.Error("Failed to call OpenAI for anomaly detection", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("anomaly detection failed: %w", err)
	}

	result, err := ai.parseAnomalyDetectionResult(response)
	if err != nil {
		utils.Error("Failed to parse anomaly detection result", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}

// Package Analysis for comprehensive lab panels
func (ai *AIService) AnalyzeLabPackage(ctx context.Context, packageData LabPackageData) (*LabPackageAnalysis, error) {
	prompt := ai.buildLabPackagePrompt(packageData)
	
	response, err := ai.callOpenAI(ctx, prompt, "You are a medical AI assistant specializing in comprehensive lab package analysis. Provide holistic interpretation of multiple related tests.")
	if err != nil {
		utils.Error("Failed to call OpenAI for lab package analysis", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("lab package analysis failed: %w", err)
	}

	analysis, err := ai.parseLabPackageAnalysis(response)
	if err != nil {
		utils.Error("Failed to parse lab package analysis", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return analysis, nil
}

// Automated Report Generation
func (ai *AIService) GenerateAutomatedReport(ctx context.Context, reportData ReportGenerationData) (*AutomatedReport, error) {
	prompt := ai.buildAutomatedReportPrompt(reportData)
	
	response, err := ai.callOpenAI(ctx, prompt, "You are a medical AI assistant specializing in generating comprehensive medical reports. Create professional, accurate, and well-structured medical reports.")
	if err != nil {
		utils.Error("Failed to call OpenAI for automated report generation", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("automated report generation failed: %w", err)
	}

	report, err := ai.parseAutomatedReport(response)
	if err != nil {
		utils.Error("Failed to parse automated report", utils.LogField{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return report, nil
}

// Helper methods for building prompts
func (ai *AIService) buildMedicalImagePrompt(imageURL, imageType, bodyPart string, patientAge int, patientGender string) string {
	return fmt.Sprintf(`Please analyze this medical image and provide a comprehensive assessment:

Image Details:
- Type: %s
- Body Part: %s
- Patient Age: %d
- Patient Gender: %s
- Image URL: %s

Please provide your analysis in the following JSON format:
{
  "image_type": "%s",
  "body_part": "%s",
  "findings": [
    {
      "finding": "description of finding",
      "location": "anatomical location",
      "severity": "mild/moderate/severe",
      "description": "detailed description",
      "confidence": 0.85
    }
  ],
  "abnormalities": [
    {
      "type": "type of abnormality",
      "location": "location",
      "size": "size if applicable",
      "description": "description",
      "significance": "clinical significance"
    }
  ],
  "measurements": [
    {
      "parameter": "measured parameter",
      "value": "measurement value",
      "unit": "unit of measurement",
      "normal": true/false
    }
  ],
  "recommendations": ["list of recommendations"],
  "urgency_level": "low/medium/high/emergency",
  "confidence": 0.85,
  "requires_review": true/false
}

Important: This analysis is for informational purposes only and requires professional radiologist review.`,
		imageType, bodyPart, patientAge, patientGender, imageURL, imageType, bodyPart)
}

func (ai *AIService) buildAnomalyDetectionPrompt(data []float64, dataType string, patientProfile PatientProfile) string {
	dataStr := fmt.Sprintf("%v", data)
	
	return fmt.Sprintf(`Please analyze this medical data for anomalies:

Data Type: %s
Patient Profile:
- Age: %d
- Gender: %s
- Medical History: %v
- Current Medications: %v

Data Points: %s

Please provide your analysis in the following JSON format:
{
  "anomalies_detected": true/false,
  "anomalies": [
    {
      "data_point": "name of data point",
      "value": actual_value,
      "expected_range": "normal range",
      "severity": "mild/moderate/severe",
      "description": "explanation of anomaly"
    }
  ],
  "overall_risk": "low/medium/high",
  "recommendations": ["list of recommendations"],
  "confidence": 0.85,
  "data_quality": "good/fair/poor"
}

Focus on medically significant anomalies that may require attention.`,
		dataType, patientProfile.Age, patientProfile.Gender, patientProfile.MedicalHistory, patientProfile.Medications, dataStr)
}

func (ai *AIService) buildLabPackagePrompt(packageData LabPackageData) string {
	resultsStr := ""
	for testName, value := range packageData.TestResults {
		refRange := packageData.ReferenceRanges[testName]
		resultsStr += fmt.Sprintf("- %s: %v (Reference: %s)\n", testName, value, refRange)
	}

	return fmt.Sprintf(`Please analyze this comprehensive lab package:

Package Type: %s
Patient Profile:
- Age: %d
- Gender: %s
- Medical History: %v

Test Results:
%s

Please provide your analysis in the following JSON format:
{
  "package_type": "%s",
  "overall_assessment": "comprehensive assessment",
  "key_findings": [
    {
      "category": "category name",
      "finding": "key finding",
      "significance": "clinical significance",
      "impact": "potential impact"
    }
  ],
  "system_analysis": {
    "cardiovascular": {
      "system": "cardiovascular",
      "status": "normal/abnormal",
      "findings": ["list of findings"],
      "concerns": ["list of concerns"]
    }
  },
  "correlations": [
    {
      "tests": ["test1", "test2"],
      "relationship": "relationship description",
      "significance": "clinical significance"
    }
  ],
  "recommendations": ["list of recommendations"],
  "follow_up_tests": ["recommended additional tests"],
  "risk_factors": ["identified risk factors"]
}

Provide a holistic interpretation considering all test interactions and patient context.`,
		packageData.PackageType, packageData.PatientProfile.Age, packageData.PatientProfile.Gender, 
		packageData.PatientProfile.MedicalHistory, resultsStr, packageData.PackageType)
}

func (ai *AIService) buildAutomatedReportPrompt(reportData ReportGenerationData) string {
	testsStr := ""
	for _, test := range reportData.TestResults {
		testsStr += fmt.Sprintf("- %s: %s %s (Ref: %s) - %s\n", 
			test.TestName, test.Value, test.Unit, test.ReferenceRange, test.Status)
	}

	return fmt.Sprintf(`Please generate a comprehensive medical report:

Report Type: %s
Report Purpose: %s
Target Audience: %s

Patient Information:
- Age: %d
- Gender: %s
- Medical History: %v

Test Results:
%s

Clinical Data: %v

Please provide the report in the following JSON format:
{
  "report_id": "generated_report_id",
  "report_type": "%s",
  "patient_summary": "brief patient summary",
  "executive_summary": "executive summary of key findings",
  "detailed_findings": [
    {
      "title": "section title",
      "content": "detailed content",
      "findings": ["key findings"],
      "significance": "clinical significance"
    }
  ],
  "recommendations": ["list of recommendations"],
  "conclusions": "overall conclusions",
  "metadata": {
    "confidence": 0.85,
    "review_required": true/false
  },
  "generated_at": "timestamp"
}

Create a professional, comprehensive report suitable for the target audience.`,
		reportData.ReportType, reportData.ReportPurpose, reportData.TargetAudience,
		reportData.PatientInfo.Age, reportData.PatientInfo.Gender, reportData.PatientInfo.MedicalHistory,
		testsStr, reportData.ClinicalData, reportData.ReportType)
}

// Parsing methods for AI responses
func (ai *AIService) parseMedicalImageAnalysis(response string) (*MedicalImageAnalysis, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var analysis MedicalImageAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal medical image analysis: %w", err)
	}
	
	return &analysis, nil
}

func (ai *AIService) parseAnomalyDetectionResult(response string) (*AnomalyDetectionResult, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var result AnomalyDetectionResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal anomaly detection result: %w", err)
	}
	
	return &result, nil
}

func (ai *AIService) parseLabPackageAnalysis(response string) (*LabPackageAnalysis, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var analysis LabPackageAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lab package analysis: %w", err)
	}
	
	return &analysis, nil
}

func (ai *AIService) parseAutomatedReport(response string) (*AutomatedReport, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var report AutomatedReport
	if err := json.Unmarshal([]byte(jsonStr), &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal automated report: %w", err)
	}
	
	// Set generated timestamp
	report.GeneratedAt = time.Now().Format(time.RFC3339)
	report.ReportID = uuid.New().String()
	
	return &report, nil
}
// Cache-related helper methods

// generateCacheKey creates a cache key for the given operation and input
func (ai *AIService) generateCacheKey(operation string, input interface{}) string {
	if ai.cache == nil {
		return ""
	}
	
	prefix := ai.keyConfig.GetKeyPrefixForOperation(operation)
	hash := ai.cache.GenerateCacheKey(prefix, input)
	return hash
}

// GetCacheStats returns cache performance statistics
func (ai *AIService) GetCacheStats(ctx context.Context) (*cache.CacheStats, error) {
	if ai.cache == nil {
		return nil, fmt.Errorf("cache not available")
	}
	
	stats, err := ai.cache.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	
	// Merge with metrics
	metricsStats := ai.metrics.GetStats()
	stats.Hits = metricsStats.Hits
	stats.Misses = metricsStats.Misses
	stats.HitRate = metricsStats.HitRate
	
	return stats, nil
}

// ClearCache clears all cached AI responses
func (ai *AIService) ClearCache(ctx context.Context) error {
	if ai.cache == nil {
		return fmt.Errorf("cache not available")
	}
	
	return ai.cache.Clear(ctx)
}

// InvalidateCacheKey removes a specific cache entry
func (ai *AIService) InvalidateCacheKey(ctx context.Context, operation string, input interface{}) error {
	if ai.cache == nil {
		return fmt.Errorf("cache not available")
	}
	
	cacheKey := ai.generateCacheKey(operation, input)
	return ai.cache.Delete(ctx, cacheKey)
}
// Generic caching wrapper for AI operations
func (ai *AIService) WithCache(ctx context.Context, operation string, input interface{}, aiCall func() (interface{}, error)) (interface{}, error) {
	// Try cache first if available
	if ai.cache != nil {
		cacheKey := ai.generateCacheKey(operation, input)
		start := time.Now()
		
		if cached, found, err := ai.cache.Get(ctx, cacheKey); err == nil && found {
			ai.metrics.RecordHit(time.Since(start))
			utils.Info("Cache hit", 
				utils.LogField{Key: "operation", Value: operation},
				utils.LogField{Key: "cache_key", Value: cacheKey})
			return cached, nil
		} else if err == nil {
			ai.metrics.RecordMiss(time.Since(start))
		} else {
			ai.metrics.RecordError(time.Since(start))
			utils.Warn("Cache error", 
				utils.LogField{Key: "operation", Value: operation},
				utils.LogField{Key: "error", Value: err.Error()})
		}
	}

	// Cache miss or no cache - call AI function
	result, err := aiCall()
	if err != nil {
		return nil, err
	}

	// Cache the result if cache is available
	if ai.cache != nil {
		cacheKey := ai.generateCacheKey(operation, input)
		ttl := ai.cacheConfig.GetTTLForOperation(operation)
		start := time.Now()
		
		if err := ai.cache.Set(ctx, cacheKey, result, ttl); err != nil {
			ai.metrics.RecordError(time.Since(start))
			utils.Warn("Failed to cache result", 
				utils.LogField{Key: "operation", Value: operation},
				utils.LogField{Key: "error", Value: err.Error()})
		} else {
			ai.metrics.RecordSet(time.Since(start))
			utils.Info("Cached result", 
				utils.LogField{Key: "operation", Value: operation},
				utils.LogField{Key: "cache_key", Value: cacheKey},
				utils.LogField{Key: "ttl", Value: ttl.String()})
		}
	}

	return result, nil
}
