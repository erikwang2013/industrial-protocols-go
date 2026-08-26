// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package wirelesshart

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
		{"F", []byte{0x0F}},
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
	if _, err := parseHex("NOPE"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeWriteFormat(t *testing.T) {
	c, _ := New().NewCodec("wireless")
	raw, err := c.Encode(&kernel.Request{Function: "write", Address: "TAG07", Data: []byte{0xAB}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "write TAG07 AB\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestDecodeTagMetadata(t *testing.T) {
	c, _ := New().NewCodec("wireless")
	resp, err := c.Decode([]byte("TAG: DEV42"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["device"] != "TAG: DEV42" {
		t.Errorf("device = %v", resp.Metadata["device"])
	}
}
