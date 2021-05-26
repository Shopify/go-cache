package cache

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, client Client, encoding Encoding) {
	r := rand.NewSource(time.Now().UnixNano())
	ctx := context.Background()

	t.Run("get not existing", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		var data interface{}
		err := client.Get(ctx, testKey, &data)
		require.Zero(t, data)
		require.EqualError(t, err, "cache miss")
	})

	t.Run("set", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		for _, data := range []int{123, 124} {
			err := client.Set(ctx, testKey, data, time.Now().Add(1*time.Second))
			require.NoError(t, err)

			var loaded int
			err = client.Get(ctx, testKey, &loaded)
			require.NoError(t, err)
			require.Equal(t, data, loaded)
		}
	})

	t.Run("add", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		data := 123

		err := client.Add(ctx, testKey, data, time.Now().Add(1*time.Second))
		require.NoError(t, err)

		var loaded int
		err = client.Get(ctx, testKey, &loaded)
		require.NoError(t, err)
		require.Equal(t, data, loaded)

		data2 := 124
		err = client.Add(ctx, testKey, data2, time.Now().Add(1*time.Second))
		require.EqualError(t, err, "not stored")

		var loaded2 int
		err = client.Get(ctx, testKey, &loaded2)
		require.NoError(t, err)
		require.Equal(t, data, loaded2)
	})

	t.Run("delete", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := 123

		err := client.Set(ctx, testKey, stored, time.Now().Add(1*time.Second))
		require.NoError(t, err)

		var loaded int
		err = client.Get(ctx, testKey, &loaded)
		require.NoError(t, err)
		require.Equal(t, loaded, loaded)

		err = client.Delete(ctx, testKey)
		require.NoError(t, err)

		var loaded2 int
		err = client.Get(ctx, testKey, &loaded2)
		require.Zero(t, loaded2)
		require.EqualError(t, err, "cache miss")
	})

	t.Run("expire", func(t *testing.T) {
		if encoding != GobEncoding {
			t.Skip("only run expire test once")
		}

		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := 123

		err := client.Set(ctx, testKey, stored, time.Now().Add(1*time.Second))
		require.NoError(t, err)

		var loaded int
		err = client.Get(ctx, testKey, &loaded)
		require.NoError(t, err)
		require.Equal(t, loaded, loaded)

		time.Sleep(1500 * time.Millisecond)

		var loaded2 int
		err = client.Get(ctx, testKey, &loaded2)
		require.Zero(t, loaded2)
		require.EqualError(t, err, "cache miss")
	})

	t.Run("incr", func(t *testing.T) {
		if encoding == GobEncoding {
			t.Skip("gob encoding does not support increment")
		}
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())

		newVal, err := client.Increment(ctx, testKey, 123)
		require.NoError(t, err)
		require.Equal(t, uint64(123), newVal)

		newVal, err = client.Increment(ctx, testKey, 10)
		require.NoError(t, err)
		require.Equal(t, uint64(133), newVal)
	})

	t.Run("incr overflow", func(t *testing.T) {
		if encoding == GobEncoding {
			t.Skip("gob encoding does not support increment")
		}
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())

		newVal, err := client.Increment(ctx, testKey, uint64(math.MaxUint64))
		require.NoError(t, err)
		require.Equal(t, uint64(math.MaxUint64), newVal)

		newVal, err = client.Increment(ctx, testKey, 10)
		require.NoError(t, err)
		require.Equal(t, uint64(9), newVal)
	})

	t.Run("decr", func(t *testing.T) {
		if encoding == GobEncoding {
			t.Skip("gob encoding does not support decrement")
		}
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())

		err := client.Set(ctx, testKey, uint64(123), time.Now().Add(1*time.Second))
		require.NoError(t, err)

		newVal, err := client.Decrement(ctx, testKey, 10)
		require.NoError(t, err)
		require.Equal(t, uint64(113), newVal)
	})

	t.Run("decr overflow", func(t *testing.T) {
		if encoding == GobEncoding {
			t.Skip("gob encoding does not support decrement")
		}
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())

		newVal, err := client.Decrement(ctx, testKey, 10)
		require.NoError(t, err)
		require.Equal(t, math.MaxUint64-uint64(9), newVal)
	})
}
