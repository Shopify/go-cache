package cachelock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMockForKey(t *testing.T) {
	locker := NewMockLocker()
	ctx := context.Background()

	lock := ExpectAcquireAndRelease(locker, "foo")

	acquired, err := locker.Acquire(ctx, "foo")
	require.NoError(t, err)
	require.Equal(t, lock, acquired)

	err = acquired.Release(ctx)
	require.NoError(t, err)

	locker.AssertExpectations(t)
	lock.AssertExpectations(t)
}
