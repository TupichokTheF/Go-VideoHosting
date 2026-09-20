package schemas

import (
	"time"

	"github.com/google/uuid"
)

type VideoUploadedSchema struct {
	VideoID   uuid.UUID
	OwnerID   int
	CreatedAt time.Time
}
