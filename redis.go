package cache

import (
	"context"
	"time"

	"github.com/Shopify/go-encoding"
	"github.com/go-redis/redis/v8"
)

var _ Client = &redisClient{}

func NewRedisClient(c *redis.Client, enc encoding.ValueEncoding) *redisClient {
	return &redisClient{client: c, encoding: enc}
}

type redisClient struct {
	client   *redis.Client
	encoding encoding.ValueEncoding
}

func (c *redisClient) Get(key string, data interface{}) error {
	cmd := c.client.Get(context.Background(), key)
	b, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return c.encoding.Decode(b, data)
}

func (c *redisClient) Set(key string, data interface{}, expiration time.Time) error {
	data, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}

	cmd := c.client.Set(context.Background(), key, data, TtlForExpiration(expiration))
	return cmd.Err()
}

func (c *redisClient) Add(key string, data interface{}, expiration time.Time) error {
	b, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}

	cmd := c.client.SetNX(context.Background(), key, b, TtlForExpiration(expiration))
	if !cmd.Val() {
		return ErrNotStored
	}
	return cmd.Err()
}

func (c *redisClient) Delete(key string) error {
	err := c.client.Del(context.Background(), key)
	return err.Err()
}

func (c *redisClient) Increment(key string, delta uint64) (uint64, error) {
	cmd := c.client.IncrBy(context.Background(), key, int64(delta))
	val, err := cmd.Result()
	return uint64(val), err
}

func (c *redisClient) Decrement(key string, delta uint64) (uint64, error) {
	cmd := c.client.DecrBy(context.Background(), key, int64(delta))
	val, err := cmd.Result()
	return uint64(val), err
}
