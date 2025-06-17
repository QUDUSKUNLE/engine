package config

import (
	"log"
	"sync"

	"github.com/medicue/adapters/ex"
	"github.com/medicue/core/jobs"
	"github.com/medicue/core/ports"
)

// CronConfig manages cron jobs and services configuration
type CronConfig struct {
	ReminderJob *jobs.ReminderJob
	started     bool
	mu          sync.Mutex
}

var (
	once sync.Once
	cfg  *CronConfig
)

// GetConfig initializes and returns the configuration for cron jobs and services.
// It ensures that initialization happens only once using sync.Once.
func GetConfig(
	useRepo ports.UserRepository,
	diagnosticCentreRepo ports.DiagnosticRepository,
	appointmentPort ports.AppointmentRepository,
) *CronConfig {
	once.Do(func() {
		// Initialize email service
		// Initialize email service
		emailService := ex.NewSendGridEmailService()
		if emailService == nil {
			log.Fatal("Failed to initialize email service")
		}

		// Initialize repositories
		notificationSvc := &ex.EmailAdapter{}

		// Initialize reminder job
		reminderJob := jobs.NewReminderJob(
			appointmentPort,
			notificationSvc,
			useRepo,
			diagnosticCentreRepo,
		)

		if reminderJob == nil {
			log.Fatal("Failed to initialize reminder job")
		}

		cfg = &CronConfig{
			ReminderJob: reminderJob,
			started:     false,
		}
	})
	return cfg
}

// Start initializes and starts all services and jobs
func (c *CronConfig) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return nil // Already started
	}

	// Start reminder job
	c.ReminderJob.Start()

	c.started = true
	return nil
}

// Stop cleanly shuts down all services
func (c *CronConfig) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.started {
		return nil // Already stopped
	}

	if c.ReminderJob != nil {
		c.ReminderJob.Stop()
	}

	c.started = false
	return nil
}

// Cleanup performs necessary cleanup operations for the configuration
func (c *CronConfig) Cleanup() {
	_ = c.Stop() // We're shutting down anyway, ignore errors
}
