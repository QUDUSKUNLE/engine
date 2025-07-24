package services

import (
	"github.com/medivue/adapters/config"
	"github.com/medivue/adapters/db/repository"
)

type Service struct {
	Core *ServicesHandler
}

// InitializeServices creates and returns all services
func InitializeServices(repos *repository.Repositories, cfg *config.EnvConfiguration) *Service {
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
	return &Service{
		Core: core,
	}
}
