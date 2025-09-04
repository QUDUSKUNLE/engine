package repository

import (
	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides access to the database
type Repository struct {
	database *db.Queries
	pool     *pgxpool.Pool
}

// NewRepository creates a new repository instance
func NewRepository(
	database *db.Queries,
	pool *pgxpool.Pool,
) *Repository {
	return &Repository{
		database: database,
		pool:     pool,
	}
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

func NewTestPriceRepository(
	store *db.Queries,
) ports.TestPriceRepository {
	return &Repository{database: store}
}

func NewPaymentRepository(
	store *db.Queries,
	pool *pgxpool.Pool,
) ports.PaymentRepository {
	return &Repository{database: store, pool: pool}
}

func NewAppointmentRepository(
	store *db.Queries,
	pool *pgxpool.Pool,
) ports.AppointmentRepository {
	return &Repository{database: store, pool: pool}
}

func NewDiagnosticCentreRepository(
	store *db.Queries,
	pool *pgxpool.Pool,
) ports.DiagnosticRepository {
	return &Repository{database: store, pool: pool}
}
