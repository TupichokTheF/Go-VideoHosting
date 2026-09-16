package messaging

import (
	"context"
	"fmt"
	"log/slog"
	app_ports "transcoder/internal/application/ports"
	"transcoder/internal/presentation/events"
)

type Bus struct {
	consumer consumerInterface
	handlers events.EventsManager
}

type consumerInterface interface {
	FetchMessage() <-chan app_ports.Message
}

func NewBus(consumer consumerInterface, handlers events.EventsManager) *Bus {
	return &Bus{
		consumer: consumer,
		handlers: handlers,
	}
}

func (bus *Bus) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Context was cancelled")
		case msg := <-bus.consumer.FetchMessage():
			bus.processMessage(ctx, msg)
		}
	}
}

func (bus *Bus) processMessage(ctx context.Context, msg app_ports.Message) {
	handler, ok := bus.handlers[msg.Type]

	if !ok {
		return
	}

	go func() {
		if err := handler.Handle(ctx, msg); err != nil {
			slog.Error("error while handle message", "error", err)

			return
		}

		msg.Callback(ctx)
	}()
}
