package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func DatabaseConnection(DB_URL string) (*Queries, *pgx.Conn, error) {
	// Open a connection using pgx directly
	db, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening database connection: %v", err)
	}

	return New(db), db, nil
}
