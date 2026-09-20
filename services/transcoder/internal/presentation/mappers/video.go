package mappers

import (
	"time"
	"transcoder/internal/application/dtos"

	"github.com/google/uuid"
)

func FromUploadedEventToDTO(payload map[string]any) dtos.TranscodeVideo {
	return dtos.TranscodeVideo{
		VideoID:   payload["video_id"].(uuid.UUID),
		OwnerID:   payload["owner_id"].(int),
		CreatedAt: payload["create_at"].(time.Time),
	}
}
