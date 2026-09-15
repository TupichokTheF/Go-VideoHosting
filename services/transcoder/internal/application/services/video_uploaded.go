package services

import (
	"context"
	"fmt"
)

type VideoUploadedService struct {
}

func NewVideoUploadedService() *VideoUploadedService {
	return &VideoUploadedService{}
}

func (service *VideoUploadedService) Transcode(ctx context.Context) error {
	fmt.Println("Transcoding...")

	return nil
}
