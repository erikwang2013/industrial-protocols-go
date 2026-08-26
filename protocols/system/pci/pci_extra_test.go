// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pci

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestCodecNilData(t *testing.T) {
	c, _ := New().NewCodec("pci")
	out, err := c.Encode(&kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Errorf("Encode(nil data) = %v, want nil", out)
	}
	resp, err := c.Decode(nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data != nil {
		t.Errorf("Decode(nil) data = %v, want nil", resp.Data)
	}
}

func TestCodecPassthroughBytes(t *testing.T) {
	c, _ := New().NewCodec("pci")
	in := []byte{0xAA, 0x55}
	out, _ := c.Encode(&kernel.Request{Data: in})
	if !bytes.Equal(out, in) {
		t.Errorf("Encode = %v, want %v", out, in)
	}
}
