package pipeline

import (
	"context"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type Handler func(ctx context.Context, req *kernel.Request) (*kernel.Response, error)

type Middleware func(next Handler) Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
