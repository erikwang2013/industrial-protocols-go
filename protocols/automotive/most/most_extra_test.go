// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package most

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
		{"48 65 6C 6C 6F", []byte("Hello")},
		{"0x48 0x65", []byte{0x48, 0x65}},
		{"0X48", []byte{0x48}},
		{"A", []byte{0x0A}}, // odd length left-pads with 0
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

func TestParseHexInvalid(t *testing.T) {
	if _, err := parseHex("ZZ"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestEncodeWriteDefaultAddress(t *testing.T) {
	c, _ := New().NewCodec("optical")
	raw, err := c.Encode(&kernel.Request{Function: "write", Data: []byte{0x01, 0x02}})
	if err != nil {
		t.Fatal(err)
	}
	want := "AT+WRITE=0x0000,0102\r\n"
	if string(raw) != want {
		t.Errorf("Encode = %q, want %q", raw, want)
	}
}

func TestEncodeReadCount(t *testing.T) {
	c, _ := New().NewCodec("optical")
	raw, err := c.Encode(&kernel.Request{Function: "read", Address: "0x1000", Count: 4})
	if err != nil {
		t.Fatal(err)
	}
	want := "AT+READ=0x1000,4\r\n"
	if string(raw) != want {
		t.Errorf("Encode = %q, want %q", raw, want)
	}
}

func TestDecodeOKInvalidHexKeepsRaw(t *testing.T) {
	c, _ := New().NewCodec("optical")
	resp, err := c.Decode([]byte("+OK:NOTHEX"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["raw"] != "NOTHEX" {
		t.Errorf("raw = %v, want NOTHEX", resp.Metadata["raw"])
	}
}

func TestDecodeEmptyLine(t *testing.T) {
	c, _ := New().NewCodec("optical")
	resp, err := c.Decode([]byte(""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["raw"] != "" {
		t.Errorf("raw = %v, want empty", resp.Metadata["raw"])
	}
}
