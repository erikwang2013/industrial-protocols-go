// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package asinterface

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "asinterface" {
		t.Errorf("expected asinterface, got %s", p.Name())
	}
	if p.DefaultPort() != 2003 {
		t.Errorf("expected 2003, got %d", p.DefaultPort())
	}
}

func TestNewCodecInvalid(t *testing.T) {
	_, err := New().NewCodec("serial")
	if err == nil {
		t.Error("expected error for unsupported variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "1A"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty output")
	}
}

func TestEncodeWrite(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Address:  "5B",
		Data:     []byte{0x0F},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty output")
	}
}

func TestDecode(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("OK\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}
