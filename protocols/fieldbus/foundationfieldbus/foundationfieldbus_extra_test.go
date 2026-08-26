// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package foundationfieldbus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncodeWithData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "write",
		Address:  "FIC101",
		Data:     []byte{0x0A, 0x0B},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write FIC101 0a0b\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "write FIC101 0a0b\r\n")
	}
}

func TestExtraEncodeWithoutData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "read", Address: "TT001"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "read TT001\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "read TT001\r\n")
	}
}

func TestExtraDecodePassthrough(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	in := []byte("ACK: 0001")
	resp, err := c.Decode(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != string(in) {
		t.Errorf("data = %q, want %q", resp.Data, in)
	}
}
