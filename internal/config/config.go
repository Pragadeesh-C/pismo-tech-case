// Package config loads process configuration once from environment variables.
// DATABASE_URL is required; other fields have defaults (see Load).
package config

import (
	"errors"
	"os"
	"sync"
)

var (
	cfg  *Config
	once sync.Once
)

// Config holds server, logging, and database settings used at startup.
type Config struct {
	Server   Server
	Log      Log
	Database Database
}

type Log struct {
	Env   string
	Level string
}

type Server struct {
	Port    string
	GinMode string
}

type Database struct {
	URL string
}

// Get returns the config loaded by Load. Call Load first.
func Get() *Config {
	return cfg
}

// Load reads environment variables once (sync.Once). Safe for concurrent use after first call.
func Load() (*Config, error) {
	// Load environment variables from .env file. For local development. If docker compose is used, the environment variables are set in the docker compose file.
	var loadCfgErr error
	// Load environment variables once for concurrent use.
	once.Do(func() {
		// Set the config.
		cfg = &Config{
			Server: Server{
				Port:    envOrDefault("SERVER_PORT", "8080"),
				GinMode: envOrDefault("SERVER_GIN_MODE", "debug"),
			},
			Log: Log{
				Env:   envOrDefault("LOG_ENV", "production"),
				Level: envOrDefault("LOG_LEVEL", "info"),
			},
			Database: Database{
				URL: os.Getenv("DATABASE_URL"),
			},
		}

		if cfg.Database.URL == "" {
			loadCfgErr = errors.New("DATABASE_URL is required")
			return
		}
	})
	return cfg, loadCfgErr
}

// envOrDefault returns the environment variable value if it is set, otherwise it returns the default value.
func envOrDefault(env, defaultValue string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return defaultValue
}
