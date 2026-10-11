package pres_events

import pres_messaging "project/internal/presentation/events/messaging"

func WithUplaodedSuccess(handler pres_messaging.Handler) Option {
	return func(manager EventManager) {
		manager["video.uploaded.success"] = handler
	}
}
