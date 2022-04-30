package datastore

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient will create new redis instance
func NewRedisClient(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return rdb
}
