package services

import (
	"context"
	"fmt"
	"project/internal/application/dtos"
	app_errors "project/internal/application/errors"
	app_ports "project/internal/application/ports"
	"project/internal/domain/video"
)

type VideoService struct {
	videoRepo    video.Repository
	videoStorage app_ports.Storage
}

func NewVideoService(videoRepo video.Repository, videoStorage app_ports.Storage) *VideoService {
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

	if err := videoService.videoRepo.AddVideo(ctx, newVideo); err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	url, err := videoService.videoStorage.PresignedURLCreate(ctx, newVideo.ID().String())
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	return &dtos.PresignedURL{
		VideoID: newVideo.ID().String(),
		URL:     url,
	}, nil
}

func (videoService *VideoService) GetVideo(ctx context.Context, videoData *dtos.GetVideo) (*dtos.PresignedURL, error) {
	url, err := videoService.videoStorage.PresignedURLGet(ctx, videoData.VideoID)
	if err != nil {
		return nil, fmt.Errorf("get video: %w", err)
	}

	return &dtos.PresignedURL{
		VideoID: videoData.VideoID,
		URL:     url,
	}, nil
}

func (videoService *VideoService) CompleteVideo(ctx context.Context, completeVideoData *dtos.CompleteVideo) error {
	v, err := videoService.videoRepo.GetVideoByID(ctx, completeVideoData.VideoID)
	if err != nil {
		return fmt.Errorf("complete video: %w", err)
	}

	if completeVideoData.UserID != v.OwnerID() {
		return fmt.Errorf("complete video: %w", app_errors.ErrForbidden)
	}

	size, err := videoService.videoStorage.Stat(ctx, v.ID().String())
	if err != nil {
		return fmt.Errorf("complete video: %w", err)
	}
	if size == 0 {
		return fmt.Errorf("complete video: %w", video.ErrVideoNotLoaded)
	}

	if err := v.MarkUploaded(size); err != nil {
		return fmt.Errorf("complete video: %w", err)
	}

	if err := videoService.videoRepo.UpdateVideo(ctx, v); err != nil {
		return fmt.Errorf("complete video: %w", err)
	}

	return nil
}
