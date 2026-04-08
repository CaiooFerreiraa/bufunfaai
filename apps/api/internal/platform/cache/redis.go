package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
)

func Connect(ctx context.Context, cfg config.Config) (*redis.Client, error) {
	options, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
