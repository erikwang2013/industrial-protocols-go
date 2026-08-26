// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package iolink

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncodeWithData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "iowrite",
		Address:  "1/2",
		Data:     []byte{0x0A, 0x0B},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "iowrite 1/2 0a0b\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "iowrite 1/2 0a0b\r\n")
	}
}

func TestExtraEncodeWithoutData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "ioread", Address: "3/0"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "ioread 3/0\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "ioread 3/0\r\n")
	}
}

func TestExtraDecodePassthrough(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	in := []byte("OK")
	resp, err := c.Decode(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != string(in) {
		t.Errorf("data = %q, want %q", resp.Data, in)
	}
}
