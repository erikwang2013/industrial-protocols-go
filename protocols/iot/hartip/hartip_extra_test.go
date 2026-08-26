// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package hartip

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncodeHeaderLayout(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, err := c.Encode(&kernel.Request{Function: "read_unique_id"})
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != 13 {
		t.Fatalf("len = %d, want 13 (8 header + 5 HART payload)", len(raw))
	}
	if raw[0] != 1 {
		t.Errorf("version = %d, want 1", raw[0])
	}
	if raw[1] != 0 {
		t.Errorf("message type = %d, want 0 (request)", raw[1])
	}
	if got := binary.BigEndian.Uint16(raw[6:8]); got != 5 {
		t.Errorf("length = %d, want 5", got)
	}
	// HART payload begins at offset 8 with the delimiter.
	if raw[8] != 0x02 || raw[9] != 0x80 || raw[10] != 0x00 {
		t.Errorf("hart payload = % X, want 02 80 00", raw[8:11])
	}
}

func TestExtraEncodeCommand3(t *testing.T) {
	c, _ := New().NewCodec("udp")
	raw, err := c.Encode(&kernel.Request{Function: "read_dynamic_variables"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[10] != 3 { // header(8) + delimiter(1) + address(1)
		t.Errorf("hart cmd = %d, want 3", raw[10])
	}
}

func TestExtraDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{Function: "read_unique_id"})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["command"] != 0 {
		t.Errorf("command = %v, want 0", resp.Metadata["command"])
	}
	if resp.Metadata["long_frame"] != false {
		t.Errorf("long_frame = %v, want false", resp.Metadata["long_frame"])
	}
}

func TestExtraDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	if _, err := c.Decode(make([]byte, 7)); err == nil {
		t.Fatal("expected error for 7-byte frame")
	}
}

func TestExtraDecodeLengthMismatch(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// Length field says 10 but only 5 bytes follow the header.
	frame := []byte{0x01, 0x01, 0, 0, 0, 0, 0x00, 0x0A, 0x02, 0x80, 0x00, 0x00, 0x82}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for length mismatch")
	}
}
