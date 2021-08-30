package cachelock

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/go-cache/v2"
)

const tokenLength = 16

type cacheLocker struct {
	client        cache.Client
	expiration    time.Duration
	retryStrategy func() RetryAttempt
}

func New(client cache.Client, expiration time.Duration, retryStrategy RetryStrategy) Locker {
	return &cacheLocker{
		client:        client,
		expiration:    expiration,
		retryStrategy: retryStrategy,
	}
}

func (l *cacheLocker) Acquire(ctx context.Context, key string) (Lock, error) {
	retry := l.retryStrategy()
	token, err := randomToken(tokenLength)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	expiration := time.Now().Add(l.expiration)
	ctx, cancel := context.WithDeadline(ctx, expiration)
	defer cancel()

	var timer *time.Timer
	for {
		err := l.client.Add(ctx, key, token, expiration)
		if err == nil {
			return &cacheLock{client: l.client, key: key, token: token}, nil
		} else if !errors.Is(err, cache.ErrNotStored) {
			return nil, fmt.Errorf("locking: %w", err)
		}

		backoff := retry.NextBackoff()
		if backoff < 1 {
			return nil, ErrNotAcquired
		}

		if timer == nil {
			timer = time.NewTimer(backoff)
			defer timer.Stop()
		} else {
			timer.Reset(backoff)
		}

		select {
		case <-ctx.Done():
			return nil, ErrNotAcquired
		case <-timer.C:
		}
	}
}

func randomToken(length int) (string, error) {
	tok := make([]byte, length)

	if _, err := rand.Read(tok); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tok), nil

}

type cacheLock struct {
	client cache.Client
	key    string
	token  string
}

func (r *cacheLock) Release(ctx context.Context) error {
	var stored string
	err := r.client.Get(ctx, r.key, &stored)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			err = ErrNotReleased
		}
		return err
	}

	if stored != r.token {
		return ErrNotReleased
	}

	// This implementation checks for the lock being held and still having the same token,
	// but it's _possible_ there's a race condition here and the lock will be tampered with before we delete it here.
	// However, that should be extremely unlikely, since only one thread should be able to lock or attempt unlocking at the same time.
	// Therefore, this particular implementation takes the shortcut of not dealing with that race condition.

	return r.client.Delete(ctx, r.key)
}
