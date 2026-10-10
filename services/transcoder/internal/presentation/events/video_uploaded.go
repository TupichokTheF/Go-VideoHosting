package pres_events

import (
	"context"
	"encoding/json"
	"time"
	domain_event "transcoder/internal/domain/event"
	"transcoder/internal/domain/video"
	"transcoder/internal/presentation/messaging"
	"transcoder/internal/presentation/middleware"

	"github.com/google/uuid"
)

type publisherInterface interface {
	PublishEvents(ctx context.Context, events []domain_event.Interface) error
}

func WithVideoUploadedEvent(handler messaging.Handler, publisher publisherInterface) Option {
	return func(manager EventsManager) {
		handlerWithRetry := func(ctx context.Context, msg messaging.Message) error {
			wrappedHandler := middleware.Retry(5, 50*time.Millisecond)(handler)

			if err := wrappedHandler.Handle(ctx, msg); err != nil {
				videoID := videoIDFromPayload(msg.Payload)
				eventToPublish := video.UploadedFailedEvent{
					Base: domain_event.Base{
						CreatedAt: time.Now(),
					},
					VideoID: videoID,
				}

				publisher.PublishEvents(ctx, []domain_event.Interface{&eventToPublish})
			}

			return nil
		}

		manager["video.uploaded"] = messaging.HandlerFunc(handlerWithRetry)
	}
}

func videoIDFromPayload(payload []byte) uuid.UUID {
	var videoID videoIDPayload
	_ = json.Unmarshal(payload, &videoID)

	return videoID.VideoID
}

type videoIDPayload struct {
	VideoID uuid.UUID `json:"video_id"`
}
