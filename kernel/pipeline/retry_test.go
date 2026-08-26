// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestRetry_SuccessFirstAttempt(t *testing.T) {
	wrapped := Chain(Retry(3, LinearBackoff(10*time.Millisecond)))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRetry_EventualSuccess(t *testing.T) {
	attempts := 0
	flaky := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		attempts++
		if attempts < 2 {
			return nil, errors.New("transient error")
		}
		return &kernel.Response{}, nil
	}
	wrapped := Chain(Retry(3, LinearBackoff(10*time.Millisecond)))(flaky)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetry_Exhausted(t *testing.T) {
	alwaysFail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("always fail")
	}
	wrapped := Chain(Retry(2, LinearBackoff(10*time.Millisecond)))(alwaysFail)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if !errors.Is(err, kernel.ErrRetryExhausted) {
		t.Errorf("expected ErrRetryExhausted, got %v", err)
	}
}

func TestBackoffFuncs(t *testing.T) {
	lin := LinearBackoff(2 * time.Millisecond)
	if got := lin(3); got != 6*time.Millisecond {
		t.Errorf("LinearBackoff(2ms)(3) = %v, want 6ms", got)
	}

	exp := ExponentialBackoff(2 * time.Millisecond)
	if got := exp(1); got != 2*time.Millisecond {
		t.Errorf("ExponentialBackoff(2ms)(1) = %v, want 2ms", got)
	}
	if got := exp(3); got != 8*time.Millisecond {
		t.Errorf("ExponentialBackoff(2ms)(3) = %v, want 8ms", got)
	}
	if got := exp(4); got != 16*time.Millisecond {
		t.Errorf("ExponentialBackoff(2ms)(4) = %v, want 16ms", got)
	}

	jit := JitterBackoff(2 * time.Millisecond)
	if got := jit(2); got != 6*time.Millisecond {
		t.Errorf("JitterBackoff(2ms)(2) = %v, want 6ms", got)
	}
}

func TestRetry_ZeroAttempts(t *testing.T) {
	calls := 0
	h := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		calls++
		return nil, errors.New("x")
	}
	wrapped := Chain(Retry(0, nil))(h)
	if _, err := wrapped(context.Background(), &kernel.Request{}); !errors.Is(err, kernel.ErrRetryExhausted) {
		t.Errorf("expected ErrRetryExhausted, got %v", err)
	}
	if calls != 0 {
		t.Errorf("handler should not be called with maxAttempts=0, got %d calls", calls)
	}
}

func TestRetry_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	alwaysFail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("always fail")
	}
	wrapped := Chain(Retry(3, LinearBackoff(time.Hour)))(alwaysFail)
	_, err := wrapped(ctx, &kernel.Request{})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
