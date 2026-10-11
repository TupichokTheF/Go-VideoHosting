package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"project/internal/app"
	"project/internal/core"
	"syscall"
)

// @title    Go Video Hosting API
// @version  1.0
// @host     localhost:8080
// @BasePath /api/v1
func main() {
	start()
}

func start() error {
	cfg := core.LoadConfig()

	logger := core.SetupLogger()
	slog.SetDefault(logger)

	app, err := app.SetupApplication(cfg)
	if err != nil {
		return fmt.Errorf("error while setting up application: %w", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			slog.Error("error while closing deps", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	return app.Run(ctx)
}
