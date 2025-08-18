package repository

import (
	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *db.Queries
	conn     *pgxpool.Pool
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
	conn *pgxpool.Pool,
) ports.PaymentRepository {
	return &Repository{database: store, conn: conn}
}

func NewAppointmentRepository(
	store *db.Queries,
	conn *pgxpool.Pool,
) ports.AppointmentRepository {
	return &Repository{database: store, conn: conn}
}

func NewDiagnosticCentreRepository(
	store *db.Queries,
	conn *pgxpool.Pool,
) ports.DiagnosticRepository {
	return &Repository{database: store, conn: conn}
}
