package cache

import (
	"reflect"
	"time"
)

// Client is inspired from *memcached.Client
type Client interface {
	// Get gets the item for the given key.
	// Returns nil for a cache miss.
	// The key must be at most 250 bytes in length.
	Get(key string, data interface{}) error

	// Set writes the given item, unconditionally.
	Set(key string, data interface{}, expiration time.Time) error

	// Add writes the given item, if no value already exists for its
	// key. ErrNotStored is returned if that condition is not met.
	Add(key string, data interface{}, expiration time.Time) error

	// Delete deletes the item with the provided key, if it exists.
	Delete(key string) error

	// Increment adds delta to the currently stored number
	// If the key does not exist, it will be initialized to 0 with no expiration.
	// If the value overflows, it will loop around from 0
	// For many client implementations, you need to be using a LiteralEncoding for this feature to work
	Increment(key string, delta uint64) (uint64, error)

	// Decrement subtracts delta to the currently stored number
	// If the key does not exist, it will be initialized to 0 with no expiration, which will make it overflow.
	// If the value overflows, it will loop around from MaxUint64
	// For many client implementations, you need to be using a LiteralEncoding for this feature to work
	Decrement(key string, delta uint64) (uint64, error)
}

var NeverExpire time.Time

// TtlForExpiration returns the duration until t
// Deprecated: Do not use the exported function
// nolint:golint
func TtlForExpiration(t time.Time) time.Duration {
	return ttlForExpiration(t)
}

func ttlForExpiration(t time.Time) time.Duration {
	if t.IsZero() {
		return 0
	}

	return time.Until(t)
}

func setPointerValue(ptr interface{}, value interface{}) error {
	if !isPointer(ptr) {
		return ErrNotAPointer
	}

	reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(value))
	return nil
}

func isPointer(data interface{}) bool {
	switch reflect.ValueOf(data).Kind() {
	case reflect.Ptr, reflect.Interface:
		return true
	default:
		return false
	}
}
