package main

import (
	"transcoder/internal/core"
	app_kafka "transcoder/internal/infrastructure/brokers/kafka"
	"transcoder/internal/infrastructure/messaging"
)


func main() {
	start()
}

func start() {
	cfg := core.LoadConfig()
	
	kafkaConsumer := app_kafka.NewConsumer(&cfg.KafkaConfig)

	

	bus := messaging.NewBus(kafkaConsumer)
}