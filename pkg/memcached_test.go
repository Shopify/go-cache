package cache

import (
	"net"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
)

func ExampleNewMemcacheClient() {
	memcacheClient := memcache.New("localhost:11211")
	NewMemcacheClient(memcacheClient)
}

func Test_coalesceTimeoutError(t *testing.T) {
	assert.Nil(t, coalesceTimeoutError(nil))

	timeoutError := memcache.ConnectTimeoutError{Addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}}
	if err, ok := coalesceTimeoutError(&timeoutError).(net.Error); ok {
		assert.Equal(t, "connect tcp 127.0.0.1:1234: memcache: connect timeout", err.Error())
		assert.Equal(t, true, err.Timeout())
		assert.Equal(t, true, err.Temporary())
	} else {
		assert.Fail(t, "should be a net.Error")
	}
}
