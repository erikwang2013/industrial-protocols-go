// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package isa100

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestParseHexVectors(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{"DEAD", []byte{0xDE, 0xAD}},
		{"0xDEAD", []byte{0xDE, 0xAD}},
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
}

func TestParseHexInvalid(t *testing.T) {
	if _, err := parseHex("xyz"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeWriteFormat(t *testing.T) {
	c, _ := New().NewCodec("wireless")
	raw, err := c.Encode(&kernel.Request{Function: "write", Address: "TAG01", Data: []byte{0x01}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write TAG01 01\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestDecodeDeviceMetadata(t *testing.T) {
	c, _ := New().NewCodec("wireless")
	resp, err := c.Decode([]byte("DEVICE: 1"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["device"] != "DEVICE: 1" {
		t.Errorf("device = %v", resp.Metadata["device"])
	}
}
