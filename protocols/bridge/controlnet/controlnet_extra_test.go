// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package controlnet

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestParseHexVectors(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{"48656C6C6F", []byte("Hello")},
		{"0x48656C6C6F", []byte("Hello")},
		{"48 65", []byte{0x48, 0x65}},
		{"A", []byte{0x0A}},
	}
	for _, tt := range tests {
		got, err := parseHex(tt.in)
		if err != nil {
			t.Errorf("parseHex(%q) error: %v", tt.in, err)
			continue
		}
		if string(got) != string(tt.want) {
			t.Errorf("parseHex(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseHexEmpty(t *testing.T) {
	if _, err := parseHex(""); err == nil {
		t.Error("expected error for empty input")
	}
	if _, err := parseHex("0x"); err == nil {
		t.Error("expected error for bare prefix")
	}
}

func TestParseHexInvalid(t *testing.T) {
	if _, err := parseHex("GG"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeWriteDefaultAddress(t *testing.T) {
	c, _ := New().NewCodec("coax")
	raw, err := c.Encode(&kernel.Request{Function: "write", Data: []byte{0x01, 0x02}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write 0x00 0102\n" {
		t.Errorf("Encode = %q, want %q", raw, "write 0x00 0102\n")
	}
}

func TestDecodeRawPassthrough(t *testing.T) {
	c, _ := New().NewCodec("coax")
	resp, err := c.Decode([]byte("some freeform text"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["raw"] != "some freeform text" {
		t.Errorf("raw = %v", resp.Metadata["raw"])
	}
}
