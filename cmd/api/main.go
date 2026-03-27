package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/pragadeesh-c/pismo-tech-case/internal/config"
	"github.com/pragadeesh-c/pismo-tech-case/internal/database"
	"github.com/pragadeesh-c/pismo-tech-case/internal/logger"
	"github.com/pragadeesh-c/pismo-tech-case/internal/server"
	"github.com/rs/zerolog/log"
)

// Command api is the HTTP service entrypoint: loads config, connects to PostgreSQL,
// and runs the Gin server until SIGINT/SIGTERM.
func main() {
	// Load environment variables from .env file. For local development. If docker compose is used, the environment variables are set in the docker compose file.
	_ = godotenv.Load(".env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load config")
	}

	// Initialize the logger.
	logger.New(cfg.Log.Env, cfg.Log.Level)

	// Create a channel to receive the shutdown signal.
	sigs := make(chan os.Signal, 1)
	// Notify the channel when SIGINT or SIGTERM is received.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Root WaitGroup is passed through context so the overall application can register its goroutines with main's lifecycle.
	rootCtx, cancel := context.WithCancel(context.Background())

	// Shutdown signal handler. Handles SIGINT and SIGTERM.
	go func() {
		sig := <-sigs
		log.Info().Msgf("received signal -> %v", sig)
		cancel()
	}()

	// Connect to the database.
	pool, err := database.NewPool(rootCtx, cfg.Database.URL)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect database")
	}
	defer pool.Close()

	if err := database.RunMigrations(pool, cfg.Database.URL); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to run migration")
	}

	// Create a new server struct instance.
	srv := server.NewServer(cfg, pool)

	// Start the API server.
	if err := srv.Start(rootCtx); err != nil {
		log.Fatal().
			Err(err).
			Msg("server exited with error")
	}
}
