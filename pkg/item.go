package cache

import "time"

// Item is an item to be got or stored.
type Item struct {
	Expiration time.Time
	Data       interface{}
}

func (item *Item) Duration() time.Duration {
	if item.Expiration.IsZero() {
		return 0
	}

	return time.Until(item.Expiration)
}
