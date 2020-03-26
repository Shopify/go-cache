package cache

import (
	"reflect"
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

func (c *memoryClient) Get(key string, data interface{}) error {
	if item, ok := c.data.Load(key); ok {
		mItem := item.(memoryData)
		if mItem.expiration.IsZero() || mItem.expiration.After(time.Now()) {
			reflect.ValueOf(data).Elem().Set(reflect.ValueOf(mItem.data))
			return nil
		}
	}
	return ErrCacheMiss
}

func (c *memoryClient) Set(key string, data interface{}, expiration time.Time) error {
	c.data.Store(key, memoryData{
		data:       data,
		expiration: expiration,
	})
	return nil
}

func (c *memoryClient) Add(key string, data interface{}, expiration time.Time) error {
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

func (c *memoryClient) Delete(key string) error {
	c.data.Delete(key)
	return nil
}
