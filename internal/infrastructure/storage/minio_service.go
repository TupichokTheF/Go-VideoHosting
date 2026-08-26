package storage

import (
	"context"
	"fmt"
	"net/url"
	"project/internal/domain/video"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	client *minio.Client
	bucket string
	ttl    time.Duration
}

func NewMinioService(client *minio.Client, bucket string, ttl time.Duration) *MinioService {
	return &MinioService{
		client: client,
		bucket: bucket,
		ttl:    ttl,
	}
}

func (service *MinioService) PresignedURLCreate(ctx context.Context, key string) (string, error) {
	url, err := service.client.PresignedPutObject(ctx, service.bucket, key, service.ttl)
	if err != nil {
		return "", fmt.Errorf("presign put %q: %w", key, err)
	}

	return url.String(), nil
}

func (service *MinioService) PresignedURLGet(ctx context.Context, key string) (string, error) {
	url, err := service.client.PresignedGetObject(ctx, service.bucket, key, service.ttl, url.Values{})
	if err != nil {
		return "", fmt.Errorf("presign get %q: %w", key, err)
	}

	return url.String(), nil
}

func (s *MinioService) Stat(ctx context.Context, key string) (int64, error) {
	info, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return 0, fmt.Errorf("stat %q: %w", key, video.ErrNotFound)
		}

		return 0, fmt.Errorf("stat %q: %w", key, err)
	}
	return info.Size, nil
}
