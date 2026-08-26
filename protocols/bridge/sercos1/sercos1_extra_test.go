// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package sercos1

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestParseHexVectors(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{"1234", []byte{0x12, 0x34}},
		{"0x1234", []byte{0x12, 0x34}},
		{"12 34", []byte{0x12, 0x34}},
		{"1", []byte{0x01}},
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
	if _, err := parseHex("XY"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeReadDefaultsFiber(t *testing.T) {
	c, _ := New().NewCodec("fiber")
	raw, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "read 0x0000 1\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestDecodeDriveMetadata(t *testing.T) {
	c, _ := New().NewCodec("fiber")
	resp, err := c.Decode([]byte("DRIVE: 2"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["state"] != "DRIVE: 2" {
		t.Errorf("state = %v", resp.Metadata["state"])
	}
}
