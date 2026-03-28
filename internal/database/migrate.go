package database

// Blank imports register golang-migrate's "file" source and "postgres" database drivers.

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// fixDirtyState checks for and clears any dirty migration state directly via SQL.
// This must happen before migrate.New(), which itself fails on dirty state.
func fixDirtyState(pool *pgxpool.Pool) {
	ctx := context.Background()

	var dirty bool
	var version int
	err := pool.QueryRow(ctx,
		"SELECT version, dirty FROM schema_migrations LIMIT 1",
	).Scan(&version, &dirty)

	if err != nil {
		// Table might not exist yet (first run)
		return
	}

	if dirty {
		log.Warn().
			Str("component", "migration").
			Int("version", version).
			Msg("dirty state detected, clearing")
		_, _ = pool.Exec(ctx, "UPDATE schema_migrations SET dirty = false")
	}
}

// RunMigrations applies all pending .sql files under migrations/ using dbURL.
// ErrNoChange is treated as success (already up to date).
func RunMigrations(pool *pgxpool.Pool, dbURL string) error {
	// Fix dirty state before migrate.New() touches the DB
	fixDirtyState(pool)

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().
				Str("component", "migration").
				Msg("no new migrations to apply")
			return nil
		}
		return err
	}

	log.Info().
		Str("component", "migration").
		Msg("applied successfully")
	return nil
}
