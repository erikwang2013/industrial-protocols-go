// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package isa100

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "isa100" {
		t.Errorf("expected isa100, got %s", p.Name())
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
	for _, v := range []string{"wireless", "cmd"} {
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
		Address:  "DEV001",
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "read DEV001\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeReadDefault(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "read DEV001\n"
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
		Address:  "DEV001",
		Data:     []byte{0xA5, 0x5A},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "write DEV001 A55A\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeList(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "list"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "list\n"
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

func TestDecodeDevice(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("DEVICE: DEV001, RSSI: -45 dBm"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["device"] == nil {
		t.Error("expected device metadata")
	}
}

func TestDecodeHex(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("A5 5A"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0xA5, 0x5A}
	if string(resp.Data) != string(expected) {
		t.Errorf("expected %X, got %X", expected, resp.Data)
	}
}

func TestDriverNewCmdDriver(t *testing.T) {
	b, c, err := NewCmdDriver("/usr/bin/yfgw410_cli")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}
