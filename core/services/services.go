package services

import (
	"strconv"

	"github.com/diagnoxix/adapters/config"
	"github.com/diagnoxix/adapters/ex"
	"github.com/diagnoxix/adapters/ex/ai"
	"github.com/diagnoxix/adapters/ex/paystack"
	"github.com/diagnoxix/core/ports"
)

type ServicesHandler struct {
	userPort         ports.UserRepository
	schedulePort     ports.ScheduleRepository
	diagnosticPort   ports.DiagnosticRepository
	recordPort       ports.RecordRepository
	availabilityPort ports.AvailabilityRepository
	paymentPort      ports.PaymentRepository
	appointmentPort  ports.AppointmentRepository
	testPricePort    ports.TestPriceRepository
	filePort         ports.FileService
	notificationPort ports.NotificationService
	notificationRepo ports.NotificationRepository
	Config           config.EnvConfiguration
	// Payment Gateway
	paymentService ports.PaymentProviderService

	aiPort *ai.AIAdaptor
}

func ServicesAdapter(
	useRepo ports.UserRepository,
	schedulePort ports.ScheduleRepository,
	diagnosticCentreRepo ports.DiagnosticRepository,
	availabilityPort ports.AvailabilityRepository,
	record ports.RecordRepository,
	paymentPort ports.PaymentRepository,
	appointmentPort ports.AppointmentRepository,
	testPriceRepo ports.TestPriceRepository,
	notificationRepo ports.NotificationRepository,
	conn config.EnvConfiguration,
) *ServicesHandler {
	return &ServicesHandler{
		userPort:         useRepo,
		schedulePort:     schedulePort,
		diagnosticPort:   diagnosticCentreRepo,
		availabilityPort: availabilityPort,
		paymentPort:      paymentPort,
		appointmentPort:  appointmentPort,
		recordPort:       record,
		testPricePort:    testPriceRepo,
		notificationPort: ex.NewNotificationAdapter(&ex.EmailConfig{
			Host:      conn.EMAIL_HOST,
			Port:      func() int { p, _ := strconv.Atoi(conn.EMAIL_PORT); return p }(),
			Username:  conn.EMAIL_USERNAME,
			Password:  conn.EMAIL_APP_PASSWORD,
			From:      conn.EMAIL_FROM_ADDRESS,
			EmailType: ex.EmailType(conn.EMAIL_TYPE),
		}),
		notificationRepo: notificationRepo,
		paymentService: paystack.NewPaystackAdapter(&paystack.PaystackConfig{
			SecretKey: conn.PAYSTACK_SECRET_KEY,
			BaseURL:   conn.PAYSTACK_BASE_URL,
		}),
		Config: conn,
		aiPort: ai.NewAIAdaptor(
			conn.OPEN_API_KEY,
			ai.WithOCR(ai.NewTesseractOCR()),
			ai.WithAnomalyDetector(ai.NewAnomalyDetection()),
			ai.WithReportGenerator(ai.NewAutomatedReport()),
			ai.WithDecisionSupport(ai.NewDecisionSupport()),
			ai.WithImageAnalyzer(ai.NewImageAnalysis()),
			ai.WithPackageAnalyzer(ai.NewPackageAnalysis()),
		),
		filePort: ex.NewLocalFileService(),
	}
}
