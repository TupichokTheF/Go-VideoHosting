package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenCache struct {
	client     *redis.Client
	refreshTTL time.Duration
}

func NewTokenCache(cli *redis.Client, ttl time.Duration) *TokenCache {
	return &TokenCache{
		client:     cli,
		refreshTTL: ttl,
	}
}

func (cache *TokenCache) MarkAsRevoked(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("revoked:%s", jti)

	if err := cache.client.Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf("set refresh token: %w", err)
	}

	return nil
}

func (cache *TokenCache) IsRevoked(ctx context.Context, jti string) (bool, error) {
	n, err := cache.client.Exists(ctx, "revoked:"+jti).Result()
	if err != nil {
		return false, fmt.Errorf("check revoked %s: %w", jti, err)
	}
	return n == 1, nil
}
