package video

import (
	"project/internal/domain/event"

	"github.com/google/uuid"
)

type UploadedEvent struct {
	event.Base
	VideoID uuid.UUID
	OwnerID int
}

func (_ UploadedEvent) EventName() string {
	return "video.uploaded"
}

func (event UploadedEvent) Payload() map[string]any {
	return map[string]any{
		"video_id": event.VideoID,
		"owner_id": event.OwnerID,
	}
}
