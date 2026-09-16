package app_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	app_ports "transcoder/internal/application/ports"
	"transcoder/internal/core"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader

	msgChan chan app_ports.Message
}

func NewConsumer(cfg *core.KafkaConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Address()},
		GroupID:     strconv.Itoa(cfg.GroupID),
		GroupTopics: cfg.Topics,
	})

	return &Consumer{
		reader:  reader,
		msgChan: make(chan app_ports.Message, 1),
	}
}

func (consumer *Consumer) FetchMessage() <-chan app_ports.Message {
	return consumer.msgChan
}

func (consumer *Consumer) StartReading(ctx context.Context) error {
	for {
		kafkaMsg, err := consumer.reader.ReadMessage(ctx)
		if err != nil {
			close(consumer.msgChan)
			return fmt.Errorf("error while reading messages from kafka: %w", err)
		}

		var payload = make(map[string]any)
		if err := json.Unmarshal(kafkaMsg.Value, &payload); err != nil {
			slog.Warn("invalid data of message")
			continue
		}

		msg := app_ports.Message{
			Type:    string(kafkaMsg.Value),
			Payload: payload,
			Callback: func(ctx context.Context) error {
				if err := consumer.reader.CommitMessages(ctx, kafkaMsg); err != nil {
					return fmt.Errorf("error while commiting message: %w", err)
				}

				return nil
			},
		}

		consumer.msgChan <- msg
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
