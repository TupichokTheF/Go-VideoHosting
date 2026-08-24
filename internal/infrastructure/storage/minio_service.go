package storage

import (
	"context"
	"fmt"
	"net/url"
	"time"

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

func (service *MinioService) PresignedURLCreate(ctx context.Context, key string, ttl time.Duration) (string, error) {
	url, err := service.client.PresignedPutObject(ctx, service.bucket, key, ttl)
	if err != nil {
		return "", fmt.Errorf("presign put %q: %w", key, err)
	}

	return url.String(), nil
}

func (service *MinioService) PresignedURLGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	url, err := service.client.PresignedGetObject(ctx, service.bucket, key, ttl, url.Values{})
	if err != nil {
		return "", fmt.Errorf("presign get %q: %w", key, err)
	}

	return url.String(), nil
}
