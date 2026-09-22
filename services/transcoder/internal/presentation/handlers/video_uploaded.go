package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"
	app_ports "transcoder/internal/application/ports"
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
		service:   service,
		publisher: publisher,
	}
}

func (handler *VideoUploadedHandler) Handle(ctx context.Context, msg messaging.Message) error {
	dto, err := mappers.FromUploadedEventToDTO(msg.Payload)
	if err != nil {
		return fmt.Errorf("error while video transcoding: %w", err)
	}

	err = handler.service.Transcode(ctx, dto)
	switch {
	case err == nil:
		return nil

	case errors.Is(err, app_ports.ErrUnprocessiable):
		errorEvent := video.UploadedFailedEvent{
			Base:    domain_event.Base{CreatedAt: time.Now()},
			VideoID: dto.VideoID,
		}
		if err := handler.publisher.PublishEvents(ctx, []domain_event.Interface{&errorEvent}); err != nil {
			return fmt.Errorf("video transcoding handler: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("video transcoding handler: %w", err)
	}
}
