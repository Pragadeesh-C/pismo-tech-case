// Package database provides PostgreSQL access: a pgx connection pool (NewPool) and
// golang-migrate file-based migrations (RunMigrations, migrations/ on disk).
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool parses databaseURL, applies pool limits, pings once, and returns the pool.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating pool: %w", err)
	}

	// Ping the database to check if the connection is successful.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return pool, nil
}
