package app_kafka

import (
	"project/internal/core"

	"github.com/segmentio/kafka-go"
)

func NewProducer(cfg *core.KafkaConfig) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Address()),
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireAll,
	}

	return writer
}

func CloseProducer(producer *kafka.Writer) error {
	return producer.Close()
}
