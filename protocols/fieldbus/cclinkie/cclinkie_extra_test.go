// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cclinkie

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraCRC16CCITTKnownVector(t *testing.T) {
	// Published CRC-16/CCITT-FALSE check value.
	if got := crc16CCITT([]byte("123456789")); got != 0x29B1 {
		t.Errorf("crc16CCITT(123456789) = 0x%04X, want 0x29B1", got)
	}
}

func TestExtraEncodeWriteLayout(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0xAB},
		Metadata: map[string]any{"station": float64(2)},
	})
	if err != nil {
		t.Fatal(err)
	}
	// magic(2) | cmd(1) | station(2) | length(2) | payload(1) | crc(2, BE)
	want := []byte{0xCC, 0x1E, 0x02, 0x00, 0x02, 0x00, 0x01, 0xAB, 0x7F, 0xB3}
	if !bytes.Equal(raw, want) {
		t.Errorf("Encode = % X, want % X", raw, want)
	}
}

func TestExtraEncodeReadDefault(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[2] != cclinkIECmdRead {
		t.Errorf("cmd = 0x%02X, want 0x01", raw[2])
	}
	// Default station 1.
	if raw[3] != 0x00 || raw[4] != 0x01 {
		t.Errorf("station = %02X %02X, want 00 01", raw[3], raw[4])
	}
}

func TestExtraDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, _ := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0x01, 0x02},
		Metadata: map[string]any{"station": uint16(7)},
	})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, []byte{0x01, 0x02}) {
		t.Errorf("data = % X, want 01 02", resp.Data)
	}
	if resp.Metadata["station"] != 7 || resp.Metadata["command"] != int(cclinkIECmdWrite) {
		t.Errorf("metadata = %v", resp.Metadata)
	}
}

func TestExtraDecodeErrors(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	if _, err := c.Decode(make([]byte, 8)); err == nil {
		t.Error("expected error for 8-byte frame")
	}
	bad := make([]byte, 9)
	bad[0] = 0xDE
	if _, err := c.Decode(bad); err == nil {
		t.Error("expected error for invalid magic")
	}
	// length = 5 but no payload follows.
	raw, _ := c.Encode(&kernel.Request{Function: "read"})
	raw[5] = 0x00
	raw[6] = 0x05
	if _, err := c.Decode(raw); err == nil {
		t.Error("expected error for data length mismatch")
	}
	// Corrupt CRC.
	raw, _ = c.Encode(&kernel.Request{Function: "read"})
	raw[len(raw)-1] ^= 0xFF
	if _, err := c.Decode(raw); err == nil {
		t.Error("expected error for CRC mismatch")
	}
}
