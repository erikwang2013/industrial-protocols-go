// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package ethernetip

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestEncodeRegisterSessionLayout(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{Function: "register_session"})
	if binary.LittleEndian.Uint16(raw[0:2]) != CmdRegisterSession {
		t.Error("wrong command")
	}
	if binary.LittleEndian.Uint16(raw[2:4]) != 4 {
		t.Errorf("length = %d, want 4", binary.LittleEndian.Uint16(raw[2:4]))
	}
	if binary.LittleEndian.Uint32(raw[18:22]) != 1 {
		t.Errorf("options = %d, want 1", binary.LittleEndian.Uint32(raw[18:22]))
	}
	if len(raw) != 28 {
		t.Fatalf("len = %d, want 28", len(raw))
	}
	// Protocol version 1.0 at bytes 24-28.
	if raw[24] != 0x01 || raw[25] != 0x00 || raw[26] != 0x00 || raw[27] != 0x00 {
		t.Errorf("protocol version = % X, want 01 00 00 00", raw[24:28])
	}
}

func TestEncodeUnregisterSession(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, _ := c.Encode(&kernel.Request{
		Function: "unregister_session",
		Metadata: map[string]any{"session": float64(0xDEADBEEF)},
	})
	if binary.LittleEndian.Uint16(raw[0:2]) != CmdUnRegisterSession {
		t.Error("wrong command")
	}
	if binary.LittleEndian.Uint32(raw[4:8]) != 0xDEADBEEF {
		t.Errorf("session = 0x%08X, want 0xDEADBEEF", binary.LittleEndian.Uint32(raw[4:8]))
	}
	if len(raw) != 24 {
		t.Errorf("len = %d, want 24", len(raw))
	}
}

func TestEncodeReadTagLayout(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, err := c.Encode(&kernel.Request{
		Function: "read_tag",
		Address:  "TAG1",
		Metadata: map[string]any{"session": float64(7)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint16(raw[0:2]) != CmdSendRRData {
		t.Error("wrong command")
	}
	if binary.LittleEndian.Uint32(raw[4:8]) != 7 {
		t.Errorf("session = %d, want 7", binary.LittleEndian.Uint32(raw[4:8]))
	}
	if binary.LittleEndian.Uint16(raw[22:24]) != 1 {
		t.Errorf("item count = %d, want 1", binary.LittleEndian.Uint16(raw[22:24]))
	}
	if binary.LittleEndian.Uint16(raw[24:26]) != 0x00B2 {
		t.Errorf("item id = 0x%04X, want 0x00B2", binary.LittleEndian.Uint16(raw[24:26]))
	}
	// CIP service at payload offset 28.
	if raw[28] != CIPServiceReadTag {
		t.Errorf("CIP service = 0x%02X, want 0x4C", raw[28])
	}
	if raw[34] != byte(len("TAG1")) {
		t.Errorf("tag length = %d, want 4", raw[34])
	}
	if string(raw[35:35+4]) != "TAG1" {
		t.Errorf("tag = %q", raw[35:39])
	}
}

func TestEncodeReadTagInvokeIDIncrements(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw1, _ := c.Encode(&kernel.Request{Function: "read_tag", Address: "A"})
	raw2, _ := c.Encode(&kernel.Request{Function: "read_tag", Address: "B"})
	id1 := binary.LittleEndian.Uint16(raw1[30:32])
	id2 := binary.LittleEndian.Uint16(raw2[30:32])
	if id2 != id1+1 {
		t.Errorf("invoke id went %d -> %d, want increment", id1, id2)
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	if _, err := c.Encode(&kernel.Request{Function: "bogus"}); err == nil {
		t.Fatal("expected error for unknown function")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	if _, err := c.Decode(make([]byte, 23)); err == nil {
		t.Fatal("expected error for 23-byte frame")
	}
}

func TestDecodeRegisterSessionOK(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	frame := make([]byte, 24)
	binary.LittleEndian.PutUint16(frame[0:2], CmdRegisterSession)
	binary.LittleEndian.PutUint32(frame[4:8], 0x1234)
	binary.LittleEndian.PutUint32(frame[8:12], 0)
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data != nil {
		t.Errorf("data = %v, want nil for successful register session", resp.Data)
	}
	if resp.Metadata["session"] != uint32(0x1234) {
		t.Errorf("session = %v", resp.Metadata["session"])
	}
}

func TestDecodeErrorStatusIncludesData(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	frame := make([]byte, 26)
	binary.LittleEndian.PutUint16(frame[0:2], CmdSendRRData)
	binary.LittleEndian.PutUint32(frame[8:12], 0x100) // nonzero status
	frame[24] = 0xAB
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["status"] != uint32(0x100) {
		t.Errorf("status = %v", resp.Metadata["status"])
	}
	if len(resp.Data) != 2 || resp.Data[0] != 0xAB {
		t.Errorf("data = %v, want [AB 00]", resp.Data)
	}
}
