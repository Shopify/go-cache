package cache

import "github.com/pkg/errors"

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrNotStored = errors.New("not stored")
)
