package main

import (
	"transcoder/internal/application/services"
	"transcoder/internal/core"
	app_kafka "transcoder/internal/infrastructure/brokers/kafka"
	"transcoder/internal/infrastructure/messaging"
	"transcoder/internal/presentation/events"
	"transcoder/internal/presentation/handlers"
	pres_ports "transcoder/internal/presentation/ports"
)

func main() {
	start()
}

func start() {
	cfg := core.LoadConfig()

	kafkaConsumer := app_kafka.NewConsumer(&cfg.KafkaConfig)

	videoUploadedService := services.NewVideoUploadedService()

	videoUploadedHandler := handlers.NewVideoUploadedHandler(videoUploadedService)

	eventsHandlers := []events.Option{
		events.WithVideoUploadedEvent(videoUploadedHandler),
	}

	eventsManager := events.NewEventsManager(eventsHandlers...)

	bus := messaging.NewBus(kafkaConsumer, eventsManager)
}
