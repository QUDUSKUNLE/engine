package services

import (
	"strconv"

	"github.com/medivue/adapters/config"
	"github.com/medivue/adapters/ex"
	"github.com/medivue/adapters/ex/paystack"
	"github.com/medivue/core/ports"
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
	Config              config.Config
	// Payment Gateway
	paymentService ports.PaymentProviderService
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
	conn config.Config,
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
		notificationService: ex.NewNotificationAdapter(&ex.GmailConfig{
			Host:     conn.EMAIL_HOST,
			Port:     func() int { p, _ := strconv.Atoi(conn.EMAIL_PORT); return p }(),
			Username: conn.GMAIL_USERNAME,
			Password: conn.GMAIL_APP_PASSWORD,
			From:     conn.EMAIL_FROM_ADDRESS,
		}),
		paymentService: paystack.NewPaystackAdapter(&paystack.PaystackConfig{
			SecretKey: conn.PAYSTACK_SECRET_KEY,
			BaseURL:   conn.PAYSTACK_BASE_URL,
		}),
		Config: conn,
	}
}
