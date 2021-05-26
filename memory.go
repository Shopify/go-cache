package cache

import (
	"context"
	"sync"
	"time"
)

type memoryData struct {
	data       interface{}
	expiration time.Time
}

type memoryClient struct {
	data sync.Map
}

// NewMemoryClient returns a Client that only stores in memory.
// Useful for stubbing tests.
func NewMemoryClient() Client {
	return &memoryClient{}
}

func (c *memoryClient) Get(ctx context.Context, key string, data interface{}) error {
	if item, ok := c.data.Load(key); ok {
		mItem := item.(memoryData)
		if mItem.expiration.IsZero() || mItem.expiration.After(time.Now()) {
			return setPointerValue(data, mItem.data)
		}
	}
	return ErrCacheMiss
}

func (c *memoryClient) Set(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	c.data.Store(key, memoryData{
		data:       data,
		expiration: expiration,
	})
	return nil
}

func (c *memoryClient) Add(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	_, loaded := c.data.LoadOrStore(key, memoryData{
		data:       data,
		expiration: expiration,
	})
	if loaded {
		// TODO: handle when the conflicting data is expired
		return ErrNotStored
	}
	return nil
}

func (c *memoryClient) Delete(ctx context.Context, key string) error {
	c.data.Delete(key)
	return nil
}

func (c *memoryClient) getInt(key string) (uint64, time.Time, error) {
	var curr uint64
	var expiration time.Time
	if item, ok := c.data.Load(key); ok {
		mItem := item.(memoryData)
		if curr, ok = mItem.data.(uint64); !ok {
			return 0, expiration, ErrNotANumber
		}
		expiration = mItem.expiration
	}
	return curr, expiration, nil
}

func (c *memoryClient) Increment(ctx context.Context, key string, delta uint64) (uint64, error) {
	// TODO: definitely not thread-safe
	curr, expiration, err := c.getInt(key)
	if err != nil {
		return 0, err
	}

	curr += delta
	err = c.Set(context.Background(), key, curr, expiration)
	return curr, err
}

func (c *memoryClient) Decrement(ctx context.Context, key string, delta uint64) (uint64, error) {
	// TODO: definitely not thread-safe
	curr, expiration, err := c.getInt(key)
	if err != nil {
		return 0, err
	}

	curr -= delta
	err = c.Set(context.Background(), key, curr, expiration)
	return curr, err
}
