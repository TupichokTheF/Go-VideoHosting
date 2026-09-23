package app_ports

import (
	"context"
	"errors"
)

type Storage interface {
	PresignedURLCreate(ctx context.Context, key string) (string, error)
	PresignedURLGet(ctx context.Context, key string) (string, error)
	Stat(ctx context.Context, key string) (int64, error)
}

var (
	ErrObjectNotFound = errors.New("Object was not found")
)