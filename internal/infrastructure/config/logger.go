package config

import (
	"log/slog"
	"os"
)

func setNewDefaultLogger(logLevel slog.Level) {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})),
	)
}
