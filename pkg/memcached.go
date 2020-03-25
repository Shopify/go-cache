package cache

import (
	"math"
	"net"

	"github.com/bradfitz/gomemcache/memcache"
)

var _ Client = &memcacheClient{}

func NewMemcacheClient(c *memcache.Client) *memcacheClient {
	return &memcacheClient{client: c, encoding: GobEncoding}
}

type memcacheClient struct {
	client   *memcache.Client
	encoding Encoding
}

func (c *memcacheClient) Get(key string) (*Item, error) {
	mItem, err := c.client.Get(key)
	if err != nil {
		// Abstract the memcache-specific error
		if err == memcache.ErrCacheMiss {
			err = nil
		}
		return nil, coalesceTimeoutError(err)
	}

	return c.decodeItem(mItem)
}

func (c *memcacheClient) Set(key string, item *Item) error {
	mItem, err := c.encodeItem(key, item)
	if err != nil {
		return err
	}
	return coalesceTimeoutError(c.client.Set(mItem))
}

func (c *memcacheClient) Add(key string, item *Item) error {
	mItem, err := c.encodeItem(key, item)
	if err != nil {
		return err
	}
	err = c.client.Add(mItem)

	if err == memcache.ErrNotStored {
		// Abstract the memcache-specific error
		return ErrNotStored
	}
	return coalesceTimeoutError(err)
}

func (c *memcacheClient) Delete(key string) error {
	err := c.client.Delete(key)
	if err == memcache.ErrCacheMiss {
		// Deleting a missing entry is not an actual issue.
		return nil
	}
	return coalesceTimeoutError(err)
}

func (c *memcacheClient) decodeItem(mItem *memcache.Item) (*Item, error) {
	return c.encoding.Decode(mItem.Value)
}

func (c *memcacheClient) encodeItem(key string, item *Item) (*memcache.Item, error) {
	encoded, err := c.encoding.Encode(item)
	if err != nil {
		return nil, err
	}

	mItem := &memcache.Item{
		Value: encoded,
		Key:   key,
	}
	if item.Duration() != 0 {
		mItem.Expiration = int32(math.Max(item.Duration().Seconds(), 1))
	}

	return mItem, nil
}

type connectTimeoutError struct{}

func (connectTimeoutError) Error() string   { return "memcache: connect timeout" }
func (connectTimeoutError) Timeout() bool   { return true }
func (connectTimeoutError) Temporary() bool { return true }

func coalesceTimeoutError(err error) error {
	// For some reason, gomemcache decided to replace the standard net.Error.
	// Coalesce into a generic net.Error so that client don't have to deal with memcache-specific errors.
	switch typed := err.(type) {
	case *memcache.ConnectTimeoutError:
		return &net.OpError{
			Err:  &connectTimeoutError{},
			Addr: typed.Addr,
			Net:  typed.Addr.Network(),
			Op:   "connect",
		}
	default:
		// This also work if err is nil
		return err
	}
}
