package cache

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/require"
)

func ExampleNewRedisClient() {
	opts, err := redis.ParseURL("")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opts)
	NewRedisClient(client, DefaultEncoding)
}

func testRedis(t *testing.T) Client {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	url := os.Getenv("REDIS_URL")
	if len(url) == 0 {
		t.Skip("redis client not configured")
		return nil
	}

	opts, err := redis.ParseURL(url)
	require.NoError(t, err)

	return NewRedisClient(redis.NewClient(opts), DefaultEncoding)
}

func Test_redisClient(t *testing.T) {
	testClient(t, testRedis(t))
}
