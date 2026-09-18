package main

import (
	"context"
	"log/slog"
	"transcoder/internal/application/services"
	"transcoder/internal/core"
	"transcoder/internal/presentation/events"
	"transcoder/internal/presentation/handlers"
	"transcoder/internal/presentation/messaging"
	app_kafka "transcoder/internal/presentation/messaging/brokers/kafka"

	"golang.org/x/sync/errgroup"
)

func main() {
	start()
}

func start() error {
	cfg := core.LoadConfig()

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

	videoUploadedService := services.NewVideoUploadedService()

	videoUploadedHandler := handlers.NewVideoUploadedHandler(videoUploadedService)

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
