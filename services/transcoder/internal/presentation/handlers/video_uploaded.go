package handlers

import (
	"context"
	"fmt"
	"time"
	domain_event "transcoder/internal/domain/event"
	"transcoder/internal/domain/video"
	"transcoder/internal/presentation/mappers"
	"transcoder/internal/presentation/messaging"
	pres_ports "transcoder/internal/presentation/ports"
)

type VideoUploadedHandler struct {
	service   pres_ports.VideoUploadedService
	publisher pres_ports.Publisher
}

func NewVideoUploadedHandler(service pres_ports.VideoUploadedService) *VideoUploadedHandler {
	return &VideoUploadedHandler{
		service: service,
	}
}

func (handler *VideoUploadedHandler) Handle(ctx context.Context, msg messaging.Message) error {
	dto := mappers.FromUploadedEventToDTO(msg.Payload)

	err := handler.service.Transcode(ctx, dto)
	if err != nil {
		errorEvent := video.UploadedFailedEvent{
			Base:    domain_event.Base{CreatedAt: time.Now()},
			VideoID: dto.VideoID,
		}
		handler.publisher.PublishEvents(ctx, []domain_event.Interface{&errorEvent})
		return fmt.Errorf("error while video uploading: %w", err)
	}

	return nil
}
