package messaging

import (
	"context"
	"project/internal/domain/event"

	"github.com/segmentio/kafka-go"
)


type Publisher struct {
	consumer *kafka.Writer
}

func NewPublisher(consumer *kafka.Writer) *Publisher {
	return &Publisher{
		consumer: consumer,
	}
}

func PublicEvents(ctx context.Context, events []event.Interface) error {
	return nil
}