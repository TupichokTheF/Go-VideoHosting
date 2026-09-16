package handlers

import (
	"context"
	app_ports "transcoder/internal/application/ports"
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

func (handler *VideoUploadedHandler) Handle(ctx context.Context, msg app_ports.Message) error {
	return handler.service.Transcode(ctx)
}
