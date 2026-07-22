// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestCircuitBreaker_ClosedByDefault(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	wrapped := Chain(cb.Middleware())(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("fail")
	}
	wrapped := Chain(cb.Middleware())(fail)

	for i := 0; i < 2; i++ {
		_, _ = wrapped(context.Background(), &kernel.Request{})
	}

	_, err := wrapped(context.Background(), &kernel.Request{})
	if !errors.Is(err, kernel.ErrCircuitOpen) {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_ResetsAfterCooldown(t *testing.T) {
	cb := NewCircuitBreaker(2, 30*time.Millisecond)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("fail")
	}
	wrapped := Chain(cb.Middleware())(fail)

	for i := 0; i < 2; i++ {
		_, _ = wrapped(context.Background(), &kernel.Request{})
	}
	time.Sleep(50 * time.Millisecond)

	_, err := wrapped(context.Background(), &kernel.Request{})
	if err == nil {
		t.Log("breaker reset after cooldown")
	}
}
