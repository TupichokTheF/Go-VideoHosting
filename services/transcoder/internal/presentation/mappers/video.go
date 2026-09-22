package mappers

import (
	"encoding/json"
	"fmt"
	"time"
	"transcoder/internal/application/dtos"

	"github.com/google/uuid"
)

type videoUploadedPayload struct {
	VideoID   uuid.UUID `json:"video_id"`
	OwnerID   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func FromUploadedEventToDTO(payload []byte) (dtos.TranscodeVideo, error) {
	var uploadedPayload videoUploadedPayload
	if err := json.Unmarshal(payload, &uploadedPayload); err != nil {
		return dtos.TranscodeVideo{}, fmt.Errorf("transform from uploaded video payload to dto: %w", err)
	}

	return dtos.TranscodeVideo{
		VideoID:   uploadedPayload.VideoID,
		OwnerID:   uploadedPayload.OwnerID,
		CreatedAt: uploadedPayload.CreatedAt,
	}, nil
}
