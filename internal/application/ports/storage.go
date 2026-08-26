package infra_ports

import (
	"context"
)

type Storage interface {
	PresignedURLCreate(ctx context.Context, key string) (string, error)
	PresignedURLGet(ctx context.Context, key string) (string, error)
}
