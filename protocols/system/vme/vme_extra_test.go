// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package vme

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestNewCodecAcceptsAnyVariant(t *testing.T) {
	for _, v := range []string{"vme", "", "bogus"} {
		if _, err := New().NewCodec(v); err != nil {
			t.Errorf("NewCodec(%q) unexpected error: %v", v, err)
		}
	}
}

func TestCodecPassthroughBytes(t *testing.T) {
	c, _ := New().NewCodec("vme")
	in := []byte{0x01, 0x02, 0x03}
	out, _ := c.Encode(&kernel.Request{Data: in})
	if !bytes.Equal(out, in) {
		t.Errorf("Encode = %v, want %v", out, in)
	}
}

func TestCodecNilData(t *testing.T) {
	c, _ := New().NewCodec("vme")
	if out, _ := c.Encode(&kernel.Request{}); out != nil {
		t.Errorf("Encode(nil) = %v, want nil", out)
	}
}
