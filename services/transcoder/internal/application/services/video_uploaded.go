package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"transcoder/internal/application/dtos"
	app_ports "transcoder/internal/application/ports"
	"transcoder/internal/domain/video"
)

type VideoUploadedService struct {
	transcoder app_ports.Transcoder
	storage    app_ports.Storage
}

func NewVideoUploadedService(transcoder app_ports.Transcoder, storage app_ports.Storage) *VideoUploadedService {
	return &VideoUploadedService{
		transcoder: transcoder,
		storage:    storage,
	}
}

func (service *VideoUploadedService) Transcode(ctx context.Context, videoData dtos.TranscodeVideo) error {
	outputKey := fmt.Sprintf("/vide/%s/720p.mp4", videoData.VideoID)
	sourceKey := fmt.Sprintf("/video/%s/source", videoData.VideoID)

	ok, err := service.storage.IsExist(ctx, outputKey)
	if ok == false {
		if err != nil {
			return fmt.Errorf("transcode video service: %w", err)
		}

		return fmt.Errorf("transcode video service: %w", video.ErrNotFound)
	}

	dir, err := os.MkdirTemp("", "transcode")
	if err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}

	srcDir := filepath.Join(dir, "source")
	dstDir := filepath.Join(dir, "720p")

	if err := service.storage.Download(ctx, sourceKey, srcDir); err != nil {
		return fmt.Errorf("trnascode video service: %w", err)
	}
	if err := service.transcoder.Transcode(ctx, srcDir, dstDir); err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}
	if err := service.storage.Upload(ctx, outputKey, dstDir); err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}

	return nil
}
