package cache

// Client is insired from *memcached.Client
type Client interface {
	// Get gets the item for the given key.
	// Returns nil for a cache miss.
	// The key must be at most 250 bytes in length.
	Get(key string) (*Item, error)

	// Set writes the given item, unconditionally.
	Set(key string, item *Item) error

	// Add writes the given item, if no value already exists for its
	// key. ErrNotStored is returned if that condition is not met.
	Add(key string, item *Item) error

	// Delete deletes the item with the provided key, if it exists.
	Delete(key string) error
}
