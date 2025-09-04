package db

import (
	"context"
	"fmt"
	"time"

	"github.com/diagnoxix/adapters/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBConfig holds database connection configuration
type DBConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnTimeout       time.Duration
}

// DefaultDBConfig returns default database configuration
func DefaultDBConfig() DBConfig {
	return DBConfig{
		MaxConns:          10,
		MinConns:          2,
		MaxConnLifetime:   30 * time.Minute,
		MaxConnIdleTime:   5 * time.Minute,
		HealthCheckPeriod: time.Minute,
		ConnTimeout:       15 * time.Second,
	}
}

// DatabaseConnection establishes a connection pool and returns a Queries instance.
func DatabaseConnection(
	ctx context.Context,
	dbURL string,
	cfg ...DBConfig,
) (*Queries, *pgxpool.Pool, error) {
	// Use provided config or default
	config := DefaultDBConfig()
	if len(cfg) > 0 {
		config = cfg[0]
	}

	// Use a context with timeout to avoid hanging connections
	ctx, cancel := context.WithTimeout(ctx, config.ConnTimeout)
	defer cancel()

	// Parse config to set additional settings
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = config.HealthCheckPeriod

	start := time.Now()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
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

// WithTx wraps a function with transaction handling
func WithTx[T any](ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) (T, error)) (T, error) {
	var result T

	// Start transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		metrics.DBConnectionErrors.Inc()
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			// Rollback on error
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				err = fmt.Errorf("error rolling back transaction: %v, original error: %w", rbErr, err)
			}
		} else {
			// Commit if no error
			err = tx.Commit(ctx)
			if err != nil {
				metrics.DBConnectionErrors.Inc()
				err = fmt.Errorf("error committing transaction: %w", err)
			}
		}
	}()

	// Execute transaction function
	result, err = fn(tx)
	return result, err
}

// MonitorConnectionHealth starts a goroutine that periodically checks database health
func MonitorConnectionHealth(pool *pgxpool.Pool, checkInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for range ticker.C {
			stats := pool.Stat()
			metrics.DBPoolActiveConnections.Set(float64(stats.AcquiredConns()))
			metrics.DBPoolIdleConnections.Set(float64(stats.IdleConns()))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := pool.Ping(ctx); err != nil {
				metrics.DBConnectionErrors.Inc()
				fmt.Printf("Database health check failed: %v\n", err)
			}
			cancel()
		}
	}()
}
