// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package profinet

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraDecodeBlockLenPanic(t *testing.T) {
	// BUG: Decode slices data[10:10+blockLen] without bounds checking; a
	// frame whose blockLen exceeds the remaining bytes panics (profinet.go).
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Decode panicked on oversized blockLen: %v", r)
		}
	}()
	c, _ := New().NewCodec("nrt")
	frame := []byte{0xFE, 0xFE, 0x04, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x20} // blockLen = 32
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected parse error for oversized blockLen")
	}
}

func TestExtraEncodeDCPSetLayout(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	raw, err := c.Encode(&kernel.Request{
		Function: "dcp_set",
		Metadata: map[string]any{"device_name": "plc1", "ip_address": "192.168.1.10"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if binary.BigEndian.Uint16(raw[0:2]) != FrameIDDCPSet {
		t.Errorf("frame id = 0x%04X, want 0xFEFE", binary.BigEndian.Uint16(raw[0:2]))
	}
	if raw[2] != ServiceIDSet || raw[3] != ServiceTypeRequest {
		t.Errorf("service = %d/%d, want 4/0", raw[2], raw[3])
	}
	if xid := binary.BigEndian.Uint32(raw[4:8]); xid != 2 {
		t.Errorf("xid = %d, want 2 (starts at 1, incremented per encode)", xid)
	}
	if bl := binary.BigEndian.Uint16(raw[8:10]); bl != 24 {
		t.Errorf("block len = %d, want 24", bl)
	}
	// Name block: [02 00 02 04 p l c 1]
	wantName := []byte{0x02, 0x00, 0x02, 0x04, 'p', 'l', 'c', '1'}
	if !bytes.Equal(raw[10:18], wantName) {
		t.Errorf("name block = % X, want % X", raw[10:18], wantName)
	}
	// IP block: [01 00 0E 0C C0 A8 01 0A FF FF FF 00 C0 A8 01 01]
	wantIP := []byte{0x01, 0x00, 0x0E, 0x0C, 192, 168, 1, 10, 255, 255, 255, 0, 192, 168, 1, 1}
	if !bytes.Equal(raw[18:34], wantIP) {
		t.Errorf("ip block = % X, want % X", raw[18:34], wantIP)
	}
}

func TestExtraDCPSetDefaults(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	raw, err := c.Encode(&kernel.Request{Function: "dcp_set"})
	if err != nil {
		t.Fatal(err)
	}
	// Default name "device" with a zeros IP block.
	wantName := []byte{0x02, 0x00, 0x02, 0x06, 'd', 'e', 'v', 'i', 'c', 'e'}
	if !bytes.Equal(raw[10:20], wantName) {
		t.Errorf("default name block = % X, want % X", raw[10:20], wantName)
	}
	if !bytes.Equal(raw[20:24], []byte{0x01, 0x00, 0x0E, 0x0C}) {
		t.Errorf("ip block header = % X, want 01 00 0E 0C", raw[20:24])
	}
	if !bytes.Equal(raw[24:36], make([]byte, 12)) {
		t.Errorf("ip block address = % X, want 12 zeros", raw[24:36])
	}
}

func TestExtraBuildNameBlock(t *testing.T) {
	if got := buildNameBlock(""); !bytes.Equal(got, []byte{0x02, 0x00, 0x02, 0x06, 'd', 'e', 'v', 'i', 'c', 'e'}) {
		t.Errorf("buildNameBlock(\"\") = % X", got)
	}
	if got := buildNameBlock("x"); !bytes.Equal(got, []byte{0x02, 0x00, 0x02, 0x01, 'x'}) {
		t.Errorf("buildNameBlock(x) = % X", got)
	}
}

func TestExtraBuildIPBlock(t *testing.T) {
	empty := buildIPBlock("")
	if len(empty) != 16 || !bytes.Equal(empty[4:], make([]byte, 12)) {
		t.Errorf("buildIPBlock(\"\") = % X, want 4-byte header + 12 zeros", empty)
	}
	got := buildIPBlock("10.0.0.5")
	want := []byte{0x01, 0x00, 0x0E, 0x0C, 10, 0, 0, 5, 255, 255, 255, 0, 10, 0, 0, 1}
	if !bytes.Equal(got, want) {
		t.Errorf("buildIPBlock = % X, want % X", got, want)
	}
}

func TestExtraRecordReadLayout(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	raw, err := c.Encode(&kernel.Request{Function: "record_read", Metadata: map[string]any{"index": float64(5)}})
	if err != nil {
		t.Fatal(err)
	}
	if binary.BigEndian.Uint16(raw[0:2]) != FrameIDRecordData {
		t.Errorf("frame id = 0x%04X, want 0xFEFC", binary.BigEndian.Uint16(raw[0:2]))
	}
	if raw[2] != 0x01 {
		t.Errorf("service = %d, want 1 (RecordDataRead)", raw[2])
	}
	if bl := binary.BigEndian.Uint16(raw[8:10]); bl != 6 {
		t.Errorf("block len = %d, want 6", bl)
	}
	if got := binary.BigEndian.Uint32(raw[12:16]); got != 5 {
		t.Errorf("index = %d, want 5", got)
	}
}

func TestExtraDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	if _, err := c.Decode(make([]byte, 9)); err == nil {
		t.Fatal("expected error for 9-byte frame")
	}
}

func TestExtraDecodeDCPResponseMetadata(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	frame := []byte{0xFE, 0xFD, 0x05, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x01, 0x02}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["frame_id"] != int(FrameIDDCPIdentify) {
		t.Errorf("frame_id = %v, want %d", resp.Metadata["frame_id"], FrameIDDCPIdentify)
	}
	if resp.Metadata["service_id"] != ServiceIDIdentify {
		t.Errorf("service_id = %v, want %d", resp.Metadata["service_id"], ServiceIDIdentify)
	}
	if resp.Metadata["service_type"] != ServiceTypeResponse {
		t.Errorf("service_type = %v, want %d", resp.Metadata["service_type"], ServiceTypeResponse)
	}
	if resp.Metadata["xid"] != uint32(1) {
		t.Errorf("xid = %v, want 1", resp.Metadata["xid"])
	}
	if len(resp.Data) != 2 || resp.Data[0] != 0x01 || resp.Data[1] != 0x02 {
		t.Errorf("data = %v, want [01 02]", resp.Data)
	}
}
