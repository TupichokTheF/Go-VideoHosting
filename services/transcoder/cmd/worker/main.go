package main

import (
	"context"
	"fmt"
	"log/slog"
	"transcoder/internal/application/services"
	"transcoder/internal/core"
	app_kafka "transcoder/internal/infrastructure/brokers/kafka"
	"transcoder/internal/infrastructure/storage"
	"transcoder/internal/infrastructure/transcoder"
	"transcoder/internal/presentation/events"
	"transcoder/internal/presentation/handlers"
	"transcoder/internal/presentation/messaging"

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

	minioClient, err := storage.NewMinioClient(&cfg.MinioConfig)
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	transcoder := transcoder.NewTranscoderCmd(cfg.TranscoderConfig.Bin)
	storage := storage.NewMinioService(minioClient, cfg.MinioConfig.Bucket)
	kafkaProducer := app_kafka.NewProducer(&cfg.KafkaConfig)

	videoUploadedService := services.NewVideoUploadedService(transcoder, storage)

	publisher := messaging.NewPublisher(kafkaProducer)

	videoUploadedHandler := handlers.NewVideoUploadedHandler(videoUploadedService, publisher)

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
