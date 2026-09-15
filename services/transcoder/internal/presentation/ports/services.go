package pres_ports

import "context"

type VideoUploadedService interface {
	Transcode(ctx context.Context) error
}
