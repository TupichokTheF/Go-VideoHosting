package core

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(handler)
}
