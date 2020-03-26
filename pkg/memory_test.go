package cache

import (
	"testing"
)

func TestMemoryClient(t *testing.T) {
	testClient(t, NewMemoryClient(), nil)
}
