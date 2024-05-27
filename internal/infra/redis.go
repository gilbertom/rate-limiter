package infra

import (
	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client.
func NewRedisClient(redisAddr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
  })
	return rdb
}