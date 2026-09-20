package dtos

import (
	"time"

	"github.com/google/uuid"
)

type TranscodeVideo struct {
	VideoID   uuid.UUID
	OwnerID   int
	CreatedAt time.Time
}
