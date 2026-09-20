package app_ports

import "context"

type Transcoder interface {
	Transcode(ctx context.Context, inputFilePath, outputFilePath string) error
}

type Storage interface {
	Download(ctx context.Context, key string, dstDir string) error
	IsExist(ctx context.Context, key string) (bool, error)
	Upload(ctx context.Context, key string, srcDir string) error
}
