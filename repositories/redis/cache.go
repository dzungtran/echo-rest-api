package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheRepository represent the redis repositories
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type redisRepository struct {
	client *redis.Client
}

// NewCacheRepository will create an object that represent the CacheRepository interface
func NewCacheRepository(client *redis.Client) CacheRepository {
	return &redisRepository{
		client: client,
	}
}

// Set attaches the redis repository and set the data
func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return r.client.Set(ctx, key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
