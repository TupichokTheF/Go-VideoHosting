package messaging

import (
	"context"
	"log/slog"
	"sync"
)

type Bus struct {
	consumer  consumerInterface
	handlers  map[string]Handler
	wg        sync.WaitGroup
	semaphore chan struct{}
}

type consumerInterface interface {
	FetchMessage() <-chan Message
}

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}

type Message struct {
	Type     string
	Payload  []byte
	Callback CommitMessage
}

type CommitMessage func(ctx context.Context) error

func NewBus(consumer consumerInterface, handlers map[string]Handler) *Bus {
	return &Bus{
		consumer:  consumer,
		handlers:  handlers,
		semaphore: make(chan struct{}, 10),
	}
}

func (bus *Bus) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			bus.wg.Wait()
			return nil
		case msg, ok := <-bus.consumer.FetchMessage():
			if !ok {
				bus.wg.Wait()
				return nil
			}
			bus.processMessage(ctx, msg)
		}
	}
}

func (bus *Bus) processMessage(ctx context.Context, msg Message) {
	handler, ok := bus.handlers[msg.Type]

	if !ok {
		return
	}

	bus.wg.Add(1)
	bus.semaphore <- struct{}{}
	go func() {
		defer func() {
			bus.wg.Done()
			<-bus.semaphore
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
