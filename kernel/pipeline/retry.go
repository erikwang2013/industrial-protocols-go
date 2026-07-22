// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// BackoffFunc 根据当前尝试次数计算等待时长。
type BackoffFunc func(attempt int) time.Duration

// LinearBackoff 线性退避：每次等待 base * attempt
func LinearBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return base * time.Duration(attempt)
	}
}

// ExponentialBackoff 指数退避：每次等待 base * 2^(attempt-1)
func ExponentialBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := base
		for i := 1; i < attempt; i++ {
			d *= 2
		}
		return d
	}
}

// JitterBackoff 抖动退避：在线性退避基础上增加随机偏移
func JitterBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := base * time.Duration(attempt)
		return d + time.Duration(d/2)
	}
}

// Retry 重试中间件。在 maxAttempts 次内重试失败的请求。
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
