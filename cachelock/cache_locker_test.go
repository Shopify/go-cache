package cachelock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Shopify/go-cache/v2"
)

func Test_Lock_Release(t *testing.T) {
	ctx := context.Background()
	r := newTestCacheLock()

	lock, err := r.Acquire(ctx, "foo")
	require.NoError(t, err)

	err = lock.Release(ctx)
	require.NoError(t, err)

	lock, err = r.Acquire(ctx, "foo")
	require.NoError(t, err)

	err = lock.Release(ctx)
	require.NoError(t, err)
}

func Test_Lock_Different_Key(t *testing.T) {
	ctx := context.Background()
	r := newTestCacheLock()

	_, err := r.Acquire(ctx, "foo")
	require.NoError(t, err)

	_, err = r.Acquire(ctx, "bar")
	require.NoError(t, err)
}

func Test_Release(t *testing.T) {
	ctx := context.Background()
	r := newTestCacheLock()

	lock, err := r.Acquire(ctx, "foo")
	require.NoError(t, err)

	err = lock.Release(ctx)
	require.NoError(t, err)

	err = lock.Release(ctx)
	require.Equal(t, ErrNotReleased, err)
}

func Test_Lock_Concurrent_With_Retry(t *testing.T) {
	ctx := context.Background()
	r := newTestCacheLock()

	_, err := r.Acquire(ctx, "foo")
	require.NoError(t, err)

	_, err = r.Acquire(ctx, "foo")
	require.Equal(t, ErrNotAcquired, err)
}

func newTestCacheLock() Locker {
	return New(cache.NewMemoryClient(), DefaultLockExpiration, NoRetryStrategy)
}
