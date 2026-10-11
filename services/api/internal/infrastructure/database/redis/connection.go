package app_redis

import (
	"context"
	"log/slog"
	"project/internal/core"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *core.RedisConfig) (*redis.Client, error) {
	opts := redis.Options{
		Addr:       cfg.Address(),
		MaxRetries: 5,
	}

	client := redis.NewClient(&opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func CloseClient(client *redis.Client) {
	if err := client.Close(); err != nil {
		slog.Error("error while closing redis connection", "error", err)
	}
}
