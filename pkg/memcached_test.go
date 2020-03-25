package cache

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/require"
)

func ExampleNewMemcacheClient() {
	memcacheClient := memcache.New("localhost:11211")
	NewMemcacheClient(memcacheClient)
}

const defaultMemcachedPort = 11211

func testMemcached(t *testing.T) Client {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	serversStr := os.Getenv("MEMCACHED_SERVERS")
	if len(serversStr) == 0 {
		t.Skip("memcache client not configured")
		return nil
	}

	servers := strings.Split(serversStr, ",")
	for i, server := range servers {
		if !strings.ContainsRune(server, ':') {
			servers[i] = fmt.Sprintf("%s:%d", server, defaultMemcachedPort)
		}
	}

	return NewMemcacheClient(memcache.New(servers...))
}

func Test_memcacheClient(t *testing.T) {
	cache := testMemcached(t)

	t.Run("get not existing", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", rand.Int())
		item, err := cache.Get(testKey)
		require.Nil(t, item)
		require.NoError(t, err)
	})

	t.Run("set", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", rand.Int())
		for _, data := range []int{123, 124} {
			stored := &Item{
				Expiration: time.Now().Add(1 * time.Second),
				Data:       data,
			}

			err := cache.Set(testKey, stored)
			require.NoError(t, err)

			loaded, err := cache.Get(testKey)
			require.NoError(t, err)
			require.Equal(t, stored.Data, loaded.Data)
		}
	})

	t.Run("add", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", rand.Int())
		stored := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       123,
		}

		err := cache.Add(testKey, stored)
		require.NoError(t, err)

		loaded, err := cache.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)

		stored2 := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       124,
		}
		err = cache.Add(testKey, stored2)
		require.EqualError(t, err, "not stored")

		loaded, err = cache.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)
	})

	t.Run("delete", func(t *testing.T) {
		testKey := fmt.Sprintf("go-cache-test-%d", rand.Int())
		stored := &Item{
			Expiration: time.Now().Add(1 * time.Second),
			Data:       123,
		}

		err := cache.Set(testKey, stored)
		require.NoError(t, err)

		loaded, err := cache.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, stored.Data, loaded.Data)

		err = cache.Delete(testKey)
		require.NoError(t, err)

		loaded, err = cache.Get(testKey)
		require.NoError(t, err)
		require.Nil(t, loaded)
	})
}

func Test_coalesceTimeoutError(t *testing.T) {
	require.Nil(t, coalesceTimeoutError(nil))

	timeoutError := memcache.ConnectTimeoutError{Addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}}
	if err, ok := coalesceTimeoutError(&timeoutError).(net.Error); ok {
		require.Equal(t, "connect tcp 127.0.0.1:1234: memcache: connect timeout", err.Error())
		require.Equal(t, true, err.Timeout())
		require.Equal(t, true, err.Temporary())
	} else {
		require.Fail(t, "should be a net.Error")
	}
}
