package cache

import "github.com/pkg/errors"

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrNotStored = errors.New("not stored")
	ErrNotAPointer = errors.New("argument to Get() must be a pointer")
)
