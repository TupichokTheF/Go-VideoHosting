package services

import (
	"context"
	"fmt"
	app_ports "transcoder/internal/application/ports"
)

type VideoUploadedService struct {
	transcoder app_ports.Transcoder
}

func NewVideoUploadedService(transcoder app_ports.Transcoder) *VideoUploadedService {
	return &VideoUploadedService{
		transcoder: transcoder,
	}
}

func (service *VideoUploadedService) Transcode(ctx context.Context) error {
	fmt.Println("Transcoding...")

	return nil
}
