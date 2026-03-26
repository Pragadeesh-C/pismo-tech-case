// Package server builds the Gin HTTP engine, registers routes, and manages graceful shutdown.
package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragadeesh-c/pismo-tech-case/internal/api/middleware"
	"github.com/pragadeesh-c/pismo-tech-case/internal/config"
	"github.com/rs/zerolog/log"
)

// Server wraps the standard library HTTP server configured with Gin.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a Gin engine with recovery and CORS, registers routes, and returns a Server.
// The DB pool parameter is reserved for route handlers that will be wired in next.
func NewServer(cfg *config.Config, _ *pgxpool.Pool) *Server {
	gin.SetMode(cfg.Server.GinMode)

	r := gin.New()
	r.Use(
		gin.Recovery(),
		middleware.CORSMiddleware(),
	)

	r.GET("/health", HealthHandler())

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: r,
		},
	}
}

// Start runs ListenAndServe in a goroutine registered on rootCtx's WaitGroup.
// When ctx is cancelled (signal), it shuts down the HTTP server with a 5s timeout.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("addr", s.httpServer.Addr).Msg("server started")

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start server")
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("graceful shutdown failed")
			return err
		}
		log.Info().Msg("server stopped gracefully")
		return nil
	case err := <-errCh:
		return err
	}
}
