// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package modbusplus

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestModbusCRC16KnownVector(t *testing.T) {
	// CRC-16/ARC (Modbus) check value for "123456789" is 0x4B37.
	if got := modbusCRC16([]byte("123456789")); got != 0x4B37 {
		t.Errorf("modbusCRC16(\"123456789\") = 0x%04X, want 0x4B37", got)
	}
}

func TestEncodeWriteDestMetadata(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0x11},
		Metadata: map[string]any{"dest": float64(9)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if binary.BigEndian.Uint16(raw[0:2]) != mbpMagic {
		t.Error("bad magic")
	}
	if raw[2] != 9 {
		t.Errorf("dest = %d, want 9", raw[2])
	}
	if raw[3] != mbpCmdWrite {
		t.Errorf("cmd = 0x%02X, want 0x02", raw[3])
	}
	crc := binary.LittleEndian.Uint16(raw[len(raw)-2:])
	if crc != modbusCRC16(raw[:len(raw)-2]) {
		t.Error("bad CRC")
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{Function: "write", Data: []byte{0xDE, 0xAD}, Metadata: map[string]any{"dest": float64(5)}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, []byte{0xDE, 0xAD}) {
		t.Errorf("data = %v, want [DE AD]", resp.Data)
	}
	if resp.Metadata["dest"] != 5 || resp.Metadata["command"] != 2 {
		t.Errorf("metadata = %v", resp.Metadata)
	}
}

func TestDecodeBadMagic(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	frame := []byte{0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0xAA, 0x00, 0x00, 0x00}
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
	frame := []byte{0x4D, 0x42, 0x01, 0x01, 0x00, 0x64, 0xAA, 0x00, 0x00, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected length mismatch error")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	if _, err := c.Decode([]byte{0x4D, 0x42, 0x01}); err == nil {
		t.Fatal("expected error for short frame")
	}
}
