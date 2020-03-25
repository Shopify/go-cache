package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, client Client) {
	r := rand.NewSource(time.Now().UnixNano())

	t.Run("get not existing", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		item, err := client.Get(testKey)
		require.Nil(t, item)
		require.NoError(t, err)
	})

	t.Run("set", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		for _, data := range []int{123, 124} {
			stored := &Item{
				Expiration: time.Now().Add(1 * time.Second),
				Data:       data,
			}

			err := client.Set(testKey, stored)
			require.NoError(t, err)

			loaded, err := client.Get(testKey)
			require.NoError(t, err)
			require.Equal(t, stored.Data, loaded.Data)
		}
	})

	t.Run("add", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       123,
		}

		err := client.Add(testKey, stored)
		require.NoError(t, err)

		loaded, err := client.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)

		stored2 := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       124,
		}
		err = client.Add(testKey, stored2)
		require.EqualError(t, err, "not stored")

		loaded, err = client.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)
	})

	t.Run("delete", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       123,
		}

		err := client.Set(testKey, stored)
		require.NoError(t, err)

		loaded, err := client.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)

		err = client.Delete(testKey)
		require.NoError(t, err)

		loaded, err = client.Get(testKey)
		require.NoError(t, err)
		require.Nil(t, loaded)
	})

	t.Run("expire", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       123,
		}

		err := client.Set(testKey, stored)
		require.NoError(t, err)

		loaded, err := client.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)

		time.Sleep(1 * time.Second)

		loaded, err = client.Get(testKey)
		require.NoError(t, err)
		require.Nil(t, loaded)
	})
}
