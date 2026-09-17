package main

import (
	"context"
	"log"
	"log/slog"
	"transcoder/internal/application/services"
	"transcoder/internal/core"
	"transcoder/internal/presentation/events"
	"transcoder/internal/presentation/handlers"
	"transcoder/internal/presentation/messaging"
	app_kafka "transcoder/internal/presentation/messaging/brokers/kafka"
)

func main() {
	start()
}

func start() {
	cfg := core.LoadConfig()

	kafkaConsumer := app_kafka.NewConsumer(&cfg.KafkaConfig)
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			slog.Error("error while closing kafka consumer", "error", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := kafkaConsumer.StartReading(ctx); err != nil {
			slog.Error("error while starting reading messages from kafka", "error", err)
			cancel()
		}
	}()

	videoUploadedService := services.NewVideoUploadedService()

	videoUploadedHandler := handlers.NewVideoUploadedHandler(videoUploadedService)

	eventsHandlers := []events.Option{
		events.WithVideoUploadedEvent(videoUploadedHandler),
	}

	eventsManager := events.NewEventsManager(eventsHandlers...)

	bus := messaging.NewBus(kafkaConsumer, eventsManager)
	if err := bus.Listen(ctx); err != nil {
		log.Fatal("error while starting listen messages")
	}
}
