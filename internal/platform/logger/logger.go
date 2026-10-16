package logger

import (
	"log/slog"
	"os"
)

func New(environment string) *slog.Logger {
	options := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: environment == "dev",
	}

	if environment == "prod" {
		return slog.New(slog.NewJSONHandler(os.Stdout, options))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, options))
}
