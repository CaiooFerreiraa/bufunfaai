package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type ReadinessProbe struct {
	client *redis.Client
}

func NewReadinessProbe(client *redis.Client) *ReadinessProbe {
	return &ReadinessProbe{client: client}
}

func (probe *ReadinessProbe) Name() string {
	return "redis"
}

func (probe *ReadinessProbe) Check(ctx context.Context) error {
	return probe.client.Ping(ctx).Err()
}
