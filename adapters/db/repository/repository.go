package repository

import (
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
	"context"
)

type Repository struct {
	database *db.Queries
}

func NewUserRepository(
	store *db.Queries,
) ports.UserRepository {
	return &Repository{database: store}
}

func NewScheduleRepository(
	store *db.Queries,
) ports.ScheduleRepository {
	return &Repository{database: store}
}

func NewDiagnosticCentreRepository(
	store *db.Queries,
) ports.DiagnosticRepository {
	return &Repository{database: store}
}

func NewRecordRepository(
	store *db.Queries,
) ports.RecordRepository {
	return &Repository{database: store}
}

func NewAvailabilityRepository(
	store *db.Queries,
) ports.AvailabilityRepository {
	return &Repository{database: store}
}

func NewPaymentRepository(
	store *db.Queries,
) ports.PaymentRepository {
	return &Repository{database: store}
}

func (r *Repository) BeginTx(ctx context.Context) (ports.AppointmentTx, error) {
	// Implementation needed
	return nil, nil
}

func NewApppointmentRepository(
	store *db.Queries,
) ports.AppointmentRepository {
	return &Repository{database: store}
}
