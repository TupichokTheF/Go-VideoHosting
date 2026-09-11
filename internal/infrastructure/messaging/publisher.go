package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"project/internal/domain/event"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	producer *kafka.Writer
}

func NewPublisher(producer *kafka.Writer) *Publisher {
	return &Publisher{
		producer: producer,
	}
}

func (p *Publisher) PublicEvents(ctx context.Context, events []event.Interface) error {
	if len(events) == 0 {
		return nil
	}

	messages := make([]kafka.Message, 0, len(events))

	for _, event := range events {
		body, err := json.Marshal(event.Payload())
		if err != nil {
			return fmt.Errorf("public events: %w", err)
		}

		messages = append(messages, kafka.Message{
			Topic: p.topicFromEventName(event.EventName()),
			Value: body,
			Time:  event.OccurredAt(),
		})
	}

	if err := p.producer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("public events: %w", err)
	}

	return nil
}

func (p *Publisher) topicFromEventName(event string) string {
	entityName := strings.Split(event, ".")
	topicName := entityName[0] + ".events"

	return topicName
}
