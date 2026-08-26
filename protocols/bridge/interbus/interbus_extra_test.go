// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package interbus

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestEncodeWrite(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0xAB, 0xCD},
		Metadata: map[string]any{"slave": float64(3)},
	})
	if err != nil {
		t.Fatal(err)
	}
	// magic(2) + cmd(1) + slave(1) + length(2) + payload(2) + crc(1)
	if len(raw) != 9 {
		t.Fatalf("len = %d, want 9", len(raw))
	}
	if binary.BigEndian.Uint16(raw[0:2]) != interbusMagic {
		t.Error("bad magic")
	}
	if raw[2] != interbusCmdWrite {
		t.Errorf("cmd = 0x%02X, want 0x02", raw[2])
	}
	if raw[3] != 3 {
		t.Errorf("slave = %d, want 3", raw[3])
	}
	if binary.BigEndian.Uint16(raw[4:6]) != 2 {
		t.Errorf("length = %d, want 2", binary.BigEndian.Uint16(raw[4:6]))
	}
	if !bytes.Equal(raw[6:8], []byte{0xAB, 0xCD}) {
		t.Errorf("payload = %v, want [AB CD]", raw[6:8])
	}
	if raw[8] != checksum8(raw[:8]) {
		t.Error("bad checksum")
	}
}

func TestEncodeReadDefaultSlave(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[2] != interbusCmdRead || raw[3] != 1 {
		t.Errorf("cmd/slave = %02X/%02X, want 01/01", raw[2], raw[3])
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{Function: "write", Data: []byte{1, 2, 3}, Metadata: map[string]any{"slave": float64(7)}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, []byte{1, 2, 3}) {
		t.Errorf("data = %v, want [1 2 3]", resp.Data)
	}
	if resp.Metadata["slave"] != 7 || resp.Metadata["command"] != 2 {
		t.Errorf("metadata = %v", resp.Metadata)
	}
}

func TestDecodeBadMagic(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	frame := []byte{0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x01}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for bad magic")
	}
}

func TestDecodeChecksumMismatch(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{Function: "read"})
	raw[len(raw)-1] ^= 0xFF
	if _, err := c.Decode(raw); err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestDecodeLengthMismatch(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	// Declared length 100 but only 1 payload byte present.
	frame := []byte{0x1B, 0x55, 0x01, 0x01, 0x00, 0x64, 0xAA, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected length mismatch error")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	if _, err := c.Decode([]byte{0x1B, 0x55}); err == nil {
		t.Fatal("expected error for short frame")
	}
}

func TestChecksum8(t *testing.T) {
	if got := checksum8([]byte{0x1B, 0x55, 0x01}); got != 0x1B^0x55^0x01 {
		t.Errorf("checksum8 = 0x%02X, want 0x%02X", got, 0x1B^0x55^0x01)
	}
	// XOR checksum of identical bytes cancels out.
	if checksum8([]byte{0xAA, 0xAA}) != 0 {
		t.Error("XOR of duplicate byte should be 0")
	}
}
