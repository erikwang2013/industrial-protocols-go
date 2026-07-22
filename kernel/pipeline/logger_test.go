package pipeline

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	wrapped := Chain(Logger(logger))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{Function: "read", Address: "40001"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "read") {
		t.Errorf("log should contain function name, got: %s", buf.String())
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("boom")
	}
	wrapped := Chain(Logger(logger))(fail)
	_, _ = wrapped(context.Background(), &kernel.Request{Function: "write"})
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("log should contain ERROR, got: %s", buf.String())
	}
}
