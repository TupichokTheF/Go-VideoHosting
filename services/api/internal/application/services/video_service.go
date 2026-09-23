package services

import (
	"context"
	"fmt"
	"log/slog"
	"project/internal/application/dtos"
	app_errors "project/internal/application/errors"
	app_ports "project/internal/application/ports"
	"project/internal/domain/event"
	"project/internal/domain/video"
	"time"
)

type VideoService struct {
	videoRepo    video.Repository
	videoStorage app_ports.Storage
	publisher    app_ports.Publisher
}

func NewVideoService(videoRepo video.Repository, videoStorage app_ports.Storage, publisher app_ports.Publisher) *VideoService {
	return &VideoService{
		videoRepo:    videoRepo,
		videoStorage: videoStorage,
		publisher:    publisher,
	}
}

func (videoService *VideoService) CreateVideo(ctx context.Context, videoData *dtos.CreateVideo) (*dtos.PresignedURL, error) {
	newVideo, err := video.New(videoData.OwnerID, videoData.Title, videoData.Description)
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	video_id, err := videoService.videoRepo.AddVideo(ctx, newVideo)
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	storageKey := fmt.Sprintf("videos/%s/source", newVideo.ID().String())

	url, err := videoService.videoStorage.PresignedURLCreate(ctx, storageKey)
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	return &dtos.PresignedURL{
		VideoID: video_id,
		URL:     url,
	}, nil
}

func (videoService *VideoService) GetVideo(ctx context.Context, videoData *dtos.GetVideo) (*dtos.PresignedURL, error) {
	video, err := videoService.videoRepo.GetVideoByID(ctx, videoData.VideoID)
	if err != nil {
		return nil, fmt.Errorf("get video: %w", err)
	}

	if err := video.CanBeShowed(); err != nil {
		return nil, fmt.Errorf("get video: %w", err)
	}

	key := fmt.Sprintf("videos/%s/source", videoData.VideoID.String())
	url, err := videoService.videoStorage.PresignedURLGet(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get video: %w", err)
	}

	return &dtos.PresignedURL{
		VideoID: videoData.VideoID.String(),
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
	
	key := fmt.Sprintf("videos/%s/source", v.ID().String())
	size, err := videoService.videoStorage.Stat(ctx, key)
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

	uploadedEvent := video.UploadedEvent{
		Base:    event.Base{CreatedAt: time.Now()},
		VideoID: v.ID(),
		OwnerID: v.OwnerID(),
	}
	if err := videoService.publisher.PublishEvents(ctx, []event.Interface{uploadedEvent}); err != nil {
		slog.Error("error from kafka while public events", "error", err)
	}

	return nil
}
