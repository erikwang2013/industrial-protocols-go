// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package lightbus

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestEncodeWriteWithMetadata(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0x01},
		Metadata: map[string]any{"module": float64(4), "channel": float64(0x0203)},
	})
	if err != nil {
		t.Fatal(err)
	}
	// magic(2) + cmd(1) + module(1) + channel(2) + length(2) + payload(1) + crc(1)
	if len(raw) != 10 {
		t.Fatalf("len = %d, want 10", len(raw))
	}
	if binary.BigEndian.Uint16(raw[0:2]) != lightbusMagic {
		t.Error("bad magic")
	}
	if raw[2] != lightbusCmdWrite {
		t.Errorf("cmd = 0x%02X, want 0x02", raw[2])
	}
	if raw[3] != 4 {
		t.Errorf("module = %d, want 4", raw[3])
	}
	if binary.BigEndian.Uint16(raw[4:6]) != 0x0203 {
		t.Errorf("channel = %04X, want 0203", binary.BigEndian.Uint16(raw[4:6]))
	}
	if binary.BigEndian.Uint16(raw[6:8]) != 1 {
		t.Errorf("length = %d, want 1", binary.BigEndian.Uint16(raw[6:8]))
	}
	if raw[9] != checksum8XOR(raw[:9]) {
		t.Error("bad checksum")
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{
		Function: "read",
		Data:     []byte{0x10, 0x20},
		Metadata: map[string]any{"module": float64(2)},
	})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, []byte{0x10, 0x20}) {
		t.Errorf("data = %v, want [10 20]", resp.Data)
	}
	if resp.Metadata["module"] != 2 || resp.Metadata["channel"] != 0 {
		t.Errorf("metadata = %v", resp.Metadata)
	}
}

func TestDecodeBadMagic(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	if _, err := c.Decode([]byte{0x00, 0x00, 0x01, 0x01, 0, 0, 0, 1, 0x00, 0x00}); err == nil {
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
	// Declared length 100, only 1 payload byte present.
	frame := []byte{0x4C, 0x42, 0x01, 0x01, 0, 0, 0x00, 0x64, 0xAA, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected length mismatch error")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	if _, err := c.Decode([]byte{0x4C, 0x42}); err == nil {
		t.Fatal("expected error for short frame")
	}
}

func TestChecksum8XOR(t *testing.T) {
	if got := checksum8XOR([]byte{0x4C, 0x42, 0x01}); got != 0x4C^0x42^0x01 {
		t.Errorf("checksum8XOR = 0x%02X, want 0x%02X", got, 0x4C^0x42^0x01)
	}
	if checksum8XOR([]byte{0x01, 0x01, 0x01}) != 0x01 {
		t.Error("XOR of odd duplicate bytes should equal the byte")
	}
}
