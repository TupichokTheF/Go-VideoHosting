package pres_ports

import (
	"context"
	"transcoder/internal/application/dtos"
)

type VideoUploadedService interface {
	Transcode(ctx context.Context, videoData dtos.TranscodeVideo) error
}
