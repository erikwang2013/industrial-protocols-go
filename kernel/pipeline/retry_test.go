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
