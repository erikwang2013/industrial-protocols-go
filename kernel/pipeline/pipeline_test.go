// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func echoHandler(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
	return &kernel.Response{Address: req.Address, Data: req.Data}, nil
}

func TestChainOrder(t *testing.T) {
	order := []string{}

	mk := func(name string) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
				order = append(order, name)
				return next(ctx, req)
			}
		}
	}

	wrapped := Chain(mk("a"), mk("b"), mk("c"))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}

	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Errorf("unexpected order: %v", order)
	}
}

func TestEmptyChain(t *testing.T) {
	wrapped := Chain()(echoHandler)
	resp, err := wrapped(context.Background(), &kernel.Request{Address: "x", Data: []byte("y")})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Address != "x" {
		t.Errorf("expected x, got %s", resp.Address)
	}
}
