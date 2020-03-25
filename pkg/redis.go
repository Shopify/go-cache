package cache

import (
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

func (c *redisClient) Get(key string) (*Item, error) {
	cmd := c.client.Get(key)
	data, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	item, err := c.encoding.Decode(data)
	return item, err
}

func (c *redisClient) Set(key string, item *Item) error {
	data, err := c.encoding.Encode(item)
	if err != nil {
		return err
	}

	cmd := c.client.Set(key, data, item.Duration())
	return cmd.Err()
}

func (c *redisClient) Add(key string, item *Item) error {
	data, err := c.encoding.Encode(item)
	if err != nil {
		return err
	}

	cmd := c.client.SetNX(key, data, item.Duration())
	if !cmd.Val() {
		return ErrNotStored
	}
	return cmd.Err()
}

func (c *redisClient) Delete(key string) error {
	err := c.client.Del(key)
	return err.Err()
}
