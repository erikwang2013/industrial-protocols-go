// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package safej1850

import (
	"errors"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestNewCodecRejected(t *testing.T) {
	// The safej1850 module declares the serial variant but does not implement
	// a codec; NewCodec must fail with the invalid-address sentinel.
	_, err := New().NewCodec("serial")
	if err == nil {
		t.Fatal("expected error from NewCodec")
	}
	if !errors.Is(err, kernel.ErrInvalidAddress) {
		t.Errorf("expected kernel.ErrInvalidAddress, got %v", err)
	}
}

func TestDefaultPort(t *testing.T) {
	if New().DefaultPort() != 0 {
		t.Errorf("DefaultPort = %d, want 0", New().DefaultPort())
	}
}

func TestImplementsKernelInterface(t *testing.T) {
	var _ kernel.Protocol = New()
}
