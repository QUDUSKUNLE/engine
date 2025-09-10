package handlers

import (
	"net/http"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/services"
	"github.com/diagnoxix/core/utils"
	"github.com/labstack/echo/v4"
)

// InterpretLabHandler analyzes lab test results using AI
// @Summary Interpret lab test results
// @Description Provides AI-powered analysis and interpretation of lab test results
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param labTest body domain.LabTest true "Lab test data"
// @Success 200 {object} map[string]interface{} "Lab interpretation results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/interpret_lab [post]
func (h *HTTPHandler) InterpretLabHandler(c echo.Context) error {
	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}

	dto, _ := c.Get(utils.ValidatedBodyDTO).(*domain.LabTest)

	interpretation, err := h.service.AI.InterpretLabResults(c.Request().Context(), *dto)
	if err != nil {
		utils.Error("Lab interpretation failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       interpretation,
		"disclaimer": "This analysis is for informational purposes only and should not replace professional medical consultation.",
	}, c)
}

// AnalyzeSymptomsHandler provides preliminary symptom analysis
// @Summary Analyze patient symptoms
// @Description Provides AI-powered preliminary analysis of patient symptoms
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param symptoms body domain.SymptomAnalysisRequest true "Symptom analysis request"
// @Success 200 {object} map[string]interface{} "Symptom analysis results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/analyze_symptoms [post]
func (h *HTTPHandler) AnalyzeSymptomsHandler(c echo.Context) error {
	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}
	dto, _ := c.Get(utils.ValidatedBodyDTO).(*domain.SymptomAnalysisRequest)

	analysis, err := h.service.AI.AnalyzeSymptoms(c.Request().Context(), dto.Symptoms, dto.Age, dto.Gender)
	if err != nil {
		utils.Error("Symptom analysis failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       analysis,
		"disclaimer": "This is preliminary analysis only. Please consult with a healthcare professional for proper diagnosis and treatment.",
	}, c)
}

// GenerateReportSummaryHandler creates summaries of medical reports
// @Summary Generate medical report summary
// @Description Creates patient-friendly or professional summaries of medical reports
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param report body domain.ReportSummaryRequest true "Report summary request"
// @Success 200 {object} map[string]interface{} "Report summary"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/summarize_report [post]
func (h *HTTPHandler) GenerateReportSummaryHandler(c echo.Context) error {

	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}

	dto, _ := c.Get(utils.ValidatedBodyDTO).(*domain.ReportSummaryRequest)

	summary, err := h.service.AI.GenerateReportSummary(c.Request().Context(), dto.MedicalReport, dto.PatientFriendly)
	if err != nil {
		utils.Error("Report summary generation failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success": true,
		"summary": summary,
		"type": func() string {
			if dto.PatientFriendly {
				return "patient-friendly"
			}
			return "professional"
		}(),
	}, c)
}

// GetAICapabilitiesHandler returns available AI features
// @Summary Get AI capabilities
// @Description Returns list of available AI-powered features
// @Tags AI Medical
// @Produce json
// @Success 200 {object} map[string]interface{} "AI capabilities"
// @Router /v1/ai/capabilities [get]
func (h *HTTPHandler) GetAICapabilitiesHandler(c echo.Context) error {
	capabilities := map[string]interface{}{
		"features": []map[string]interface{}{
			{
				"name":        "Lab Result Interpretation",
				"endpoint":    "/v1/ai/interpret_lab",
				"description": "AI-powered analysis of laboratory test results",
				"method":      "POST",
			},
			{
				"name":        "Symptom Analysis",
				"endpoint":    "/v1/ai/analyze_symptoms",
				"description": "Preliminary analysis of patient symptoms",
				"method":      "POST",
			},
			{
				"name":        "Report Summarization",
				"endpoint":    "/v1/ai/summarize_report",
				"description": "Generate patient-friendly or professional summaries of medical reports",
				"method":      "POST",
			},
			{
				"name":        "Medical Image Analysis",
				"endpoint":    "/v1/ai/analyze_medical_image",
				"description": "AI-powered analysis of medical images (X-rays, CT scans, MRIs)",
				"method":      "POST",
			},
			{
				"name":        "Anomaly Detection",
				"endpoint":    "/v1/ai/detect_anomalies",
				"description": "Detect unusual patterns in medical data",
				"method":      "POST",
			},
			{
				"name":        "Lab Package Analysis",
				"endpoint":    "/v1/ai/analyze_lab_package",
				"description": "Comprehensive analysis of lab test packages",
				"method":      "POST",
			},
			{
				"name":        "Automated Report Generation",
				"endpoint":    "/v1/ai/generate_report",
				"description": "Generate comprehensive medical reports using AI",
				"method":      "POST",
			},
		},
		"disclaimer": "All AI features are for informational purposes only and should not replace professional medical consultation.",
		"version":    "2.0",
	}

	return utils.ResponseMessage(http.StatusOK, capabilities, c)
}

