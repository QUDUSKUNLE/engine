package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
)

type Repository struct {
	database *db.Queries
	conn     *pgx.Conn
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

func NewAppointmentRepository(
	store *db.Queries,
	conn *pgx.Conn,
) ports.AppointmentRepository {
	return &Repository{database: store, conn: conn}
}

func (r *Repository) GetTestTypes(ctx context.Context) ([]string, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT enumlabel 
        FROM pg_enum 
        WHERE enumtypid = (
            SELECT oid 
            FROM pg_type 
            WHERE typname = 'test_type'
        )
        ORDER BY enumsortorder
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testTypes []string
	for rows.Next() {
		var testType string
		if err := rows.Scan(&testType); err != nil {
			return nil, err
		}
		testTypes = append(testTypes, testType)
	}
	return testTypes, nil
}
