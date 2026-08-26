// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package powerlink

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestParseHexVectors(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{"0102", []byte{0x01, 0x02}},
		{"0x0102", []byte{0x01, 0x02}},
		{"01 02", []byte{0x01, 0x02}},
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
	if _, err := parseHex("QQ"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeWriteDefaultAddress(t *testing.T) {
	c, _ := New().NewCodec("ethernet")
	raw, err := c.Encode(&kernel.Request{Function: "write", Data: []byte{0x0A}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write 0x1000 0A\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestDecodeStateMetadata(t *testing.T) {
	c, _ := New().NewCodec("ethernet")
	resp, err := c.Decode([]byte("NMT_OPERATIONAL"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["state"] != "NMT_OPERATIONAL" {
		t.Errorf("state = %v", resp.Metadata["state"])
	}
}
