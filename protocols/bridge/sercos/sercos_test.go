// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package sercos

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "sercos" {
		t.Errorf("expected sercos, got %s", p.Name())
	}
	if len(p.Variants()) == 0 {
		t.Error("should have at least one variant")
	}
	if p.DefaultPort() < 0 {
		t.Error("default port must not be negative")
	}
}

func TestProtocol_NewCodec(t *testing.T) {
	p := New()
	for _, v := range []string{"fiber", "cmd"} {
		c, err := p.NewCodec(v)
		if err != nil {
			t.Errorf("NewCodec(%q) returned error: %v", v, err)
		}
		if c == nil {
			t.Errorf("NewCodec(%q) returned nil codec", v)
		}
	}
}

func TestProtocol_RejectsUnknownVariant(t *testing.T) {
	p := New()
	_, err := p.NewCodec("invalid")
	if err == nil {
		t.Error("expected error for unknown variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "read",
		Address:  "S-0-51",
		Count:    2,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "read S-0-51 2\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeReadDefaults(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "read S-0-0 1\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeWrite(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Address:  "S-0-51",
		Data:     []byte{0xAA, 0xBB},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "write S-0-51 AABB\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodePhase(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "phase",
		Address:  "4",
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "phase 4\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodePhaseDefault(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "phase"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "phase 0\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Encode(&kernel.Request{Function: "bogus"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestDecodeState(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("PHASE: 4, STATE: OPERATIONAL"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["state"] == nil {
		t.Error("expected state metadata")
	}
}

func TestDecodeHex(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("AABB CCDD"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	if string(resp.Data) != string(expected) {
		t.Errorf("expected %X, got %X", expected, resp.Data)
	}
}

func TestDriverNewCmdDriver(t *testing.T) {
	b, c, err := NewCmdDriver("/usr/bin/netx_cli")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}
