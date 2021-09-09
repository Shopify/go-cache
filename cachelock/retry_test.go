package cachelock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetryStrategy(t *testing.T) {
	t.Run("no-retry", func(t *testing.T) {
		retry := NoRetry()
		require.Equal(t, time.Duration(0), retry.NextBackoff())
	})

	t.Run("linear backoff", func(t *testing.T) {
		retry := LinearBackoff(time.Second)
		require.Equal(t, time.Second, retry.NextBackoff())
		require.Equal(t, time.Second, retry.NextBackoff())
	})

	t.Run("exponential", func(t *testing.T) {
		retry := ExponentialBackoff(10*time.Millisecond, 300*time.Millisecond)
		require.Equal(t, 10*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 10*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 16*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 32*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 64*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 128*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 256*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 300*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 300*time.Millisecond, retry.NextBackoff())
		require.Equal(t, 300*time.Millisecond, retry.NextBackoff())
	})
}
