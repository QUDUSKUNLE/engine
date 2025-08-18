package ai

type AIAdaptor struct {
	open_api_key string
	OCR          *TesseractOCR
	Anomaly      *AnomalyDetection
	Report       *AutomatedReport
	Decision     *DecisionSupport
	Analysis     *ImageAnalysis
	Package      *PackageAnalysis
}

func NewAIAdaptor(key string, opts ...Option) *AIAdaptor {
	adaptor := &AIAdaptor{open_api_key: key}

	// apply functional options
	for _, opt := range opts {
		opt(adaptor)
	}
	return adaptor
}

type Option func(*AIAdaptor)

func WithOCR(ocr *TesseractOCR) Option {
	return func(a *AIAdaptor) {
		a.OCR = ocr
	}
}

func WithAnomalyDetector(anomaly *AnomalyDetection) Option {
	return func(a *AIAdaptor) {
		a.Anomaly = anomaly
	}
}

func WithReportGenerator(report *AutomatedReport) Option {
	return func(a *AIAdaptor) {
		a.Report = report
	}
}

func WithDecisionSupport(decision *DecisionSupport) Option {
	return func(a *AIAdaptor) {
		a.Decision = decision
	}
}

func WithImageAnalyzer(analyse *ImageAnalysis) Option {
	return func(a *AIAdaptor) {
		a.Analysis = analyse
	}
}

func WithPackageAnalyzer(pac *PackageAnalysis) Option {
	return func(a *AIAdaptor) {
		a.Package = pac
	}
}
