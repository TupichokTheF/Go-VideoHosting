package handlers

import (
	"context"
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
	return handler.service.Transcode(ctx)
}
