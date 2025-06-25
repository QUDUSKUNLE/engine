package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func DatabaseConnection(dbURL string) (*Queries, *pgx.Conn, error) {
	ctx := context.Background()

	// Create connection config
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Open a connection using pgx directly
	db, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify the connection
	if err := db.Ping(ctx); err != nil {
		// Close the connection if ping fails
		db.Close(ctx)
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := New(db)
	return queries, db, nil
}
