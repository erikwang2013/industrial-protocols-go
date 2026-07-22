// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// Handler 是一次协议读写操作的标准函数签名。
// 将 Codec.Encode → Transport.Write → Transport.Read → Codec.Decode 封装为一个函数。
type Handler func(ctx context.Context, req *kernel.Request) (*kernel.Response, error)

// Middleware 包装一个 Handler，返回新的 Handler。
// 中间件链的执行顺序与 Chain 的参数顺序一致（先传入的先执行）。
type Middleware func(next Handler) Handler

// Chain 将多个中间件串联。Chain(a, b, c) 的执行顺序为 a → b → c → handler。
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
