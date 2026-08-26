// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package modbus

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraParseAddr(t *testing.T) {
	tests := []struct {
		in   string
		want uint16
	}{
		{"40001", 40001},
		{"12", 12},
		{"0", 0},
		{"0x10", 0}, // %d stops at 'x'
		{"abc", 0},
		{"", 0},
	}
	for _, tt := range tests {
		if got := parseAddr(tt.in); got != tt.want {
			t.Errorf("parseAddr(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestExtraDecodePDUByteCountPanic(t *testing.T) {
	// Regression: decodePDU used to panic on pdu[2:2+pdu[1]] when the byte
	// count exceeded the PDU; it must return an error instead.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decodePDU panicked on byte-count overflow: %v", r)
		}
	}()
	c, _ := NewCodec("tcp")
	// FC=3, byte count=0x0A but only 2 data bytes follow.
	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x01, 0x03, 0x0A, 0x01, 0x02}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for byte count exceeding PDU length")
	}
}

func TestExtraDecodePDUExceptionIndexPanic(t *testing.T) {
	// Regression: a 1-byte exception PDU (0x81) used to panic at pdu[1];
	// it must return an error instead.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decodePDU panicked on 1-byte exception PDU: %v", r)
		}
	}()
	c, _ := NewCodec("tcp")
	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x01, 0x81, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for 1-byte exception PDU")
	}
}

func TestExtraTCPLengthField(t *testing.T) {
	c, _ := NewCodec("tcp")
	// LEN=7: UID + 6-byte PDU (FC=3, byte count 4, data 00 0A 00 14).
	frame := []byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x07, 0x05, 0x03, 0x04, 0x00, 0x0A, 0x00, 0x14}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 4 || resp.Data[0] != 0x00 || resp.Data[1] != 0x0A {
		t.Errorf("data = %v, want [00 0A 00 14]", resp.Data)
	}
	// LEN=6 but frame too short: mismatch.
	if _, err := c.Decode(frame[:10]); err == nil {
		t.Error("expected error for length mismatch")
	}
	// LEN < 2 (no PDU): rejected.
	if _, err := c.Decode([]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x03, 0x00}); err == nil {
		t.Error("expected error for length < 2")
	}
}

func TestExtraTransactionIncrements(t *testing.T) {
	c, _ := NewCodec("tcp")
	req := &kernel.Request{Function: "read_coils", Address: "0", Count: 8}
	raw1, _ := c.Encode(req)
	raw2, _ := c.Encode(req)
	if id1, id2 := binary.BigEndian.Uint16(raw1[0:2]), binary.BigEndian.Uint16(raw2[0:2]); id2 != id1+1 {
		t.Errorf("transaction id went %d -> %d, want increment", id1, id2)
	}
}

func TestExtraRTUDecodesCorruptCRC(t *testing.T) {
	// GAP: decodeRTU never verifies the CRC; a frame with a corrupted CRC
	// is silently accepted. A real RTU device would reject it.
	c, _ := NewCodec("rtu")
	raw, _ := c.Encode(&kernel.Request{Function: "read_coils", Address: "0", Count: 8, Metadata: map[string]any{"unit_id": byte(1)}})
	raw[len(raw)-1] ^= 0xFF // corrupt CRC
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatalf("corrupt-CRC frame should not error under current code, got %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("data = %v, want empty", resp.Data)
	}
}

func TestExtraWriteSingleCoilOff(t *testing.T) {
	c, _ := NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{Function: "write_single_coil", Address: "10", Data: []byte{0x00}})
	if raw[10] != 0x00 || raw[11] != 0x00 {
		t.Errorf("coil OFF = %02X%02X, want 0000", raw[10], raw[11])
	}
}

func TestExtraDecodeUnknownFunction(t *testing.T) {
	c, _ := NewCodec("tcp")
	// LEN=2, UID=0, PDU = [0x2B]: unknown function passes raw PDU through.
	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x2B, 0x00}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0] != 0x2B {
		t.Errorf("data = %v, want [2B]", resp.Data)
	}
}

func TestExtraCRC16KnownVector(t *testing.T) {
	// Published CRC-16/ARC check value.
	if got := crc16([]byte("123456789")); got != 0x4B37 {
		t.Errorf("crc16(123456789) = 0x%04X, want 0x4B37", got)
	}
}
