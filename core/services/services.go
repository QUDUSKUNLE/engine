package services

import (
	"strconv"
	"time"

	"github.com/diagnoxix/adapters/config"
	"github.com/diagnoxix/adapters/ex"
	"github.com/diagnoxix/adapters/ex/ai"
	"github.com/diagnoxix/adapters/ex/paystack"
	"github.com/diagnoxix/core/ports"
	"github.com/diagnoxix/core/services/cache"
	"github.com/diagnoxix/core/services/websocket"
	"github.com/diagnoxix/core/utils"
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
	// AI Services
	aiPort *ai.AIAdaptor
	AI     *AIService
	// WebSocket Manager
	WebSocketManager *websocket.WebSocketManager
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
		AI:               initializeAIService(conn),
		WebSocketManager: websocket.NewWebSocketManager(),
		filePort:         ex.NewLocalFileService(),
	}
}
// initializeAIService creates an AI service with caching if Redis is available
func initializeAIService(conn config.EnvConfiguration) *AIService {
	// Try to initialize cache if Redis URL is provided
	if conn.REDIS_URL != "" {
		cacheConfig := cache.CacheConfig{
			RedisURL:     conn.REDIS_URL,
			DefaultTTL:   1 * time.Hour,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
		}

		aiCache, err := cache.NewAICache(cacheConfig)
		if err != nil {
			utils.Warn("Failed to initialize AI cache, proceeding without caching", 
				utils.LogField{Key: "error", Value: err.Error()})
			return NewAIService(conn.OPEN_API_KEY)
		}

		utils.Info("AI cache initialized successfully")
		return NewAIServiceWithCache(conn.OPEN_API_KEY, aiCache)
	}

	utils.Info("No Redis URL provided, AI service running without cache")
	return NewAIService(conn.OPEN_API_KEY)
}
