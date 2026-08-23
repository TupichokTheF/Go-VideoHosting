package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)


type TokenCache struct {
	client *redis.Client
	refreshTTL time.Duration
}

func NewTokenCache(cli *redis.Client, ttl time.Duration) *TokenCache {
	return &TokenCache{
		client: cli,
		refreshTTL: ttl,
	}
}

func (cache *TokenCache) SetRefreshToken(ctx context.Context, refresh string, userID int) error {
	key := fmt.Sprintf("refresh:%v", userID)

	if err := cache.client.Set(ctx, key, refresh, cache.refreshTTL).Err(); err != nil {
		return  fmt.Errorf("set refresh token: %w", err)
	}

	return nil
}

func (cache *TokenCache) DeleteToken(ctx context.Context, userID int) error {
	key := fmt.Sprintf("refresh:%v", userID)

	if err := cache.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("set refresh token: %w", err)
	}

	return nil
}

func (cache *TokenCache) GetRefreshToken(ctx context.Context, userID int) (string, error) {
	key := fmt.Sprintf("refresh:%v", userID)

	value := cache.client.Get(ctx, key)
	refreshToken, err := value.Val(), value.Err() 
	if err != nil {
		return "", fmt.Errorf("set refresh token: %w", err)
	}

	return refreshToken, nil
}