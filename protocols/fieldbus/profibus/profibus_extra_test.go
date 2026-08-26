// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package profibus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncode(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "read", Address: "5"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "read 5\r\n" {
		t.Errorf("Encode = %q, want %q", raw, "read 5\r\n")
	}
}

func TestExtraEncodeIgnoresData(t *testing.T) {
	// Current behavior: the gateway codec drops req.Data entirely.
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "write", Address: "3", Data: []byte{0xAB}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write 3\r\n" {
		t.Errorf("Encode = %q, want %q (data dropped)", raw, "write 3\r\n")
	}
}

func TestExtraDecodePassthrough(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	in := []byte("OK: 00 01")
	resp, err := c.Decode(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != string(in) {
		t.Errorf("data = %q, want %q", resp.Data, in)
	}
}
