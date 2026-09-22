package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"transcoder/internal/application/dtos"
	app_ports "transcoder/internal/application/ports"
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

func (service *VideoUploadedService) Transcode(ctx context.Context, videoData dtos.TranscodeVideo) (err error) {
	defer func() {
		err = isUnprocessable(err)
	}()

	outputKey := fmt.Sprintf("videos/%s/720p.mp4", videoData.VideoID)
	sourceKey := fmt.Sprintf("videos/%s/source", videoData.VideoID)

	ok, err := service.storage.IsExist(ctx, outputKey)
	if err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}
	if ok {
		return nil
	}

	dir, err := os.MkdirTemp("", "transcode")
	if err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}
	defer os.RemoveAll(dir)

	srcDir := filepath.Join(dir, "source")
	dstDir := filepath.Join(dir, "720p.mp4")

	if err := service.storage.Download(ctx, sourceKey, srcDir); err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}
	if err := service.transcoder.Transcode(ctx, srcDir, dstDir); err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}
	if err := service.storage.Upload(ctx, outputKey, dstDir); err != nil {
		return fmt.Errorf("transcode video service: %w", err)
	}

	return nil
}

func isUnprocessable(err error) error {
	switch {
	case errors.Is(err, app_ports.ErrObjectNotFound):
		return fmt.Errorf("transcode video service: %w", app_ports.ErrUnprocessable)
	default:
		return err
	}
}
