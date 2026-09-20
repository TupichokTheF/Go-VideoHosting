package events

import "transcoder/internal/presentation/messaging"

type EventsManager map[string]messaging.Handler

type Option func(EventsManager)

func NewEventsManager(opts ...Option) EventsManager {
	manager := make(EventsManager)

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}
