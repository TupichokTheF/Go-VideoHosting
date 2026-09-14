package app_ports

import (
	"context"
	"time"
)

type TokenCache interface {
	MarkAsRevoked(ctx context.Context, jti string, ttl time.Duration) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}
