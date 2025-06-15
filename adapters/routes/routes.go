package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/handlers"
	"github.com/medicue/adapters/middlewares"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

// routeConfig defines the structure for route configuration
type routeConfig struct {
	method      string
	path        string
	handler     echo.HandlerFunc
	factory     func() interface{}
	description string // Added for better logging and documentation
}

// routeGroup helps organize routes by their functional area
type routeGroup struct {
	name   string
	routes []routeConfig
}

func RoutesAdaptor(public *echo.Group, handler *handlers.HTTPHandler) *echo.Group {
	// Organize routes by functional groups
	routeGroups := []routeGroup{
		{
			name: "Authentication",
			routes: []routeConfig{
				{
					method:      http.MethodPost,
					path:        "/register",
					handler:     handler.Register,
					factory:     func() interface{} { return &domain.UserRegisterDTO{} },
					description: "User registration",
				},
				{
					method:      http.MethodPost,
					path:        "/login",
					handler:     handler.SignIn,
					factory:     func() interface{} { return &domain.UserSignInDTO{} },
					description: "User login",
				},
				{
					method:      http.MethodPost,
					path:        "/request_password_reset",
					handler:     handler.RequestPasswordReset,
					factory:     func() interface{} { return &domain.RequestPasswordResetDTO{} },
					description: "Request password reset",
				},
				{
					method:      http.MethodPost,
					path:        "/reset_password",
					handler:     handler.ResetPassword,
					factory:     func() interface{} { return &domain.ResetPasswordDTO{} },
					description: "Reset password",
				},
				{
					method:      http.MethodPut,
					path:        "/account/password",
					handler:     handler.ChangePassword,
					factory:     func() interface{} { return &domain.ChangePasswordDTO{} },
					description: "Change current password",
				},
				{
					method:      http.MethodGet,
					path:        "/account/profile",
					handler:     handler.GetProfile,
					factory:     func() interface{} { return &domain.GetUserProfileParamDTO{} },
					description: "Get user profile",
				},
				{
					method:      http.MethodPut,
					path:        "/account/profile",
					handler:     handler.UpdateProfile,
					factory:     func() interface{} { return &domain.UpdateUserProfileDTO{} },
					description: "Update user profile",
				},
				{
					method:      http.MethodDelete,
					path:        "/account",
					handler:     handler.DeactivateAccount,
					factory:     func() interface{} { return &domain.DeactivateAccountDTO{} },
					description: "Deactivate user account",
				},
				{
					method:      http.MethodPost,
					path:        "/verify_email",
					handler:     handler.VerifyEmail,
					factory:     func() interface{} { return &domain.EmailVerificationDTO{} },
					description: "Verify user email address",
				},
				{
					method:      http.MethodPost,
					path:        "/resend_verification",
					handler:     handler.ResendVerification,
					factory:     func() interface{} { return &domain.ResendVerificationDTO{} },
					description: "Resend email verification",
				},
				{
					method:      http.MethodPost,
					path:        "/auth/google",
					handler:     handler.GoogleLogin,
					factory:     func() interface{} { return &domain.GoogleAuthDTO{} },
					description: "Google OAuth login",
				},
			},
		},
		{
			name: "Diagnostic Centre",
			routes: []routeConfig{
				{
					method:      http.MethodGet,
					path:        "/diagnostic_centres/:diagnostic_centre_id",
					handler:     handler.GetDiagnosticCentre,
					factory:     func() interface{} { return &domain.GetDiagnosticParamDTO{} },
					description: "Get diagnostic centre details",
				},
				{
					method:      http.MethodPost,
					path:        "/diagnostic_centre_manager",
					handler:     handler.CreateDiagnosticCentreManager,
					factory:     func() interface{} { return &domain.DiagnosticCentreManagerRegisterDTO{} },
					description: "Create diagnostic centre manager",
				},
				{
					method:      http.MethodPost,
					path:        "/diagnostic_centres",
					handler:     handler.CreateDiagnostic,
					factory:     func() interface{} { return &domain.CreateDiagnosticDTO{} },
					description: "Create diagnostic centre",
				},
				{
					method:      http.MethodGet,
					path:        "/diagnostic_centres",
					handler:     handler.SearchDiagnosticCentre,
					factory:     func() interface{} { return &domain.SearchDiagnosticCentreQueryDTO{} },
					description: "Search diagnostic centres",
				},
				{
					method:      http.MethodPut,
					path:        "/diagnostic_centres/:diagnostic_centre_id",
					handler:     handler.UpdateDiagnosticCentre,
					factory:     func() interface{} { return &domain.UpdateDiagnosticBodyDTO{} },
					description: "Update diagnostic centre",
				},
				{
					method:  http.MethodGet,
					path:    "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules",
					handler: handler.GetDiagnosticSchedulesByCentre,
					factory: func() interface{} {
						return &domain.GetDiagnosticSchedulesByCentreParamDTO{}
					},
					description: "Get a diagnostic centre schedules",
				},
				{
					method:  http.MethodGet,
					path:    "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules/:schedule_id",
					handler: handler.GetDiagnosticScheduleByCentre,
					factory: func() interface{} {
						return &domain.GetDiagnosticScheduleByCentreParamDTO{}
					},
					description: "Get a diagnostic centre schedule",
				},
				{
					method:      http.MethodPut,
					path:        "/diagnostic_centres/:diagnostic_centre_id/diagnostic_schedules/:schedule_id",
					handler:     handler.UpdateDiagnosticScheduleByCentre,
					factory:     func() interface{} { return &domain.UpdateDiagnosticScheduleByCentreDTO{} },
					description: "Update a diagnostic centre schedule",
				},
			},
		},
		{
			name: "Schedule",
			routes: []routeConfig{
				{
					method:      http.MethodPost,
					path:        "/diagnostic_schedules",
					handler:     handler.CreateSchedule,
					factory:     func() interface{} { return &domain.CreateScheduleDTO{} },
					description: "Create schedule",
				},
				{
					method:      http.MethodGet,
					path:        "/diagnostic_schedules/:schedule_id",
					handler:     handler.GetSchedule,
					factory:     func() interface{} { return &domain.GetDiagnosticScheduleParamDTO{} },
					description: "Get schedule details",
				},
				{
					method:      http.MethodPut,
					path:        "/diagnostic_schedules/:schedule_id",
					handler:     handler.UpdateSchedule,
					factory:     func() interface{} { return &domain.UpdateScheduleDTO{} },
					description: "Update a schedule",
				},
				{
					method:      http.MethodGet,
					path:        "/diagnostic_schedules",
					handler:     handler.GetSchedules,
					factory:     func() interface{} { return &domain.GetDiagnosticSchedulesQueryDTO{} },
					description: "Get all schedules",
				},
			},
		},
		{
			name: "Medical Records",
			routes: []routeConfig{
				{
					method:      http.MethodPost,
					path:        "/medical_records",
					handler:     handler.CreateMedicalRecord,
					factory:     func() interface{} { return &domain.CreateMedicalRecordDTO{} },
					description: "Create medical record",
				},
				{
					method:      http.MethodGet,
					path:        "/medical_records/:record_id",
					handler:     handler.GetMedicalRecord,
					factory:     func() interface{} { return &domain.GetMedicalRecordParamsDTO{} },
					description: "Get medical record details",
				},
				{
					method:      http.MethodGet,
					path:        "/medical_records",
					handler:     handler.GetMedicalRecords,
					factory:     func() interface{} { return &domain.GetMedicalRecordsParamQueryDTO{} },
					description: "Get all medical records",
				},
				{
					method:      http.MethodGet,
					path:        "/medical_records/:record_id/diagnostic_centre/:diagnostic_centre_id",
					handler:     handler.GetUploaderMedicalRecord,
					factory:     func() interface{} { return &domain.GetUploaderMedicalRecordParamsDTO{} },
					description: "Get uploader medical record",
				},
				{
					method:      http.MethodGet,
					path:        "/medical_records/diagnostic_centre/:diagnostic_centre_id",
					handler:     handler.GetUploaderMedicalRecords,
					factory:     func() interface{} { return &domain.GetUploaderMedicalRecordsParamQueryDTO{} },
					description: "Get uploader medical records by centre",
				},
				{
					method:      http.MethodPut,
					path:        "/medical_records",
					handler:     handler.UpdateMedicalRecord,
					factory:     func() interface{} { return &domain.UpdateMedicalRecordDTO{} },
					description: "Update a medical record",
				},
			},
		},
		{
			name: "Appointments",
			routes: []routeConfig{
				{
					method:      http.MethodPost,
					path:        "/appointments",
					handler:     handler.CreateAppointment,
					factory:     func() interface{} { return &domain.CreateAppointmentDTO{} },
					description: "Book a new appointment",
				},
				{
					method:      http.MethodGet,
					path:        "/appointments/:appointment_id",
					handler:     handler.GetAppointment,
					factory:     func() interface{} { return &domain.GetAppointmentDTO{} },
					description: "Get appointment details",
				},
				{
					method:      http.MethodGet,
					path:        "/appointments",
					handler:     handler.ListAppointments,
					factory:     func() interface{} { return &domain.ListAppointmentsDTO{} },
					description: "List appointments with filters",
				},
				{
					method:      http.MethodPut,
					path:        "/appointments/:appointment_id",
					handler:     handler.UpdateAppointment,
					factory:     func() interface{} { return &domain.UpdateAppointmentDTO{} },
					description: "Update appointment details",
				},
				{
					method:      http.MethodPut,
					path:        "/appointments/:appointment_id/cancel",
					handler:     handler.CancelAppointment,
					factory:     func() interface{} { return &domain.CancelAppointmentDTO{} },
					description: "Cancel an appointment",
				},
				{
					method:      http.MethodPut,
					path:        "/appointments/:appointment_id/reschedule",
					handler:     handler.RescheduleAppointment,
					factory:     func() interface{} { return &domain.RescheduleAppointmentDTO{} },
					description: "Reschedule an appointment",
				},
			},
		},
		{
			name: "Payments",
			routes: []routeConfig{
				{
					method:      http.MethodPost,
					path:        "/payments",
					handler:     handler.CreatePayment,
					factory:     func() interface{} { return &domain.CreatePaymentDTO{} },
					description: "Process a new payment",
				},
				{
					method:      http.MethodGet,
					path:        "/payments/:payment_id",
					handler:     handler.GetPayment,
					factory:     func() interface{} { return &domain.GetPaymentDTO{} },
					description: "Get payment details",
				},
				{
					method:      http.MethodGet,
					path:        "/payments",
					handler:     handler.ListPayments,
					factory:     func() interface{} { return &domain.ListPaymentsDTO{} },
					description: "List payments with filters",
				},
				{
					method:      http.MethodPost,
					path:        "/payments/:payment_id/refund",
					handler:     handler.RefundPayment,
					factory:     func() interface{} { return &domain.RefundPaymentDTO{} },
					description: "Process payment refund",
				},
				{
					method:      http.MethodPost,
					path:        "/payments/webhook",
					handler:     handler.PaymentWebhook,
					factory:     func() interface{} { return &domain.PaymentWebhookDTO{} },
					description: "Handle payment provider webhook",
				},
			},
		},
		{
			name: "Availability",
			routes: []routeConfig{
				{
					method:      http.MethodPost,
					path:        "/availability",
					handler:     handler.CreateAvailability,
					factory:     func() interface{} { return &domain.CreateAvailabilityDTO{} },
					description: "Create diagnostic centre availability",
				},
				{
					method:      http.MethodGet,
					path:        "/availability/:diagnostic_centre_id",
					handler:     handler.GetAvailability,
					factory:     func() interface{} { return &domain.GetAvailabilityDTO{} },
					description: "Get diagnostic centre availability",
				},
				{
					method:      http.MethodPut,
					path:        "/availability/:diagnostic_centre_id/:day_of_week",
					handler:     handler.UpdateAvailability,
					factory:     func() interface{} { return &domain.UpdateAvailabilityDTO{} },
					description: "Update diagnostic centre availability",
				},
				{
					method:      http.MethodDelete,
					path:        "/availability/:diagnostic_centre_id/:day_of_week",
					handler:     handler.DeleteAvailability,
					factory:     func() interface{} { return nil },
					description: "Delete diagnostic centre availability",
				},
			},
		},
	}

	// Register all routes with logging
	for _, group := range routeGroups {
		utils.Info("Registering route group",
			utils.LogField{Key: "group", Value: group.name})

		for _, route := range group.routes {
			registerRoute(public, route)
		}
	}

	return public
}

// registerRoute handles the registration of a single route with proper logging
func registerRoute(group *echo.Group, route routeConfig) {
	utils.Info("Registering route",
		utils.LogField{Key: "method", Value: route.method},
		utils.LogField{Key: "path", Value: "/v1" + route.path},
		utils.LogField{Key: "description", Value: route.description})

	switch route.method {
	case http.MethodPost:
		group.POST(
			route.path,
			route.handler,
			middlewares.BodyValidationInterceptorFor(route.factory),
		)
	case http.MethodGet:
		group.GET(
			route.path,
			route.handler,
			middlewares.BodyValidationInterceptorFor(route.factory),
		)
	case http.MethodPut:
		group.PUT(
			route.path,
			route.handler,
			middlewares.BodyValidationInterceptorFor(route.factory),
		)
	case http.MethodDelete:
		group.DELETE(
			route.path,
			route.handler,
			middlewares.BodyValidationInterceptorFor(route.factory),
		)
	default:
		utils.Error("Unsupported HTTP method",
			utils.LogField{Key: "method", Value: route.method},
			utils.LogField{Key: "path", Value: route.path})
	}
}
