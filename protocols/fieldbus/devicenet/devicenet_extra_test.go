// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package devicenet

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraOpenFrameLayout(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{Function: "open"})
	if err != nil {
		t.Fatal(err)
	}
	f := unmarshalCAN(raw)
	if !f.Ext {
		t.Error("open frame must use extended addressing")
	}
	// The extended bit is stripped on unmarshal; with macID 0 the Group-3
	// message ID 6 bits are all zero, leaving ID 0.
	if f.ID != 0 {
		t.Errorf("can id = 0x%08X, want 0x00000000", f.ID)
	}
	// Unmasked wire ID keeps the extended bit.
	if raw[3]&0x80 == 0 {
		t.Error("wire ID must have bit 7 of byte 3 set (extended frame)")
	}
	// marshalCAN carries the first 4 data bytes on the wire.
	if !bytes.Equal(raw[4:8], []byte{0x4B, 0x03, 0x01, 0x01}) {
		t.Errorf("data = % X, want 4B 03 01 01", raw[4:8])
	}
}

func TestExtraPollPadAndID(t *testing.T) {
	c := &dnCodec{macID: 3}
	raw, err := c.Encode(&kernel.Request{Function: "poll", Data: []byte{0xAA, 0xBB, 0xCC}})
	if err != nil {
		t.Fatal(err)
	}
	f := unmarshalCAN(raw)
	if f.ID != 0x403 {
		t.Errorf("can id = 0x%X, want 0x403", f.ID)
	}
	// Poll data is padded to 8 in the frame, then truncated to 4 by the wire format.
	if !bytes.Equal(raw[4:8], []byte{0xAA, 0xBB, 0xCC, 0}) {
		t.Errorf("data = % X, want AA BB CC 00", raw[4:8])
	}
}

func TestExtraPollTruncatesTo8(t *testing.T) {
	c := &dnCodec{macID: 0}
	raw, _ := c.Encode(&kernel.Request{Function: "poll", Data: make([]byte, 12)})
	if len(raw) != 8 {
		t.Errorf("frame len = %d, want 8", len(raw))
	}
}

func TestExtraGatewayCodecEncode(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	raw, err := c.Encode(&kernel.Request{Function: "read", Data: []byte{0x01, 0x02}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "read 0102\n" {
		t.Errorf("Encode = %q, want %q", raw, "read 0102\n")
	}
}

func TestExtraDecodeMetadata(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, _ := c.Encode(&kernel.Request{Function: "poll", Data: []byte{1}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["can_id"] != uint32(0x400) {
		t.Errorf("can_id = %v, want 0x400", resp.Metadata["can_id"])
	}
	if resp.Metadata["ext"] != false {
		t.Errorf("ext = %v, want false", resp.Metadata["ext"])
	}
}

func TestExtraDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("can")
	if _, err := c.Decode(make([]byte, 3)); err == nil {
		t.Fatal("expected error for 3-byte frame")
	}
}
