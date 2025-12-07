package lib

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLib interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, exp time.Duration) error
}

type redisLib struct {
	client *redis.Client
}

func NewRedisLib(client *redis.Client) RedisLib {
	return &redisLib{
		client: client,
	}
}

func (l *redisLib) Get(ctx context.Context, key string) (string, error) {
	return l.client.Get(ctx, key).Result()
}

func (l *redisLib) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return l.client.Set(ctx, key, value, exp).Err()
}
