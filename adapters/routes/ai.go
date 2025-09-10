package routes

import (
	"net/http"

	"github.com/diagnoxix/adapters/handlers"
	"github.com/diagnoxix/core/domain"
	"github.com/labstack/echo/v4"
)

// AIRoutes registers all AI-related routes
func AIRoutes(group *echo.Group, handler *handlers.HTTPHandler) {
	aiGroup := []routeConfig{
		{
			method:      http.MethodPost,
			path:        "/ai/interpret-lab",
			handler:     handler.InterpretLabHandler,
			factory:     func() interface{} { return &domain.LabTest{} },
			description: "Interpret lab test results using AI",
		},
		{
			method:      http.MethodPost,
			path:        "/ai/analyze-symptoms",
			handler:     handler.AnalyzeSymptomsHandler,
			factory:     func() interface{} { return &domain.SymptomAnalysisRequest{} },
			description: "Analyze patient symptoms using AI",
		},
		{
			method:      http.MethodPost,
			path:        "/ai/summarize-report",
			handler:     handler.GenerateReportSummaryHandler,
			factory:     func() interface{} { return &domain.ReportSummaryRequest{} },
			description: "Generate medical report summary using AI",
		},
		{
			method:      http.MethodGet,
			path:        "/ai/capabilities",
			handler:     handler.GetAICapabilitiesHandler,
			factory:     func() interface{} { return &domain.CapabilitiesDTO{} },
			description: "Get available AI capabilities",
		},
		{
			method:      http.MethodPost,
			path:        "/ai/analyze-medical-image",
			handler:     handler.AnalyzeMedicalImageHandler,
			factory:     func() interface{} { return &domain.MedicalImageAnalysisRequest{} },
			description: "Analyze medical images using AI",
		},
		{
			method:      http.MethodPost,
			path:        "/ai/detect-anomalies",
			handler:     handler.DetectAnomaliesHandler,
			factory:     func() interface{} { return &handlers.AnomalyDetectionRequest{} },
			description: "Detect anomalies in medical data",
		},
		{
			method:      http.MethodPost,
			path:        "/ai/analyze-lab-package",
			handler:     handler.AnalyzeLabPackageHandler,
			factory:     func() interface{} { return &handlers.LabPackageAnalysisRequest{} },
			description: "Analyze comprehensive lab packages",
		},
		{
			method:      http.MethodPost,
			path:        "/ai/generate-report",
			handler:     handler.GenerateAutomatedReportHandler,
			factory:     func() interface{} { return &handlers.AutomatedReportRequest{} },
			description: "Generate automated medical reports",
		},
	}

	registerRoutes(group, aiGroup)
}
