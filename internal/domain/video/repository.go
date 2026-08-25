package video

import (
	"context"
)

type Repository interface {
	AddVideo(ctx context.Context, video *Video) (int, error)
}
