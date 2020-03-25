package cache

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

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
	testClient(t, testMemcached(t))
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
