// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package foundationfieldbus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "foundationfieldbus" {
		t.Errorf("expected foundationfieldbus, got %s", p.Name())
	}
	if p.DefaultPort() != 2002 {
		t.Errorf("expected 2002, got %d", p.DefaultPort())
	}
}

func TestNewCodecInvalid(t *testing.T) {
	_, err := New().NewCodec("h1")
	if err == nil {
		t.Error("expected error for unsupported variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "TT101.PV"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty output")
	}
}

func TestEncodeWriteWithData(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Address:  "FV101.SP",
		Data:     []byte{0x41, 0xA0, 0x00, 0x00},
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
	resp, err := c.Decode([]byte("OK 41a00000\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}
