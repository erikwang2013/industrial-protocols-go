// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package dnp3

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraCRC16DNPKnownVector(t *testing.T) {
	// Published CRC-16/DNP check value for "123456789" is 0xEA82.
	// BUG: crc16DNP returns 0x69FF because the working variable is XORed
	// with the full 16-bit CRC instead of being masked to 8 bits
	// (temp := crc ^ b, missing the spec's & 0xFF).
	if got := crc16DNP([]byte("123456789")); got != 0xEA82 {
		t.Errorf("crc16DNP(123456789) = 0x%04X, want 0xEA82", got)
	}
}

func TestExtraDecodeTPDUShortLength(t *testing.T) {
	// BUG: decodeTPDU computes apduEnd = 10 + (length - 5); a length < 5
	// makes apduEnd < apduStart and data[apduStart:apduEnd] panics
	// (dnp3.go decodeTPDU).
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decodeTPDU panicked on length < 5: %v", r)
		}
	}()
	c, _ := New().NewCodec("tcp")
	frame := []byte{TPDUStart1, TPDUStart2, 0x02, 0, 0, 0, 0, 0}
	crc := crc16DNP(frame[3:8])
	frame = append(frame, byte(crc), byte(crc>>8), 0x00, 0x00)
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected parse error for length < 5")
	}
}

func TestExtraDecodeTPDUTruncatedLength(t *testing.T) {
	// BUG: decodeTPDU slices data[apduStart:apduEnd] without checking
	// apduEnd against len(data); a length field that overruns the frame
	// panics (dnp3.go decodeTPDU).
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decodeTPDU panicked on truncated length: %v", r)
		}
	}()
	c, _ := New().NewCodec("tcp")
	frame := []byte{TPDUStart1, TPDUStart2, 0x0F, 0x44, 0, 0, 0, 0} // length 15
	crc := crc16DNP(frame[3:8])
	frame = append(frame, byte(crc), byte(crc>>8), 0xAB, 0xCD)
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected parse error for truncated length")
	}
}

func TestExtraEncodeDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["function"] != int(FuncRead) {
		t.Errorf("function = %v, want 1", resp.Metadata["function"])
	}
	if resp.Metadata["fir"] != true || resp.Metadata["fin"] != true {
		t.Errorf("fir/fin = %v/%v, want true/true", resp.Metadata["fir"], resp.Metadata["fin"])
	}
	if resp.Metadata["sequence"] != 1 {
		t.Errorf("sequence = %v, want 1", resp.Metadata["sequence"])
	}
}

func TestExtraSeqIncrements(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw1, _ := c.Encode(&kernel.Request{Function: "class0_poll"})
	raw2, _ := c.Encode(&kernel.Request{Function: "class0_poll"})
	if raw1[3]&0x3F != 1 || raw2[3]&0x3F != 2 {
		t.Errorf("sequence went %d -> %d, want 1 -> 2", raw1[3]&0x3F, raw2[3]&0x3F)
	}
}

func TestExtraConfirmFunction(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{Function: "confirm"})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["function"] != int(FuncConfirm) {
		t.Errorf("function = %v, want 0", resp.Metadata["function"])
	}
}

func TestExtraUnknownFunctionDefaultsClass0Poll(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{Function: "bogus"})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["function"] != int(FuncRead) {
		t.Errorf("function = %v, want 1 (class 0 poll fallback)", resp.Metadata["function"])
	}
}

func TestExtraFrameTooShort(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	if _, err := c.Decode(make([]byte, 9)); err == nil {
		t.Fatal("expected error for 9-byte frame")
	}
}

func TestExtraHeaderCRCByteOrder(t *testing.T) {
	// Header CRC is appended little-endian per DNP3.
	c, _ := New().NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{Function: "read"})
	want := crc16DNP(raw[3:8])
	if got := binary.LittleEndian.Uint16(raw[8:10]); got != want {
		t.Errorf("header CRC = 0x%04X, want 0x%04X", got, want)
	}
}
