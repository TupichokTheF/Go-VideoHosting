package handlers

import (
	"context"
	"fmt"
	"transcoder/internal/presentation/mappers"
	"transcoder/internal/presentation/messaging"
	pres_ports "transcoder/internal/presentation/ports"
)

type VideoUploadedHandler struct {
	service pres_ports.VideoUploadedService
}

func NewVideoUploadedHandler(service pres_ports.VideoUploadedService) *VideoUploadedHandler {
	return &VideoUploadedHandler{
		service: service,
	}
}

func (handler *VideoUploadedHandler) Handle(ctx context.Context, msg messaging.Message) error {
	dto, err := mappers.FromUploadedEventToDTO(msg.Payload)
	if err != nil {
		return fmt.Errorf("error while video transcoding: %w", err)
	}

	err = handler.service.Transcode(ctx, dto)
	if err == nil {
		return nil
	}

	return fmt.Errorf("video transcoding handler: %w", err)
}
