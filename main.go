package main

import (
	"errors"
	"log/slog"
	"os"

	"web/internal/app"
	"web/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("load .env", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load configuration", "error", err)
		os.Exit(1)
	}
	if err := app.Run(cfg); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}
