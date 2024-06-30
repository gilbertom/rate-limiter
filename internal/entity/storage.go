package entity

import (
	"time"
)

// Storage represents a storage interface.
type Storage interface {
	Incr(key string) (int64, error)
	Expire(key string, seconds int) error
	UpdateTTLForKeysWithPrefix(prefix string, ttl time.Duration) error
	FindByKey(prefix string) (bool, time.Duration, error)
}
