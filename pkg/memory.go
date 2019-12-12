package cache

import "sync"

type memoryClient struct {
	data sync.Map
}

// NewMemoryClient returns a Client that only stores in memory.
// Useful for stubbing tests.
// Note that it does not honour expiration.
func NewMemoryClient() Client {
	return &memoryClient{}
}

func (c *memoryClient) Get(key string) (*Item, error) {
	item, ok := c.data.Load(key)
	if !ok {
		return nil, nil
	}
	return item.(*Item), nil

}

func (c *memoryClient) Set(key string, item *Item) error {
	c.data.Store(key, item)
	return nil
}

func (c *memoryClient) Add(key string, item *Item) error {
	_, loaded := c.data.LoadOrStore(key, item)
	if loaded {
		return ErrNotStored
	}
	return nil
}

func (c *memoryClient) Delete(key string) error {
	c.data.Delete(key)
	return nil
}
