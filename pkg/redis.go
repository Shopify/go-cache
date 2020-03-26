package cache

import (
	"time"

	"github.com/go-redis/redis/v7"
)

var _ Client = &redisClient{}

func NewRedisClient(c *redis.Client, encoding Encoding) *redisClient {
	return &redisClient{client: c, encoding: encoding}
}

type redisClient struct {
	client   *redis.Client
	encoding Encoding
}

func (c *redisClient) Get(key string, data interface{}) error {
	cmd := c.client.Get(key)
	b, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
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

	cmd := c.client.Set(key, data, TtlForExpiration(expiration))
	return cmd.Err()
}

func (c *redisClient) Add(key string, data interface{}, expiration time.Time) error {
	b, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}

	cmd := c.client.SetNX(key, b, TtlForExpiration(expiration))
	if !cmd.Val() {
		return ErrNotStored
	}
	return cmd.Err()
}

func (c *redisClient) Delete(key string) error {
	err := c.client.Del(key)
	return err.Err()
}