// AnalyzeMedicalImageHandler analyzes medical images using AI
// @Summary Analyze medical images
// @Description Provides AI-powered analysis of medical images (X-rays, CT scans, MRIs, etc.)
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param imageData body domain.MedicalImageAnalysisRequest true "Medical image analysis data"
// @Success 200 {object} map[string]interface{} "Medical image analysis results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/analyze_medical_image [post]
func (h *HTTPHandler) AnalyzeMedicalImageHandler(c echo.Context) error {
	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}
	dto, _ := c.Get(utils.ValidatedBodyDTO).(*domain.MedicalImageAnalysisRequest)
	analysis, err := h.service.AI.AnalyzeMedicalImage(
		c.Request().Context(),
		dto.ImageURL,
		dto.ImageType,
		dto.BodyPart,
		dto.PatientAge,
		dto.PatientGender,
	)
	if err != nil {
		utils.Error("Medical image analysis failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       analysis,
		"disclaimer": "This analysis is for informational purposes only and requires professional radiologist review.",
	}, c)
}

// DetectAnomaliesHandler detects anomalies in medical data
// @Summary Detect anomalies in medical data
// @Description Identifies unusual patterns in medical data that may require attention
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param anomalyData body AnomalyDetectionRequest true "Anomaly detection data"
// @Success 200 {object} map[string]interface{} "Anomaly detection results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/detect_anomalies [post]
func (h *HTTPHandler) DetectAnomaliesHandler(c echo.Context) error {
	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}
	dto, _ := c.Get(utils.ValidatedBodyDTO).(*AnomalyDetectionRequest)

	result, err := h.service.AI.DetectAnomalies(
		c.Request().Context(),
		dto.Data,
		dto.DataType,
		dto.PatientProfile,
	)
	if err != nil {
		utils.Error("Anomaly detection failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       result,
		"disclaimer": "This analysis is for informational purposes only and should not replace professional medical evaluation.",
	}, c)
}

// AnalyzeLabPackageHandler analyzes comprehensive lab packages
// @Summary Analyze lab package
// @Description Provides holistic analysis of comprehensive lab test packages
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param packageData body LabPackageAnalysisRequest true "Lab package analysis data"
// @Success 200 {object} map[string]interface{} "Lab package analysis results"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/analyze_lab_package [post]
func (h *HTTPHandler) AnalyzeLabPackageHandler(c echo.Context) error {
	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}
	dto, _ := c.Get(utils.ValidatedBodyDTO).(*LabPackageAnalysisRequest)

	analysis, err := h.service.AI.AnalyzeLabPackage(
		c.Request().Context(),
		dto.PackageData,
	)
	if err != nil {
		utils.Error("Lab package analysis failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       analysis,
		"disclaimer": "This analysis is for informational purposes only and should not replace professional medical consultation.",
	}, c)
}

// GenerateAutomatedReportHandler generates comprehensive medical reports
// @Summary Generate automated medical report
// @Description Creates comprehensive, professional medical reports using AI
// @Tags AI Medical
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param reportData body AutomatedReportRequest true "Report generation data"
// @Success 200 {object} map[string]interface{} "Generated medical report"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/ai/generate_report [post]
func (h *HTTPHandler) GenerateAutomatedReportHandler(c echo.Context) error {
	_, err := services.PrivateMiddlewareContext(c, []db.UserEnum{db.UserEnumPATIENT})
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, c)
	}
	dto, _ := c.Get(utils.ValidatedBodyDTO).(*AutomatedReportRequest)

	report, err := h.service.AI.GenerateAutomatedReport(
		c.Request().Context(),
		dto.ReportData,
	)
	if err != nil {
		utils.Error("Automated report generation failed", utils.LogField{Key: "error", Value: err.Error()})
		return utils.ErrorResponse(http.StatusInternalServerError, err, c)
	}

	return utils.ResponseMessage(http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       report,
		"disclaimer": "This report is AI-generated and should be reviewed by qualified medical professionals.",
	}, c)
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
