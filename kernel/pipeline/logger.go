// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"log"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// Logger 日志中间件。记录每次请求的函数名、地址、耗时和错误信息。
func Logger(logger *log.Logger) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			elapsed := time.Since(start)
			if err != nil {
				logger.Printf("[ERROR] %s %s %v: %v", req.Function, req.Address, elapsed, err)
			} else {
				logger.Printf("[DEBUG] %s %s %v", req.Function, req.Address, elapsed)
			}
			return resp, err
		}
	}
}
