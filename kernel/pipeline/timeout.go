// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// Timeout 超时控制中间件。通过 context.WithTimeout 为每次请求设置截止时间。
func Timeout(d time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(ctx, req)
		}
	}
}
