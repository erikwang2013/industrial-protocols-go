// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package sercos

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestParseHexVectors(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{"AABB", []byte{0xAA, 0xBB}},
		{"0xAABB", []byte{0xAA, 0xBB}},
		{"AA BB", []byte{0xAA, 0xBB}},
		{"B", []byte{0x0B}},
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
}

func TestParseHexInvalid(t *testing.T) {
	if _, err := parseHex("zz"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodePhaseAddress(t *testing.T) {
	c, _ := New().NewCodec("fiber")
	raw, err := c.Encode(&kernel.Request{Function: "phase", Address: "3"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "phase 3\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestEncodeWriteDefaultAddress(t *testing.T) {
	c, _ := New().NewCodec("fiber")
	raw, err := c.Encode(&kernel.Request{Function: "write", Data: []byte{0x01}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write S-0-0 01\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestDecodePhaseMetadata(t *testing.T) {
	c, _ := New().NewCodec("fiber")
	resp, err := c.Decode([]byte("PHASE: 4"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["state"] != "PHASE: 4" {
		t.Errorf("state = %v", resp.Metadata["state"])
	}
}
