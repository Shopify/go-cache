package cachelock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Locker = (*MockLocker)(nil)
var _ Lock = (*MockLock)(nil)

func ExpectAcquireAndRelease(locker *MockLocker, key string) *MockLock {
	lock := NewMockLock()
	locker.On("Acquire", mock.Anything, key).Return(lock, nil)
	lock.On("Release", mock.Anything).Return(nil)
	return lock
}

func NewMockLocker() *MockLocker {
	return &MockLocker{}
}

type MockLocker struct {
	mock.Mock
}

func (m *MockLocker) Acquire(ctx context.Context, key string) (Lock, error) {
	args := m.Called(ctx, key)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(Lock), args.Error(1)
}

func NewMockLock() *MockLock {
	return &MockLock{}
}

type MockLock struct {
	mock.Mock
}

func (m *MockLock) Release(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
