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
	UserRepo            ports.UserRepository
	ScheduleRepo        ports.ScheduleRepository
	DiagnosticRepo      ports.DiagnosticRepository
	RecordRepo          ports.RecordRepository
	AvailabilityRepo    ports.AvailabilityRepository
	PaymentRepo         ports.PaymentRepository
	AppointmentRepo     ports.AppointmentRepository
	TestPriceRepo       ports.TestPriceRepository
	FileRepo            ports.FileService
	notificationService ports.NotificationService
	Config              config.EnvConfiguration
	// Payment Gateway
	paymentService ports.PaymentProviderService
	aiService      ports.AIService
}

func ServicesAdapter(
	useRepo ports.UserRepository,
	scheduleRepo ports.ScheduleRepository,
	diagnosticCentreRepo ports.DiagnosticRepository,
	availabilityRepo ports.AvailabilityRepository,
	record ports.RecordRepository,
	paymentPort ports.PaymentRepository,
	appointmentPort ports.AppointmentRepository,
	testPriceRepo ports.TestPriceRepository,
	conn config.EnvConfiguration,
) *ServicesHandler {
	return &ServicesHandler{
		UserRepo:         useRepo,
		ScheduleRepo:     scheduleRepo,
		DiagnosticRepo:   diagnosticCentreRepo,
		AvailabilityRepo: availabilityRepo,
		PaymentRepo:      paymentPort,
		AppointmentRepo:  appointmentPort,
		RecordRepo:       record,
		TestPriceRepo:    testPriceRepo,
		notificationService: ex.NewNotificationAdapter(&ex.EmailConfig{
			Host:      conn.EMAIL_HOST,
			Port:      func() int { p, _ := strconv.Atoi(conn.EMAIL_PORT); return p }(),
			Username:  conn.EMAIL_USERNAME,
			Password:  conn.EMAIL_APP_PASSWORD,
			From:      conn.EMAIL_FROM_ADDRESS,
			EmailType: ex.EmailType(conn.EMAIL_TYPE),
		}),
		paymentService: paystack.NewPaystackAdapter(&paystack.PaystackConfig{
			SecretKey: conn.PAYSTACK_SECRET_KEY,
			BaseURL:   conn.PAYSTACK_BASE_URL,
		}),
		Config:    conn,
		aiService: ai.NewAIAdaptor(conn.OPEN_API_KEY),
		FileRepo: ex.NewLocalFileService(),
	}
}
