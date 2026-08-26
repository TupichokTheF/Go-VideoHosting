package video

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddVideo(ctx context.Context, video *Video) error
	GetVideoByID(ctx context.Context, videoID uuid.UUID) (*Video, error)
	UpdateVideo(ctx context.Context, video *Video) error
}
