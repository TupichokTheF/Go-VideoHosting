package app_kafka

import (
	"context"
	"fmt"
	"transcoder/internal/core"
	"transcoder/internal/presentation/messaging"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader

	msgChan chan messaging.Message
}

func NewConsumer(cfg *core.KafkaConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Address()},
		GroupID:     cfg.GroupID,
		GroupTopics: cfg.Topics,
	})

	return &Consumer{
		reader:  reader,
		msgChan: make(chan messaging.Message, 1),
	}
}

func (consumer *Consumer) FetchMessage() <-chan messaging.Message {
	return consumer.msgChan
}

func (consumer *Consumer) StartReading(ctx context.Context) error {
	defer close(consumer.msgChan)
	for {
		kafkaMsg, err := consumer.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error while reading messages from kafka: %w", err)
		}

		msg := messaging.Message{
			Type:    typeOfMessage(&kafkaMsg),
			Payload: kafkaMsg.Value,
			Callback: func(ctx context.Context) error {
				if err := consumer.reader.CommitMessages(ctx, kafkaMsg); err != nil {
					return fmt.Errorf("error while commiting message: %w", err)
				}

				return nil
			},
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("start reading: %w", ctx.Err())
		case consumer.msgChan <- msg:
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func typeOfMessage(msg *kafka.Message) string {
	for _, header := range msg.Headers {
		if header.Key == "event_type" {
			return string(header.Value)
		}
	}

	return ""
}
