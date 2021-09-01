package cache

import (
	"context"
	"time"
)

type prefixClient struct {
	client Client
	prefix string
}

// NewPrefixClient returns a Client that adds a prefix to all keys
func NewPrefixClient(client Client, prefix string) Client {
	return &prefixClient{
		client: client,
		prefix: prefix,
	}
}

func (c *prefixClient) key(key string) string {
	return c.prefix + key
}

func (c *prefixClient) Get(ctx context.Context, key string, data interface{}) error {
	return c.client.Get(ctx, c.key(key), data)
}

func (c *prefixClient) Set(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	return c.client.Set(ctx, c.key(key), data, expiration)
}

func (c *prefixClient) Add(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	return c.client.Add(ctx, c.key(key), data, expiration)
}

func (c *prefixClient) Delete(ctx context.Context, key string) error {
	return c.client.Delete(ctx, c.key(key))
}

func (c *prefixClient) Increment(ctx context.Context, key string, delta uint64) (uint64, error) {
	return c.client.Increment(ctx, c.key(key), delta)
}

func (c *prefixClient) Decrement(ctx context.Context, key string, delta uint64) (uint64, error) {
	return c.client.Decrement(ctx, c.key(key), delta)
}
