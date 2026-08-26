// Copyright (c) 2026 erik <erik@erik.xyz] — https://erik.xyz

package lonworks

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestEncodeWithData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "nv_write",
		Address:  "nvoTemp",
		Data:     []byte{0xAB, 0xCD},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "nv_write nvoTemp abcd\r\n"
	if string(raw) != want {
		t.Errorf("Encode = %q, want %q", raw, want)
	}
}

func TestEncodeWithoutData(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "nv_read", Address: "nviSetpoint"})
	if err != nil {
		t.Fatal(err)
	}
	want := "nv_read nviSetpoint\r\n"
	if string(raw) != want {
		t.Errorf("Encode = %q, want %q", raw, want)
	}
}

func TestDecodePassthrough(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	in := []byte("OK: 12 34")
	resp, err := c.Decode(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != string(in) {
		t.Errorf("data = %q, want %q", resp.Data, in)
	}
}
