package cache

import "github.com/redis/go-redis/v9"

type ClientAdapter struct {
	client *redis.Client
}

func NewClientAdapter(client *redis.Client) *ClientAdapter {
	return &ClientAdapter{client: client}
}

func (adapter *ClientAdapter) Close() error {
	if adapter == nil || adapter.client == nil {
		return nil
	}

	return adapter.client.Close()
}
