package cache

import (
	"os"
	"testing"

	"github.com/Shopify/go-encoding"
	"github.com/go-redis/redis/v8"
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

func testRedis(t *testing.T) *redis.Client {
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
	return redis.NewClient(opts)
}

func Test_redisClient(t *testing.T) {
	client := testRedis(t)
	encodings := map[string]encoding.ValueEncoding{
		"gob":          gobEncoding,
		"json":         encoding.JSONEncoding,
		"literal+gob":  encoding.NewLiteralEncoding(gobEncoding),
		"literal+json": encoding.NewLiteralEncoding(encoding.JSONEncoding),
	}
	for name, enc := range encodings {
		t.Run(name, func(t *testing.T) {
			testClient(t, NewRedisClient(client, enc), enc)
		})
	}

}
