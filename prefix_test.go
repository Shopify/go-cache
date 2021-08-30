package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixClient(t *testing.T) {
	c := NewMemoryClient()
	pc := NewPrefixClient(c, "test.")

	t.Run("prefix", func(t *testing.T) {
		ctx := context.Background()
		err := pc.Set(ctx, "foo", "bar", NeverExpire)
		require.NoError(t, err)

		var out string
		err = c.Get(ctx, "test.foo", &out)
		require.NoError(t, err)
		require.Equal(t, "bar", out)
	})

	testClient(t, pc, nil)
}
