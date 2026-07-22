// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestTimeout_NoTimeout(t *testing.T) {
	wrapped := Chain(Timeout(100 * time.Millisecond))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTimeout_Exceeded(t *testing.T) {
	slow := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return &kernel.Response{}, nil
		}
	}
	wrapped := Chain(Timeout(10 * time.Millisecond))(slow)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
