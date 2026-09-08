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
	return "video.topic"
}

func (event UploadedEvent) Payload() map[string]any {
	return map[string]any{
		"video.id": event.VideoID,
		"owner.id": event.OwnerID,
	}
}