// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package worldfip

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestXorCRC16KnownVector(t *testing.T) {
	// Same CRC-16/ARC family as Modbus; check value for "123456789" is 0x4B37.
	if got := xorCRC16([]byte("123456789")); got != 0x4B37 {
		t.Errorf("xorCRC16(\"123456789\") = 0x%04X, want 0x4B37", got)
	}
}

func TestEncodeProduce(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "produce",
		Data:     []byte{0x01, 0x02},
		Metadata: map[string]any{"tag": float64(0x1234)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if binary.BigEndian.Uint16(raw[0:2]) != worldFIPMagic {
		t.Error("bad magic")
	}
	if raw[2] != worldFIPProduce {
		t.Errorf("op = 0x%02X, want 0x01", raw[2])
	}
	if binary.BigEndian.Uint16(raw[3:5]) != 0x1234 {
		t.Errorf("tag = %04X, want 1234", binary.BigEndian.Uint16(raw[3:5]))
	}
	if binary.BigEndian.Uint16(raw[5:7]) != 2 {
		t.Errorf("length = %d, want 2", binary.BigEndian.Uint16(raw[5:7]))
	}
}

func TestEncodeReadIsConsume(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[2] != worldFIPConsume {
		t.Errorf("op = 0x%02X, want 0x02 (consume)", raw[2])
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{Function: "write", Data: []byte{0xCA, 0xFE}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, []byte{0xCA, 0xFE}) {
		t.Errorf("data = %v, want [CA FE]", resp.Data)
	}
	if resp.Metadata["operation"] != 1 {
		t.Errorf("operation = %v, want 1", resp.Metadata["operation"])
	}
}

func TestDecodeBadMagic(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	frame := []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0xAA, 0x00, 0x00, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for bad magic")
	}
}

func TestDecodeCRCMismatch(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{Function: "read"})
	raw[len(raw)-1] ^= 0xFF
	if _, err := c.Decode(raw); err == nil {
		t.Fatal("expected CRC error")
	}
}

func TestDecodeLengthMismatch(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	// Declared length 100, only 1 payload byte.
	frame := []byte{0x57, 0xF1, 0x01, 0x00, 0x00, 0x00, 0x64, 0xAA, 0x00, 0x00, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected length mismatch error")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	if _, err := c.Decode([]byte{0x57, 0xF1}); err == nil {
		t.Fatal("expected error for short frame")
	}
}
