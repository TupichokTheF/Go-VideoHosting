package app_kafka

import (
	"project/internal/core"

	"github.com/segmentio/kafka-go"
)


func NewProducer(cfg *core.KafkaConfig) *kafka.Writer {
	kafkaConfig := kafka.WriterConfig{
		Brokers: []string{cfg.Address()},
	}

	writer := kafka.NewWriter(kafkaConfig)

	return writer
}

func CloseProducer(producer *kafka.Writer) error {
	return producer.Close()
}