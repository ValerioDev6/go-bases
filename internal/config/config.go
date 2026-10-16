package config

import (
	"errors"
	"os"
	"strings"
)

const (
	defaultEnvironment = "dev"
	defaultPort        = "8080"
)

type Config struct {
	Environment string
	Port        string
	PostgresURL string
	JWTSecret   string
}

func Load() (Config, error) {
	cfg := Config{
		Environment: environment("ENV", defaultEnvironment),
		Port:        environment("PORT", defaultPort),
		PostgresURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret:   strings.TrimSpace(os.Getenv("JWT_SECRET")),
	}
	if cfg.PostgresURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must have at least 32 characters")
	}

	return cfg, nil
}

func (c Config) Address() string {
	return ":" + c.Port
}

func environment(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
