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

func NewVideoUploadedHandler(service pres_ports.VideoUploadedService, publisher pres_ports.Publisher) *VideoUploadedHandler {
	return &VideoUploadedHandler{
		service: service,
		publisher: publisher,
	}
}

func (handler *VideoUploadedHandler) Handle(ctx context.Context, msg messaging.Message) error {
	dto, err := mappers.FromUploadedEventToDTO(msg.Payload)
	if err != nil {
		return fmt.Errorf("error while video transcoding: %w", err)
	}

	err = handler.service.Transcode(ctx, dto)
	if err != nil {
		errorEvent := video.UploadedFailedEvent{
			Base:    domain_event.Base{CreatedAt: time.Now()},
			VideoID: dto.VideoID,
		}
		handler.publisher.PublishEvents(ctx, []domain_event.Interface{&errorEvent})
		return fmt.Errorf("error while video transcoding: %w", err)
	}

	return nil
}
