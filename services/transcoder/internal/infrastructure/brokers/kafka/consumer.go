package app_kafka

import (
	"strconv"
	"transcoder/internal/core"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *core.KafkaConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Address()},
		GroupID:     strconv.Itoa(cfg.GroupID),
		GroupTopics: cfg.Topics,
	})

	return &Consumer{
		reader: reader,
	}
}

func (consumer *Consumer) FetchMessage() error {
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
