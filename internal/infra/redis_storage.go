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
func NewRedisStorage(addr, port string) (*RedisStorage, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr: addr+":"+port,
    })
    
    ctx := context.Background()
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }

    // Fushdb on redis
    rdb.FlushDB(ctx)

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

// UpdateTTLForKeysWithPrefix updates the TTL for keys with the given prefix in Redis.
func (r *RedisStorage) UpdateTTLForKeysWithPrefix(prefix string, ttl time.Duration) error {
    var cursor uint64
    for {
        ctx := context.Background()
        keys, nextCursor, err := r.client.Scan(ctx, cursor, prefix+"*", 60).Result()
        if err != nil {
            return err
        }

        for _, key := range keys {
            err := r.client.Expire(ctx, key, ttl).Err()
            if err != nil {
                return err
            }
        }

        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
    return nil
}

// FindByKey finds keys with the specified prefix in Redis and returns whether any keys were found.
func (r *RedisStorage) FindByKey(prefix string) (bool, time.Duration, error) {
    var cursor uint64
    for {
        ctx := context.Background()
        // SCAN command to find keys with the specified prefix
        keys, nextCursor, err := r.client.Scan(ctx, cursor, prefix+"*", 60).Result()
        if err != nil {
            return false, 0, err
        }

        for _, key := range keys {
            // GET command to get the value of the key
            _, err := r.client.Get(ctx, key).Result()
            if err != nil {
                return false, 0, err
            }

            // TTL command to get the time to live of the key
            ttl, err := r.client.TTL(ctx, key).Result()
            if err != nil {
                return false, 0, err
            }

            // Print the key and value found
            return true, ttl, nil
        }

        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
    return false, 0, nil
}
