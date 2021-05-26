package cache

import (
	"context"
	"time"

	"github.com/Shopify/go-encoding"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(c *redis.Client, enc encoding.ValueEncoding) Client {
	return &redisClient{client: c, encoding: enc}
}

type redisClient struct {
	client   *redis.Client
	encoding encoding.ValueEncoding
}

func (c *redisClient) Get(ctx context.Context, key string, data interface{}) error {
	cmd := c.client.Get(ctx, key)
	b, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return c.encoding.Decode(b, data)
}

func (c *redisClient) Set(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	data, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}

	cmd := c.client.Set(ctx, key, data, ttlForExpiration(expiration))
	return cmd.Err()
}

func (c *redisClient) Add(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	b, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}

	cmd := c.client.SetNX(ctx, key, b, ttlForExpiration(expiration))
	if !cmd.Val() {
		return ErrNotStored
	}
	return cmd.Err()
}

func (c *redisClient) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key)
	return err.Err()
}

func (c *redisClient) Increment(ctx context.Context, key string, delta uint64) (uint64, error) {
	cmd := c.client.IncrBy(ctx, key, int64(delta))
	val, err := cmd.Result()
	return uint64(val), err
}

func (c *redisClient) Decrement(ctx context.Context, key string, delta uint64) (uint64, error) {
	cmd := c.client.DecrBy(ctx, key, int64(delta))
	val, err := cmd.Result()
	return uint64(val), err
}
