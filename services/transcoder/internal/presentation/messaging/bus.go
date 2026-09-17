package messaging

import (
	"context"
	"log/slog"
	"sync"
	app_ports "transcoder/internal/application/ports"
	"transcoder/internal/presentation/events"
)

type Bus struct {
	consumer consumerInterface
	handlers events.EventsManager
	wg sync.WaitGroup
	semaphore chan struct{}
}

type consumerInterface interface {
	FetchMessage() <-chan app_ports.Message
}

func NewBus(consumer consumerInterface, handlers events.EventsManager) *Bus {
	return &Bus{
		consumer: consumer,
		handlers: handlers,
		semaphore: make(chan struct{}, 10),
	}
}

func (bus *Bus) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			bus.wg.Wait()
			return nil
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
	
	bus.wg.Add(1)
	bus.semaphore <- struct{}{}
	go func() {
		defer func() {
			bus.wg.Done()
			<- bus.semaphore
		}()

		if err := handler.Handle(ctx, msg); err != nil {
			slog.Error("error while handle message", "error", err)

			return
		}

		if err := msg.Callback(ctx); err != nil {
			slog.Error("error while commiting message", "error", err)
		}
	}()
}
