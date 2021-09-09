package cachelock

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotAcquired = errors.New("lock not acquired")
	ErrNotReleased = errors.New("lock not released")
)

const DefaultLockExpiration = time.Minute * 5

// Locker is a service that provides a lock for a resource.
type Locker interface {
	Acquire(ctx context.Context, key string) (Lock, error)
}

// Lock is a handle on an acquired lock.
type Lock interface {
	Release(context.Context) error
}
