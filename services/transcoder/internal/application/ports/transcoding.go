package app_ports

import "context"

type Transcoder interface {
	Transcode(ctx context.Context, inputFilePath, outputFilePath string) error
}
