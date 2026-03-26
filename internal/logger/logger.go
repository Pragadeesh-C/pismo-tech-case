// Package logger configures the global zerolog logger (level, format, service label).
package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New sets log.Logger. Use env "development" for human-readable console output; otherwise JSON to stdout.
// Invalid logLevel falls back to info.
func New(env, logLevel string) {
	zerolog.TimeFieldFormat = time.RFC3339

	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}

	e := strings.ToLower(env)

	var logger zerolog.Logger
	if e == "development" {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05.000",
			NoColor:    false,
		})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	log.Logger = logger.
		Level(level).
		With().
		Timestamp().
		Caller().
		Str("service", "pismo-tech-case").
		Logger()

}
