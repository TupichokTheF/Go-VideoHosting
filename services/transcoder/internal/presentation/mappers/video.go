package mappers

import (
	"errors"
	"fmt"
	"transcoder/internal/application/dtos"

	"github.com/google/uuid"
)

func FromUploadedEventToDTO(payload map[string]any) (dtos.TranscodeVideo, error) {
	rawID, _ := payload["video_id"].(string)
    id, err := uuid.Parse(rawID)
    if err != nil {
        return dtos.TranscodeVideo{}, fmt.Errorf("video_id: %w", err)
    }

    ownerID, ok := payload["owner_id"].(float64)
    if !ok {
        return dtos.TranscodeVideo{}, errors.New("owner_id: not a number")
    }

    return dtos.TranscodeVideo{VideoID: id, OwnerID: int(ownerID)}, nil
}
