package events

import app_ports "transcoder/internal/application/ports"

func WithVideoUploadedEvent(handler app_ports.Handler) Option {
	return func(manager EventsManager) {
		manager["video.events"] = handler
	}
}
