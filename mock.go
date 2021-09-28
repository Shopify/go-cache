package cache

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ Client = &Mock{}

func (m *Mock) Get(ctx context.Context, key string, data interface{}) error {
	args := m.Called(ctx, key, data)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Mock) Set(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	args := m.Called(ctx, key, data, expiration)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Mock) Add(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	args := m.Called(ctx, key, data, expiration)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Mock) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Mock) getInt(key string) (uint64, time.Time, error) {
	args := m.Called(key)
	if args.Error(2) != nil {
		return 0, time.Time{}, args.Error(2)
	}
	return args.Get(0).(uint64), args.Get(1).(time.Time), nil
}

func (m *Mock) Increment(ctx context.Context, key string, delta uint64) (uint64, error) {
	args := m.Called(ctx, key, delta)
	if args.Error(1) != nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(uint64), nil
}

func (m *Mock) Decrement(ctx context.Context, key string, delta uint64) (uint64, error) {
	args := m.Called(ctx, key, delta)
	if args.Error(1) != nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(uint64), nil
}
