package messaging

import "transcoder/internal/presentation/events"

type Bus struct {
	consumer consumerInterface
	handlers events.EventsManager
}

type consumerInterface interface {
	FetchMessage() error
}

func NewBus(consumer consumerInterface, handlers events.EventsManager) *Bus {
	return &Bus{
		consumer: consumer,
		handlers: handlers,
	}
}

func (bus *Bus) ProcessMessage() error {
	return nil
}
