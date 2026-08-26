// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package flexray

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// TestCRC16KnownVector checks CRC-16/XMODEM against the published check value.
func TestCRC16KnownVector(t *testing.T) {
	// CRC-16/XMODEM check value for "123456789" is 0x31C3.
	if got := crc16([]byte("123456789")); got != 0x31C3 {
		t.Errorf("crc16(\"123456789\") = 0x%04X, want 0x31C3", got)
	}
}

func TestCRC16Empty(t *testing.T) {
	if got := crc16(nil); got != 0x0000 {
		t.Errorf("crc16(nil) = 0x%04X, want 0x0000", got)
	}
}

// TestEncodeFrameCycleMetadata verifies cycle number and PPI flag are packed into the frame.
func TestEncodeFrameCycleMetadata(t *testing.T) {
	c, _ := New().NewCodec("can")
	req := &kernel.Request{
		Function: "frame",
		Data:     []byte{0x11, 0x22},
		Metadata: map[string]any{"cycle": float64(42), "ppi": true},
	}
	raw, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// CAN wire format: ID(4) + data(4); FlexRay payload starts at byte 4.
	cycle := binary.LittleEndian.Uint16(raw[4:6])
	if cycle != 42 {
		t.Errorf("cycle = %d, want 42", cycle)
	}
	status := raw[6]
	if status&0x80 == 0 {
		t.Error("PPI flag (bit 7 of status byte) not set")
	}
	if raw[7] != 0x11 {
		t.Errorf("first payload byte = 0x%02X, want 0x11", raw[7])
	}
}

func TestEncodeFrameNoPPI(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{Function: "frame", Data: []byte{1}})
	if err != nil {
		t.Fatal(err)
	}
	if raw[6]&0x80 != 0 {
		t.Error("PPI flag should be clear when metadata is absent")
	}
}

func TestEncodeFrameTruncatesLongPayload(t *testing.T) {
	c, _ := New().NewCodec("can")
	big := make([]byte, 300)
	for i := range big {
		big[i] = byte(i)
	}
	raw, err := c.Encode(&kernel.Request{Function: "frame", Data: big})
	if err != nil {
		t.Fatal(err)
	}
	// Only the first 4 payload bytes are carried in the CAN frame.
	if raw[7] != big[0] {
		t.Errorf("first data byte = 0x%02X, want 0x%02X", raw[7], big[0])
	}
}

func TestEncodeStatusFrame(t *testing.T) {
	c, _ := New().NewCodec("can")
	raw, err := c.Encode(&kernel.Request{Function: "status"})
	if err != nil {
		t.Fatal(err)
	}
	// CAN ID: 0x02000000 | slotID (1).
	if id := binary.LittleEndian.Uint32(raw[0:4]); id&0x7FFFFFFF != 0x02000001 {
		t.Errorf("status CAN ID = 0x%08X, want 0x02000001", id&0x7FFFFFFF)
	}
}

func TestDecodeFrameMetadata(t *testing.T) {
	c, _ := New().NewCodec("can")
	// Encode a frame with known slot (default 1) and cycle 7, then decode it back.
	raw, err := c.Encode(&kernel.Request{
		Function: "frame",
		Data:     []byte{0xAA},
		Metadata: map[string]any{"cycle": float64(7)},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["cycle"] != uint16(7) {
		t.Errorf("decoded cycle = %v, want 7", resp.Metadata["cycle"])
	}
	if resp.Metadata["slot_id"] != uint16(1) {
		t.Errorf("decoded slot_id = %v, want 1", resp.Metadata["slot_id"])
	}
	if resp.Metadata["ext"] != true {
		t.Error("decoded ext should be true")
	}
}

func TestMarshalCANExtFlag(t *testing.T) {
	// ID with the extended bit set is masked on unmarshal, flag preserved.
	f := unmarshalCAN(marshalCAN(bridge.CANFrame{ID: 0x1FFFFFFF, Data: []byte{1, 2, 3, 4}, Ext: true}))
	if !f.Ext {
		t.Error("ext flag lost in round trip")
	}
	if f.ID != 0x1FFFFFFF {
		t.Errorf("id = 0x%08X, want 0x1FFFFFFF", f.ID)
	}
}

func TestMarshalCANNonExt(t *testing.T) {
	f := unmarshalCAN(marshalCAN(bridge.CANFrame{ID: 0x1234567, Data: []byte{9, 8, 7, 6}}))
	if f.Ext {
		t.Error("ext flag should be false for standard frame")
	}
	if f.ID != 0x1234567 {
		t.Errorf("id = 0x%08X, want 0x1234567", f.ID)
	}
}
