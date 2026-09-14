package pres_ports

import (
	"context"
	"project/internal/application/dtos"
)

type VideoService interface {
	CreateVideo(ctx context.Context, videoData *dtos.CreateVideo) (*dtos.PresignedURL, error)
	GetVideo(ctx context.Context, videoData *dtos.GetVideo) (*dtos.PresignedURL, error)
	CompleteVideo(ctx context.Context, completeVideoData *dtos.CompleteVideo) error
}
