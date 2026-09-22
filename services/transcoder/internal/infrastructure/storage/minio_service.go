package storage

import (
	"context"
	"fmt"
	app_ports "transcoder/internal/application/ports"

	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(client *minio.Client, bucket string) *MinioService {
	return &MinioService{
		client: client,
		bucket: bucket,
	}
}

func (service *MinioService) Download(ctx context.Context, key string, dstDir string) error {
	err := service.client.FGetObject(ctx, service.bucket, key, dstDir, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return fmt.Errorf("download video: %w", app_ports.ErrObjectNotFound)
		}

		return fmt.Errorf("download video: %w", err)
	}

	return nil
}

func (service *MinioService) IsExist(ctx context.Context, key string) (bool, error) {
	_, err := service.client.StatObject(ctx, service.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, fmt.Errorf("is exist: %w", app_ports.ErrObjectNotFound)
		}

		return false, fmt.Errorf("is exist: %w", err)
	}

	return true, nil
}

func (service *MinioService) Upload(ctx context.Context, key, srcDir string) error {
	_, err := service.client.FPutObject(ctx, service.bucket, key, srcDir, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("upload video: %w", err)
	}

	return nil
}
