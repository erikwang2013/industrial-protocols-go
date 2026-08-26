// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cpci

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestNewCodecAcceptsAnyVariant(t *testing.T) {
	for _, v := range []string{"cpci", "", "bogus"} {
		if _, err := New().NewCodec(v); err != nil {
			t.Errorf("NewCodec(%q) unexpected error: %v", v, err)
		}
	}
}

func TestCodecPassthroughBytes(t *testing.T) {
	c, _ := New().NewCodec("cpci")
	in := []byte{0x00, 0x01, 0x10, 0xFF}
	out, err := c.Encode(&kernel.Request{Data: in})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, in) {
		t.Errorf("Encode = %v, want %v", out, in)
	}
	resp, err := c.Decode(in)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, in) {
		t.Errorf("Decode data = %v, want %v", resp.Data, in)
	}
}

func TestCodecNilData(t *testing.T) {
	c, _ := New().NewCodec("cpci")
	out, err := c.Encode(&kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Errorf("Encode(nil data) = %v, want nil", out)
	}
}
