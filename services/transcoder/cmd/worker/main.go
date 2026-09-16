package main

import (
	"context"
	"log"
	"log/slog"
	"transcoder/internal/application/services"
	"transcoder/internal/core"
	app_kafka "transcoder/internal/infrastructure/brokers/kafka"
	"transcoder/internal/infrastructure/messaging"
	"transcoder/internal/presentation/events"
	"transcoder/internal/presentation/handlers"
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

	videoUploadedService := services.NewVideoUploadedService()

	videoUploadedHandler := handlers.NewVideoUploadedHandler(videoUploadedService)

	eventsHandlers := []events.Option{
		events.WithVideoUploadedEvent(videoUploadedHandler),
	}

	eventsManager := events.NewEventsManager(eventsHandlers...)

	bus := messaging.NewBus(kafkaConsumer, eventsManager)
	if err := bus.Listen(context.Background()); err != nil {
		log.Fatal("error while starting listen messages")
	}
}
