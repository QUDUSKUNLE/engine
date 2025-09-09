package handlers

import (
	"net/http"

	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/services"
	"github.com/diagnoxix/core/utils"
	"github.com/labstack/echo/v4"
)

// InterpretLabHandler analyzes lab test results using AI
// @Summary Interpret lab test results
// @Description Provides AI-powered analysis and interpretation of lab test results
// @Tags AI
// @Accept json
// @Produce json
// @Param labTest body domain.LabTest true "Lab test data"
// @Success 200 {object} map[string]interface{} "Lab interpretation results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/interpret-lab [post]
func (h *HTTPHandler) InterpretLabHandler(c echo.Context) error {
	var labTest domain.LabTest
	if err := c.Bind(&labTest); err != nil {
		utils.Error("Failed to bind lab test data", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&labTest); err != nil {
		utils.Error("Lab test validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	interpretation, err := h.service.AI.InterpretLabResults(c.Request().Context(), labTest)
	if err != nil {
		utils.Error("Lab interpretation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to interpret lab results",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    interpretation,
		"disclaimer": "This analysis is for informational purposes only and should not replace professional medical consultation.",
	})
}

// AnalyzeSymptomsHandler provides preliminary symptom analysis
// @Summary Analyze patient symptoms
// @Description Provides AI-powered preliminary analysis of patient symptoms
// @Tags AI
// @Accept json
// @Produce json
// @Param symptoms body domain.SymptomAnalysisRequest true "Symptom analysis request"
// @Success 200 {object} map[string]interface{} "Symptom analysis results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/analyze-symptoms [post]
func (h *HTTPHandler) AnalyzeSymptomsHandler(c echo.Context) error {
	var req domain.SymptomAnalysisRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind symptom analysis request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Symptom analysis validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	analysis, err := h.service.AI.AnalyzeSymptoms(c.Request().Context(), req.Symptoms, req.Age, req.Gender)
	if err != nil {
		utils.Error("Symptom analysis failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to analyze symptoms",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analysis,
		"disclaimer": "This is preliminary analysis only. Please consult with a healthcare professional for proper diagnosis and treatment.",
	})
}

// GenerateReportSummaryHandler creates summaries of medical reports
// @Summary Generate medical report summary
// @Description Creates patient-friendly or professional summaries of medical reports
// @Tags AI
// @Accept json
// @Produce json
// @Param report body domain.ReportSummaryRequest true "Report summary request"
// @Success 200 {object} map[string]interface{} "Report summary"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/summarize-report [post]
func (h *HTTPHandler) GenerateReportSummaryHandler(c echo.Context) error {
	var req domain.ReportSummaryRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind report summary request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Report summary validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	summary, err := h.service.AI.GenerateReportSummary(c.Request().Context(), req.MedicalReport, req.PatientFriendly)
	if err != nil {
		utils.Error("Report summary generation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to generate report summary",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"summary": summary,
		"type": func() string {
			if req.PatientFriendly {
				return "patient-friendly"
			}
			return "professional"
		}(),
	})
}

// GetAICapabilitiesHandler returns available AI features
// @Summary Get AI capabilities
// @Description Returns list of available AI-powered features
// @Tags AI
// @Produce json
// @Success 200 {object} map[string]interface{} "AI capabilities"
// @Router /v1/ai/capabilities [get]
func (h *HTTPHandler) GetAICapabilitiesHandler(c echo.Context) error {
	capabilities := map[string]interface{}{
		"features": []map[string]interface{}{
			{
				"name": "Lab Result Interpretation",
				"endpoint": "/v1/ai/interpret-lab",
				"description": "AI-powered analysis of laboratory test results",
				"method": "POST",
			},
			{
				"name": "Symptom Analysis",
				"endpoint": "/v1/ai/analyze-symptoms",
				"description": "Preliminary analysis of patient symptoms",
				"method": "POST",
			},
			{
				"name": "Report Summarization",
				"endpoint": "/v1/ai/summarize-report",
				"description": "Generate patient-friendly or professional summaries of medical reports",
				"method": "POST",
			},
			{
				"name": "Medical Image Analysis",
				"endpoint": "/v1/ai/analyze-medical-image",
				"description": "AI-powered analysis of medical images (X-rays, CT scans, MRIs)",
				"method": "POST",
			},
			{
				"name": "Anomaly Detection",
				"endpoint": "/v1/ai/detect-anomalies",
				"description": "Detect unusual patterns in medical data",
				"method": "POST",
			},
			{
				"name": "Lab Package Analysis",
				"endpoint": "/v1/ai/analyze-lab-package",
				"description": "Comprehensive analysis of lab test packages",
				"method": "POST",
			},
			{
				"name": "Automated Report Generation",
				"endpoint": "/v1/ai/generate-report",
				"description": "Generate comprehensive medical reports using AI",
				"method": "POST",
			},
		},
		"disclaimer": "All AI features are for informational purposes only and should not replace professional medical consultation.",
		"version": "2.0",
	}

	return c.JSON(http.StatusOK, capabilities)
}

// AnalyzeMedicalImageHandler analyzes medical images using AI
// @Summary Analyze medical images
// @Description Provides AI-powered analysis of medical images (X-rays, CT scans, MRIs, etc.)
// @Tags AI
// @Accept json
// @Produce json
// @Param imageData body domain.MedicalImageAnalysisRequest true "Medical image analysis data"
// @Success 200 {object} map[string]interface{} "Medical image analysis results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/analyze-medical-image [post]
func (h *HTTPHandler) AnalyzeMedicalImageHandler(c echo.Context) error {
	var req domain.MedicalImageAnalysisRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind medical image analysis request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Medical image analysis validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	analysis, err := h.service.AI.AnalyzeMedicalImage(
		c.Request().Context(),
		req.ImageURL,
		req.ImageType,
		req.BodyPart,
		req.PatientAge,
		req.PatientGender,
	)
	if err != nil {
		utils.Error("Medical image analysis failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to analyze medical image",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analysis,
		"disclaimer": "This analysis is for informational purposes only and requires professional radiologist review.",
	})
}

// DetectAnomaliesHandler detects anomalies in medical data
// @Summary Detect anomalies in medical data
// @Description Identifies unusual patterns in medical data that may require attention
// @Tags AI
// @Accept json
// @Produce json
// @Param anomalyData body AnomalyDetectionRequest true "Anomaly detection data"
// @Success 200 {object} map[string]interface{} "Anomaly detection results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/detect-anomalies [post]
func (h *HTTPHandler) DetectAnomaliesHandler(c echo.Context) error {
	var req AnomalyDetectionRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind anomaly detection request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Anomaly detection validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	result, err := h.service.AI.DetectAnomalies(
		c.Request().Context(),
		req.Data,
		req.DataType,
		req.PatientProfile,
	)
	if err != nil {
		utils.Error("Anomaly detection failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to detect anomalies",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
		"disclaimer": "This analysis is for informational purposes only and should not replace professional medical evaluation.",
	})
}

// AnalyzeLabPackageHandler analyzes comprehensive lab packages
// @Summary Analyze lab package
// @Description Provides holistic analysis of comprehensive lab test packages
// @Tags AI
// @Accept json
// @Produce json
// @Param packageData body LabPackageAnalysisRequest true "Lab package analysis data"
// @Success 200 {object} map[string]interface{} "Lab package analysis results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/analyze-lab-package [post]
func (h *HTTPHandler) AnalyzeLabPackageHandler(c echo.Context) error {
	var req LabPackageAnalysisRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind lab package analysis request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Lab package analysis validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	analysis, err := h.service.AI.AnalyzeLabPackage(
		c.Request().Context(),
		req.PackageData,
	)
	if err != nil {
		utils.Error("Lab package analysis failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to analyze lab package",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analysis,
		"disclaimer": "This analysis is for informational purposes only and should not replace professional medical consultation.",
	})
}

// GenerateAutomatedReportHandler generates comprehensive medical reports
// @Summary Generate automated medical report
// @Description Creates comprehensive, professional medical reports using AI
// @Tags AI
// @Accept json
// @Produce json
// @Param reportData body AutomatedReportRequest true "Report generation data"
// @Success 200 {object} map[string]interface{} "Generated medical report"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/generate-report [post]
func (h *HTTPHandler) GenerateAutomatedReportHandler(c echo.Context) error {
	var req AutomatedReportRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("Failed to bind automated report request", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		utils.Error("Automated report validation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	report, err := h.service.AI.GenerateAutomatedReport(
		c.Request().Context(),
		req.ReportData,
	)
	if err != nil {
		utils.Error("Automated report generation failed", utils.LogField{Key: "error", Value: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to generate automated report",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
		"disclaimer": "This report is AI-generated and should be reviewed by qualified medical professionals.",
	})
}


	// AnomalyDetectionRequest represents an anomaly detection request
	type AnomalyDetectionRequest struct {
		Data           []float64               `json:"data" validate:"required,min=1"`
		DataType       string                  `json:"data_type" validate:"required"`
		PatientProfile services.PatientProfile `json:"patient_profile" validate:"required"`
	}

	// LabPackageAnalysisRequest represents a lab package analysis request
	type LabPackageAnalysisRequest struct {
		PackageData services.LabPackageData `json:"package_data" validate:"required"`
	}

	// AutomatedReportRequest represents an automated report generation request
	type AutomatedReportRequest struct {
		ReportData services.ReportGenerationData `json:"report_data" validate:"required"`
	}
