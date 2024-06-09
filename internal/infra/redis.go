package infra

import (
	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client.
func NewRedisClient(redisAddr string) *redis.Client {
	return redis.NewClient(&redis.Options{
        Addr: redisAddr,
  })
}
