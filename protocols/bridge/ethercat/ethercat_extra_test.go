// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package ethercat

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
		{"01 02 03", []byte{0x01, 0x02, 0x03}},
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
	if _, err := parseHex("0xZZ"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeDownloadFormat(t *testing.T) {
	c, _ := New().NewCodec("cmd")
	raw, err := c.Encode(&kernel.Request{Function: "download", Data: []byte{0xDE, 0xAD}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "download 0x0000 0xDEAD\n" {
		t.Errorf("Encode = %q", raw)
	}
}

func TestDecodeRawPassthrough(t *testing.T) {
	c, _ := New().NewCodec("cmd")
	resp, err := c.Decode([]byte("slave info text"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["raw"] != "slave info text" {
		t.Errorf("raw = %v", resp.Metadata["raw"])
	}
}
