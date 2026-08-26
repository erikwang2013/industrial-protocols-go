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

func TestCircuitBreaker_SuccessResetsFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Hour)
	shouldFail := true
	h := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		if shouldFail {
			return nil, errors.New("fail")
		}
		return &kernel.Response{}, nil
	}
	wrapped := Chain(cb.Middleware())(h)

	_, _ = wrapped(context.Background(), &kernel.Request{}) // fail 1
	shouldFail = false
	_, _ = wrapped(context.Background(), &kernel.Request{}) // success: must reset the counter
	shouldFail = true
	_, _ = wrapped(context.Background(), &kernel.Request{}) // fail 1
	_, _ = wrapped(context.Background(), &kernel.Request{}) // fail 2

	if cb.open {
		t.Error("breaker should still be closed after 2 non-consecutive failures")
	}
	if cb.failures != 2 {
		t.Errorf("expected 2 failures recorded, got %d", cb.failures)
	}

	if _, err := wrapped(context.Background(), &kernel.Request{}); err == nil { // fail 3 -> open
		t.Error("expected failure from handler")
	}
	_, err := wrapped(context.Background(), &kernel.Request{}) // now open
	if !errors.Is(err, kernel.ErrCircuitOpen) {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenAllowsRequest(t *testing.T) {
	cb := NewCircuitBreaker(1, 20*time.Millisecond)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("fail")
	}
	wrapped := Chain(cb.Middleware())(fail)

	_, _ = wrapped(context.Background(), &kernel.Request{}) // open
	if _, err := wrapped(context.Background(), &kernel.Request{}); !errors.Is(err, kernel.ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen while open, got %v", err)
	}

	time.Sleep(40 * time.Millisecond) // past cooldown
	_, err := wrapped(context.Background(), &kernel.Request{})
	if errors.Is(err, kernel.ErrCircuitOpen) {
		t.Error("half-open request should reach the handler, not be rejected")
	}
	if err == nil {
		t.Error("handler should have failed again")
	}
}
