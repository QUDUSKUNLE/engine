package repository

import (
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
)

type Repository struct {
	database *db.Queries
}

func NewUserRepository(store *db.Queries) ports.UserRepository {
	return &Repository{database: store}
}

func NewScheduleRepository(store *db.Queries) ports.ScheduleRepository {
	return &Repository{database: store}
}

func NewDiagnosticCentreRepository(store *db.Queries) ports.DiagnosticRepository {
	return &Repository{database: store}
}

func NewRecordRepository(store *db.Queries) ports.RecordRepository {
	return &Repository{database: store}
}

func NewAvailabilityRepository(store *db.Queries) ports.AvailabilityRepository {
	return &Repository{database: store}
}
