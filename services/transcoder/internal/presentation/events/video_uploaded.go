package pres_events

import (
	"context"
	"encoding/json"
	"fmt"
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
	wrappedHandler := middleware.OnFailure(publishFallback(publisher))(
		middleware.Retry(5, 1*time.Second)(
			middleware.OnSuccess(publishSuccess(publisher))(handler),
		),
	)

	return func(manager EventsManager) {
		manager["video.uploaded"] = wrappedHandler
	}
}

func publishFallback(publisher publisherInterface) func(ctx context.Context, msg messaging.Message) error {
	return func(ctx context.Context, msg messaging.Message) error {
		videoID := videoIDFromPayload(msg.Payload)

		event := video.UploadedFailedEvent{
			Base: domain_event.Base{
				CreatedAt: time.Now().UTC(),
			},
			VideoID: videoID,
		}

		if err := publisher.PublishEvents(ctx, []domain_event.Interface{&event}); err != nil {
			return fmt.Errorf("publish failed event: %w", err)
		}

		return nil
	}
}

func publishSuccess(publisher publisherInterface) func(ctx context.Context, msg messaging.Message) error {
	return func(ctx context.Context, msg messaging.Message) error {
		videoID := videoIDFromPayload(msg.Payload)

		event := video.UploadedSuccessEvent{
			Base: domain_event.Base{
				CreatedAt: time.Now().UTC(),
			},
			VideoID: videoID,
		}

		return publisher.PublishEvents(ctx, []domain_event.Interface{&event})
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
