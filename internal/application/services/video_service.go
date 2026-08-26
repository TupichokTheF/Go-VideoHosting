package services

import (
	"context"
	"fmt"
	"project/internal/application/dtos"
	"project/internal/domain/video"
	infra_ports "project/internal/application/ports"
	"strconv"
)

type VideoService struct {
	videoRepo    video.Repository
	videoStorage infra_ports.Storage
}

func NewVideoService(videoRepo video.Repository, videoStorage infra_ports.Storage) *VideoService {
	return &VideoService{
		videoRepo:    videoRepo,
		videoStorage: videoStorage,
	}
}

func (videoService *VideoService) CreateVideo(ctx context.Context, videoData *dtos.CreateVideo) (*dtos.PresignedURL, error) {
	newVideo, err := video.New(videoData.OwnerID, videoData.Title, videoData.Description)
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	videoID, err := videoService.videoRepo.AddVideo(ctx, newVideo)
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	url, err := videoService.videoStorage.PresignedURLCreate(ctx, strconv.Itoa(videoID))
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	return &dtos.PresignedURL{
		URL: url,
	}, nil
}

func (videoService *VideoService) GetVideo(ctx context.Context, videoData *dtos.GetVideo) (*dtos.PresignedURL, error) {
	url, err := videoService.videoStorage.PresignedURLGet(ctx, strconv.Itoa(videoData.VideoID))
	if err != nil {
		return nil, fmt.Errorf("get video: %w", err)
	}

	return &dtos.PresignedURL{
		URL: url,
	}, nil
}
