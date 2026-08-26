// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package canopen

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

func TestExtraSDOReadDefaults(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{Function: "sdo_read"})
	if err != nil {
		t.Fatal(err)
	}
	f := unmarshalCAN(raw)
	if f.ID != 0x601 {
		t.Errorf("can id = 0x%X, want 0x601", f.ID)
	}
	if f.Data[0] != 0x40 {
		t.Errorf("ccs = 0x%02X, want 0x40", f.Data[0])
	}
	// Default index 0x1000 little-endian, sub-index 0.
	if idx := binary.LittleEndian.Uint16(f.Data[1:3]); idx != 0x1000 {
		t.Errorf("index = 0x%04X, want 0x1000", idx)
	}
	if f.Data[3] != 0 {
		t.Errorf("sub = %d, want 0", f.Data[3])
	}
}

func TestExtraSDOReadMetadata(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{
		Function: "sdo_read",
		Metadata: map[string]any{"index": float64(0x1017), "sub": float64(1)},
	})
	if err != nil {
		t.Fatal(err)
	}
	f := unmarshalCAN(raw)
	if idx := binary.LittleEndian.Uint16(f.Data[1:3]); idx != 0x1017 {
		t.Errorf("index = 0x%04X, want 0x1017", idx)
	}
	if f.Data[3] != 1 {
		t.Errorf("sub = %d, want 1", f.Data[3])
	}
}

func TestExtraSDOWriteLayout(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{Function: "sdo_write", Data: []byte{0x11, 0x22, 0x33, 0x44}})
	if err != nil {
		t.Fatal(err)
	}
	f := unmarshalCAN(raw)
	if f.Data[0] != 0x23 {
		t.Errorf("cs = 0x%02X, want 0x23", f.Data[0])
	}
	if !bytes.Equal(f.Data[1:], []byte{0x00, 0x10, 0x00}) {
		t.Errorf("index/sub = % X, want 00 10 00 (default 0x1000 sub 0)", f.Data[1:])
	}
	// GAP: marshalCAN stores only the first 4 data bytes; the SDO payload
	// (command + index + sub fill the slot) is dropped from the wire frame.
	if !bytes.Equal(raw[4:8], []byte{0x23, 0x00, 0x10, 0x00}) {
		t.Errorf("wire data = % X, want 23 00 10 00", raw[4:8])
	}
}

func TestExtraSDOWritePayloadLostOnWire(t *testing.T) {
	// GAP: sdo_write payload bytes never reach the 4-byte wire slot; a
	// 4-byte and a 6-byte payload produce identical frames.
	c, _ := New().NewCodec("can")
	raw4, _ := c.Encode(&kernel.Request{Function: "sdo_write", Data: []byte{0xAA, 0xBB, 0xCC, 0xDD}})
	raw6, _ := c.Encode(&kernel.Request{Function: "sdo_write", Data: []byte{1, 2, 3, 4, 5, 6}})
	if !bytes.Equal(raw4, raw6) {
		t.Errorf("frames differ: % X vs % X (payload should be dropped identically)", raw4, raw6)
	}
}

func TestExtraNMTCommandCodes(t *testing.T) {
	c, _ := New().NewCodec("can")
	tests := []struct {
		fn   string
		want byte
	}{
		{"nmt_start", 0x01},
		{"nmt_stop", 0x02},
		{"nmt_reset", 0x82},
	}
	for _, tt := range tests {
		raw, err := c.Encode(&kernel.Request{Function: tt.fn})
		if err != nil {
			t.Fatalf("Encode(%q) error: %v", tt.fn, err)
		}
		f := unmarshalCAN(raw)
		if f.ID != 0x000 {
			t.Errorf("nmt can id = 0x%X, want 0x000", f.ID)
		}
		if f.Data[0] != tt.want || f.Data[1] != 1 {
			t.Errorf("nmt data = % X, want %02X 01", f.Data, tt.want)
		}
	}
}

func TestExtraHeartbeatFrame(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{Function: "heartbeat"})
	if err != nil {
		t.Fatal(err)
	}
	f := unmarshalCAN(raw)
	if f.ID != 0x701 {
		t.Errorf("can id = 0x%X, want 0x701", f.ID)
	}
	if f.Data[0] != 0x05 {
		t.Errorf("state = 0x%02X, want 0x05 (operational)", f.Data[0])
	}
}

func TestExtraDecodeShortFramePanic(t *testing.T) {
	// BUG: Decode only rejects frames < 4 bytes, but unmarshalCAN slices
	// data[4:8]; frames of 4..7 bytes panic (canopen.go Decode/unmarshalCAN).
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Decode panicked on 7-byte frame: %v", r)
		}
	}()
	c, _ := New().NewCodec("can")
	if _, err := c.Decode([]byte{0x01, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00}); err == nil {
		t.Fatal("expected parse error for 7-byte frame")
	}
}

func TestExtraDecodeSDOAbort(t *testing.T) {
	c, _ := New().NewCodec("can")
	f := bridge.CANFrame{ID: 0x601, Data: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	_, err := c.Decode(marshalCAN(f))
	if err == nil {
		t.Fatal("expected SDO abort error")
	}
	var pe *kernel.ProtocolError
	if !errors.As(err, &pe) {
		t.Fatalf("expected ProtocolError, got %T", err)
	}
	if pe.Message != "SDO abort" {
		t.Errorf("message = %q, want SDO abort", pe.Message)
	}
}

func TestExtraDecodeBootupMetadata(t *testing.T) {
	c, _ := New().NewCodec("can")
	f := bridge.CANFrame{ID: 0x701, Data: []byte{0x05}}
	resp, err := c.Decode(marshalCAN(f))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "bootup" {
		t.Errorf("type = %v, want bootup", resp.Metadata["type"])
	}
	if resp.Metadata["can_id"] != uint32(0x701) {
		t.Errorf("can_id = %v", resp.Metadata["can_id"])
	}
}
