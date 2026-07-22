package pipeline

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type BackoffFunc func(attempt int) time.Duration

func LinearBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return base * time.Duration(attempt)
	}
}

func ExponentialBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := base
		for i := 1; i < attempt; i++ {
			d *= 2
		}
		return d
	}
}

func JitterBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := base * time.Duration(attempt)
		return d + time.Duration(d/2)
	}
}

func Retry(maxAttempts int, backoff BackoffFunc) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			for i := 0; i < maxAttempts; i++ {
				resp, err := next(ctx, req)
				if err == nil {
					return resp, nil
				}
				if i < maxAttempts-1 {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(backoff(i + 1)):
					}
				}
			}
			return nil, kernel.ErrRetryExhausted
		}
	}
}
