package main

import (
	"context"
	"log/slog"
	"transcoder/internal/application/services"
	"transcoder/internal/core"
	"transcoder/internal/infrastructure/transcoder"
	"transcoder/internal/presentation/events"
	"transcoder/internal/presentation/handlers"
	"transcoder/internal/presentation/messaging"
	app_kafka "transcoder/internal/presentation/messaging/brokers/kafka"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := start(); err != nil {
		slog.Error("error in main", "error", err)
	}
}

func start() error {
	cfg := core.LoadConfig()

	logger := core.SetupLogger()
	slog.SetDefault(logger)

	transcoder := transcoder.NewTranscoderCmd(cfg.TranscoderConfig.Bin)

	videoUploadedService := services.NewVideoUploadedService(transcoder)

	videoUploadedHandler := handlers.NewVideoUploadedHandler(videoUploadedService)

	kafkaConsumer := app_kafka.NewConsumer(&cfg.KafkaConfig)
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			slog.Error("error while closing kafka consumer", "error", err)
		}
	}()

	gErr, ctx := errgroup.WithContext(context.Background())

	gErr.Go(func() error {
		return kafkaConsumer.StartReading(ctx)
	})

	eventsHandlers := []events.Option{
		events.WithVideoUploadedEvent(videoUploadedHandler),
	}

	eventsManager := events.NewEventsManager(eventsHandlers...)

	bus := messaging.NewBus(kafkaConsumer, eventsManager)
	gErr.Go(func() error {
		return bus.Listen(ctx)
	})

	return gErr.Wait()
}
