// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package asinterface

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncodeWithData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "aout",
		Address:  "1",
		Data:     []byte{0x0A, 0x0B},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "aout 1 0a0b\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "aout 1 0a0b\r\n")
	}
}

func TestExtraEncodeWithoutData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "status", Address: "3"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "status 3\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "status 3\r\n")
	}
}

func TestExtraDecodePassthrough(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	in := []byte("OK: 01")
	resp, err := c.Decode(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != string(in) {
		t.Errorf("data = %q, want %q", resp.Data, in)
	}
}
