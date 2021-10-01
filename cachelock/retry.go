package cachelock

import (
	"sync/atomic"
	"time"
)

var (
	NoRetryStrategy = LinearBackoffStrategy(0)
)

type RetryStrategy func() RetryAttempt

// RetryAttempt allows to customise the lock retry strategy.
type RetryAttempt interface {
	// NextBackoff returns the next backoff duration.
	NextBackoff() time.Duration
}

// NoRetry acquire the lock only once.
func NoRetry() RetryAttempt {
	return linearBackoff(0)
}

type linearBackoff time.Duration

// LinearBackoff allows retries regularly with customized intervals
func LinearBackoff(backoff time.Duration) RetryAttempt {
	return linearBackoff(backoff)
}

func LinearBackoffStrategy(backoff time.Duration) RetryStrategy {
	attempt := linearBackoff(backoff) // idempotent, so memoize it
	return func() RetryAttempt {
		return attempt
	}
}

func (r linearBackoff) NextBackoff() time.Duration {
	return time.Duration(r)
}

type exponentialBackoff struct {
	cnt uint64

	min, max time.Duration
}

// ExponentialBackoff strategy is an optimization strategy with a retry time of 2**n milliseconds (n means number of times).
// You can set a minimum and maximum value, the recommended minimum value is not less than 16ms.
func ExponentialBackoff(min, max time.Duration) RetryAttempt {
	return &exponentialBackoff{min: min, max: max}
}

func ExponentialBackoffStrategy(min, max time.Duration) RetryStrategy {
	return func() RetryAttempt {
		return ExponentialBackoff(min, max)
	}
}

func (r *exponentialBackoff) NextBackoff() time.Duration {
	cnt := atomic.AddUint64(&r.cnt, 1)

	ms := 2 << 25
	if cnt < 25 {
		ms = 2 << cnt
	}

	if d := time.Duration(ms) * time.Millisecond; d < r.min {
		return r.min
	} else if r.max != 0 && d > r.max {
		return r.max
	} else {
		return d
	}
}

func AttemptBoundRetryStrategy(maxAttempts int, retryStrategy RetryStrategy) RetryStrategy {
	if maxAttempts < 1 {
		panic("max attempts must be greater than 0")
	}
	return func() RetryAttempt {
		return &attemptBoundRetryAttempt{maxAttempts: uint64(maxAttempts), retryAttempt: retryStrategy()}
	}
}

type attemptBoundRetryAttempt struct {
	attempts     uint64
	maxAttempts  uint64
	retryAttempt RetryAttempt
}

func (r *attemptBoundRetryAttempt) NextBackoff() time.Duration {
	if r.attempts >= r.maxAttempts {
		return NoRetry().NextBackoff()
	}
	atomic.AddUint64(&r.attempts, 1)
	return r.retryAttempt.NextBackoff()
}
