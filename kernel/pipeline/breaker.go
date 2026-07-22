// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type CircuitBreaker struct {
	threshold   int
	cooldown    time.Duration
	failures    int
	lastFailure time.Time
	open        bool
	mu          sync.Mutex
}

func (cb *CircuitBreaker) Middleware() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			cb.mu.Lock()
			if cb.open {
				if time.Since(cb.lastFailure) > cb.cooldown {
					cb.open = false
					cb.failures = 0
				} else {
					cb.mu.Unlock()
					return nil, kernel.ErrCircuitOpen
				}
			}
			cb.mu.Unlock()

			resp, err := next(ctx, req)
			if err != nil {
				cb.mu.Lock()
				cb.failures++
				cb.lastFailure = time.Now()
				if cb.failures >= cb.threshold {
					cb.open = true
				}
				cb.mu.Unlock()
			} else {
				cb.mu.Lock()
				cb.failures = 0
				cb.mu.Unlock()
			}
			return resp, err
		}
	}
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{threshold: threshold, cooldown: cooldown}
}
