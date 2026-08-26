// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package iec61850

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncodeValue(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want []byte
	}{
		{"empty defaults to bool false", nil, []byte{TagBool, 1, 0}},
		{"single byte is integer", []byte{0x05}, []byte{TagInteger, 1, 0x05}},
		{"two bytes are unsigned", []byte{0x01, 0x02}, []byte{TagUnsigned, 2, 0x01, 0x02}},
		{"four bytes are unsigned", []byte{1, 2, 3, 4}, []byte{TagUnsigned, 4, 1, 2, 3, 4}},
		{"five bytes are octet string", []byte{1, 2, 3, 4, 5}, []byte{TagOctetString, 5, 1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		if got := encodeValue(tt.in); !bytes.Equal(got, tt.want) {
			t.Errorf("encodeValue(%s) = % X, want % X", tt.name, got, tt.want)
		}
	}
}

func TestExtraDecodeValue(t *testing.T) {
	if got, err := decodeValue([]byte{TagBool, 1, 0xFF}); err != nil || !bytes.Equal(got, []byte{0xFF}) {
		t.Errorf("bool = % X, err %v", got, err)
	}
	// Float: 4-byte value preserved.
	if got, err := decodeValue([]byte{TagFloat, 4, 0x3F, 0x80, 0x00, 0x00}); err != nil || !bytes.Equal(got, []byte{0x3F, 0x80, 0x00, 0x00}) {
		t.Errorf("float = % X, err %v", got, err)
	}
	if _, err := decodeValue([]byte{TagInteger, 1}); err == nil {
		t.Error("expected error for truncated value")
	}
	if _, err := decodeValue([]byte{TagInteger}); err == nil {
		t.Error("expected error for value too short")
	}
}

func TestExtraEncodeVisibleString(t *testing.T) {
	if got := encodeVisibleString("AB"); !bytes.Equal(got, []byte{TagVisibleString, 2, 'A', 'B'}) {
		t.Errorf("encodeVisibleString = % X", got)
	}
	if got := encodeVisibleString(""); !bytes.Equal(got, []byte{TagVisibleString, 0}) {
		t.Errorf("encodeVisibleString(\"\") = % X", got)
	}
}

func TestExtraEncodeReadLayout(t *testing.T) {
	c, _ := New().NewCodec("mms")
	raw, err := c.Encode(&kernel.Request{Function: "read", Address: "VAR1"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[0] != 0xA4 || raw[1] != 0x0A || raw[2] != 1 {
		t.Errorf("header = % X, want A4 0A 01", raw[:3])
	}
	// VariableSpec: [0C 08 02 06 8A 04 V A R 1]; the inner length field is
	// len(encodeVisibleString(name)) = 6 (tag + len + 4 chars).
	wantSpec := []byte{0x0C, 0x08, 0x02, 0x06, 0x8A, 0x04, 'V', 'A', 'R', '1'}
	if !bytes.Equal(raw[3:], wantSpec) {
		t.Errorf("var spec = % X, want % X", raw[3:], wantSpec)
	}
}

func TestExtraEncodeWriteLayout(t *testing.T) {
	c, _ := New().NewCodec("mms")
	raw, err := c.Encode(&kernel.Request{Function: "write", Address: "X", Data: []byte{0x01}})
	if err != nil {
		t.Fatal(err)
	}
	// [A5 08 invokeID 02 03 8A 01 'X' 85 01 01]
	want := []byte{0xA5, 0x08, 1, 0x02, 0x03, 0x8A, 0x01, 'X', 0x85, 0x01, 0x01}
	if !bytes.Equal(raw, want) {
		t.Errorf("write = % X, want % X", raw, want)
	}
}

func TestExtraInvokeIDIncrements(t *testing.T) {
	c, _ := New().NewCodec("mms")
	raw1, _ := c.Encode(&kernel.Request{Function: "read", Address: "A"})
	raw2, _ := c.Encode(&kernel.Request{Function: "read", Address: "B"})
	if raw1[2] != 1 || raw2[2] != 2 {
		t.Errorf("invoke id went %d -> %d, want 1 -> 2", raw1[2], raw2[2])
	}
}

func TestExtraDecodeReadResponseValue(t *testing.T) {
	c, _ := New().NewCodec("mms")
	// A4 + length 3 + TagInteger(1) = 42.
	resp, err := c.Decode([]byte{0xA4, 0x03, TagInteger, 1, 42})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "read_response" {
		t.Errorf("function = %q, want read_response", resp.Function)
	}
	if len(resp.Data) != 1 || resp.Data[0] != 42 {
		t.Errorf("data = %v, want [42]", resp.Data)
	}
}

func TestExtraDecodeUnknownTag(t *testing.T) {
	c, _ := New().NewCodec("mms")
	resp, err := c.Decode([]byte{0x88, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["tag"] != int(0x88) {
		t.Errorf("tag = %v, want 136", resp.Metadata["tag"])
	}
}

func TestExtraDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("mms")
	if _, err := c.Decode([]byte{0xA4}); err == nil {
		t.Fatal("expected error for 1-byte frame")
	}
}
