// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package controlnet

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "controlnet" {
		t.Errorf("expected controlnet, got %s", p.Name())
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
	for _, v := range []string{"coax", "cmd"} {
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
		Address:  "0x10",
		Count:    8,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "read 0x10 8\n"
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
	expected := "read 0x00 1\n"
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
		Address:  "0x20",
		Data:     []byte{0x55, 0xAA},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "write 0x20 55AA\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeStatus(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "status"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "status\n"
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

func TestDecodeStatus(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("STATUS: Online, Node 3"))
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
	resp, err := c.Decode([]byte("55 AA BB CC"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x55, 0xAA, 0xBB, 0xCC}
	if string(resp.Data) != string(expected) {
		t.Errorf("expected %X, got %X", expected, resp.Data)
	}
}

func TestDriverNewCmdDriver(t *testing.T) {
	b, c, err := NewCmdDriver("/usr/bin/1784-pcic-cli")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}
