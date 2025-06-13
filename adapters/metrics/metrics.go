package metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// Database metrics
	DbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation", "table"},
	)

	// Diagnostic centre metrics
	DiagnosticSearches = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "diagnostic_searches_total",
			Help: "Total number of diagnostic centre searches",
		},
		[]string{"has_filters"},
	)

	DiagnosticCentresFound = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "diagnostic_centres_found",
			Help:    "Number of diagnostic centres found in searches",
			Buckets: []float64{0, 1, 2, 5, 10, 20, 50},
		},
		[]string{"has_filters"},
	)

	// Schedule metrics
	SchedulesCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "schedules_created_total",
		Help: "Total number of schedules created",
	})

	ScheduleStatusUpdates = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "schedule_status_updates_total",
			Help: "Total number of schedule status updates",
		},
		[]string{"from_status", "to_status"},
	)

	// Record metrics
	RecordsUploaded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "records_uploaded_total",
		Help: "Total number of medical records uploaded",
	})

	RecordUploadErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "record_upload_errors_total",
			Help: "Total number of errors during record uploads",
		},
		[]string{"error_type"},
	)

	// User metrics
	UserRegistrations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
		[]string{"user_type"},
	)

	LoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"success"},
	)
)

// ObserveHTTPRequest records metrics about an HTTP request
func ObserveHTTPRequest(c echo.Context, duration float64) {
	status := c.Response().Status
	path := c.Path()
	method := c.Request().Method

	HttpRequestsTotal.WithLabelValues(method, path, string(rune(status))).Inc()
	HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// ObserveDBQuery records metrics about a database query
func ObserveDBQuery(operation, table string, duration float64) {
	DbQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// RecordDiagnosticSearch records metrics about a diagnostic centre search
func RecordDiagnosticSearch(hasFilters bool, resultsCount int) {
	label := "false"
	if hasFilters {
		label = "true"
	}
	DiagnosticSearches.WithLabelValues(label).Inc()
	DiagnosticCentresFound.WithLabelValues(label).Observe(float64(resultsCount))
}

// RecordScheduleUpdate records a schedule status update
func RecordScheduleUpdate(fromStatus, toStatus string) {
	ScheduleStatusUpdates.WithLabelValues(fromStatus, toStatus).Inc()
}

// RecordUserRegistration records a new user registration
func RecordUserRegistration(userType string) {
	UserRegistrations.WithLabelValues(userType).Inc()
}

// RecordLoginAttempt records a login attempt
func RecordLoginAttempt(success bool) {
	if success {
		LoginAttempts.WithLabelValues("success").Inc()
	} else {
		LoginAttempts.WithLabelValues("failure").Inc()
	}
}
