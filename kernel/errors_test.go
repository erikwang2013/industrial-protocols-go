// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kernel

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	errs := []error{
		ErrTimeout,
		ErrCircuitOpen,
		ErrRetryExhausted,
		ErrTransportClosed,
		ErrInvalidAddress,
	}
	for _, e := range errs {
		if e.Error() == "" {
			t.Errorf("sentinel error %v has empty message", e)
		}
	}
}

func TestProtocolError(t *testing.T) {
	raw := []byte{0x01, 0x83, 0x02}
	pe := &ProtocolError{Code: "02", Message: "illegal data address", Raw: raw}

	if pe.Error() != "protocol error [02]: illegal data address" {
		t.Errorf("unexpected error string: %s", pe.Error())
	}
	if !errors.Is(pe, pe) {
		t.Error("ProtocolError should be errors.Is to itself")
	}
}
