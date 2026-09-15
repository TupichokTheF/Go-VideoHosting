package events

import (
	app_ports "transcoder/internal/application/ports"
)


type EventsManager map[string]app_ports.Handler

type Option func(EventsManager)

func NewEventsManager(opts ...Option) EventsManager {
	manager := make(EventsManager)
	
	for _, opt := range opts {
		opt(manager)
	}

	return manager
}