package messaging


type EventsManager map[string]Handler

type Option func(EventsManager)

func NewEventsManager(opts ...Option) EventsManager {
	manager := make(EventsManager)

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

func WithVideoUploadedEvent(handler Handler) Option {
	return func(manager EventsManager) {
		manager["video.uploaded"] = handler
	}
}
