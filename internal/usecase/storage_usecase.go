package usecases

import "github.com/gilbertom/go-rate-limiter/internal/entity"

// StorageUseCase represents a use case for storage.
type StorageUseCase struct {
    storage entity.Storage
}

// NewStorageUseCase creates a new instance of StorageUseCase.
func NewStorageUseCase(storage entity.Storage) *StorageUseCase {
    return &StorageUseCase{storage: storage}
}

// Incr increments the value for the given key in the storage.
func (uc *StorageUseCase) Incr(key string) (int64, error) {
    return uc.storage.Incr(key)
}

// Expire sets an expiration time in seconds for the given key in the storage.
func (uc *StorageUseCase) Expire(key string, seconds int) error {
    return uc.storage.Expire(key, seconds)
}