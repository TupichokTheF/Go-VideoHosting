package events

import "transcoder/internal/presentation/messaging"

func WithVideoUploadedEvent(handler messaging.Handler) Option {
	return func(manager EventsManager) {
		manager["video.uploaded"] = handler
	}
}
