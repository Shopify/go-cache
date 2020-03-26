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
		var data interface{}
		err := client.Get(testKey, &data)
		require.Nil(t, data)
		require.NoError(t, err)
	})

	t.Run("set", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		for _, data := range []int{123, 124} {
			err := client.Set(testKey, data, time.Now().Add(1 * time.Second))
			require.NoError(t, err)

			var loaded int
			err = client.Get(testKey, &loaded)
			require.NoError(t, err)
			require.Equal(t, data, loaded)
		}
	})

	t.Run("add", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		data := 123

		err := client.Add(testKey, data, time.Now().Add(1 * time.Second))
		require.NoError(t, err)

		var loaded int
		err = client.Get(testKey, &loaded)
		require.NoError(t, err)
		require.Equal(t, data, loaded)

		data2 := 124
		err = client.Add(testKey, data2, time.Now().Add(1 * time.Second))
		require.EqualError(t, err, "not stored")

		var loaded2 int
		err = client.Get(testKey, &loaded2)
		require.NoError(t, err)
		require.Equal(t, data, loaded2)
	})

	t.Run("delete", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := 123

		err := client.Set(testKey, stored, time.Now().Add(1 * time.Second))
		require.NoError(t, err)

		var loaded int
		err = client.Get(testKey, &loaded)
		require.NoError(t, err)
		require.Equal(t, loaded, loaded)

		err = client.Delete(testKey)
		require.NoError(t, err)

		var loaded2 int
		err = client.Get(testKey, &loaded2)
		require.Zero(t, loaded2)
		require.NoError(t, err)
	})

	t.Run("expire", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", r.Int63())
		stored := 123

		err := client.Set(testKey, stored, time.Now().Add(1 * time.Second))
		require.NoError(t, err)

		var loaded int
		err = client.Get(testKey, &loaded)
		require.NoError(t, err)
		require.Equal(t, loaded, loaded)

		time.Sleep(1 * time.Second)

		var loaded2 int
		err = client.Get(testKey, &loaded2)
		require.Zero(t, loaded2)
		require.NoError(t, err)
	})
}
