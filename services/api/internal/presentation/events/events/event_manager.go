package pres_events

import pres_messaging "project/internal/presentation/events/messaging"

type EventManager map[string]pres_messaging.Handler

type Option func(EventManager)

func NewEventManager(opts ...Option) EventManager {
	eventManager := make(EventManager)

	for _, opt := range opts {
		opt(eventManager)
	}

	return eventManager
}
