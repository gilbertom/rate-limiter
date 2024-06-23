package entity

// Storage represents a storage interface.
type Storage interface {
	Incr(key string) (int64, error)
	Expire(key string, seconds int) error
}