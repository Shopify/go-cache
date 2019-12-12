package cache

import "time"

// Item is an item to be got or stored.
type Item struct {
	Expiration time.Time
	Data       interface{}
}
