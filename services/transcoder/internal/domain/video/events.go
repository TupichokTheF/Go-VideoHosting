package video

import (
	domain_event "transcoder/internal/domain/event"

	"github.com/google/uuid"
)

type UploadedFailedEvent struct {
	domain_event.Base
	VideoID uuid.UUID
}

func (_ *UploadedFailedEvent) EventName() string {
	return "video.uploaded.failed"
}

func (event *UploadedFailedEvent) Payload() map[string]any {
	return map[string]any{
		"video_id":   event.VideoID,
		"created_at": event.CreatedAt,
	}
}
