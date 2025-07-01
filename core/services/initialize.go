package services

import (
	"github.com/medivue/adapters/config"
	"github.com/medivue/adapters/db/repository"
)

type Services struct {
	Core *ServicesHandler
}

// InitializeServices creates and returns all services
func InitializeServices(repos *repository.Repositories, cfg *config.Config) (*Services) {
	core := ServicesAdapter(
		repos.User,
		repos.Schedule,
		repos.Diagnostic,
		repos.Availability,
		repos.Record,
		repos.Payment,
		repos.Appointment,
		repos.TestPrice,
		*cfg,
	)
	return &Services{
		Core: core,
	}
}
