package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConnection establishes a connection pool and returns a Queries instance.
// The caller is responsible for closing the returned pool.
func DatabaseConnection(ctx context.Context, dbURL string) (*Queries, *pgxpool.Pool, error) {
	// Use a context with timeout to avoid hanging connections
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Parse config to set additional settings
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set connection pool settings
	config.MaxConns = 10                       // maximum number of connections in the pool
	config.MinConns = 2                        // minimum number of connections in the pool
	config.MaxConnLifetime = 30 * time.Minute  // maximum lifetime of a connection
	config.MaxConnIdleTime = 5 * time.Minute   // maximum idle time for a connection
	config.HealthCheckPeriod = 1 * time.Minute // how often to check health of idle connections

	start := time.Now()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	connTime := time.Since(start)

	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Printf("Database connected in %v\n", connTime)

	queries := New(pool)
	return queries, pool, nil
}
