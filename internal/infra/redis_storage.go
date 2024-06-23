package infra

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStorage is a struct that represents the Redis storage.
type RedisStorage struct {
    client *redis.Client
    ctx    context.Context
}

// NewRedisStorage creates a new instance of RedisStorage.
func NewRedisStorage(addr string) (*RedisStorage, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    
    ctx := context.Background()
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }

    return &RedisStorage{
        client: rdb,
        ctx:    ctx,
    }, nil
}

// Incr increments the value of a key in Redis.
func (r *RedisStorage) Incr(key string) (int64, error) {
    return r.client.Incr(r.ctx, key).Result()
}

// Expire sets the expiration time for a key in Redis.
func (r *RedisStorage) Expire(key string, seconds int) error {
    return r.client.Expire(r.ctx, key, time.Duration(seconds)*time.Second).Err()
}
