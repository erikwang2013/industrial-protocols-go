// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cclinkie

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "cclinkie" {
		t.Errorf("expected cclinkie, got %s", p.Name())
	}
	if p.DefaultPort() != 2005 {
		t.Errorf("expected 2005, got %d", p.DefaultPort())
	}
}

func TestNewCodecInvalid(t *testing.T) {
	_, err := New().NewCodec("ethernet")
	if err == nil {
		t.Error("expected error for unsupported variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "D100"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty output")
	}
}

func TestEncodeWrite(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Address:  "D100",
		Data:     []byte{0x11, 0x22},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 9 {
		t.Error("frame too short for write")
	}
}

func TestDecode(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	// Build a minimal valid frame: magic(CC1E) | cmd(01) | station(0001) | length(0000) | crc
	frame := []byte{0xCC, 0x1E, 0x01, 0x00, 0x01, 0x00, 0x00}
	crc := crc16CCITT(frame)
	frame = append(frame, byte(crc>>8), byte(crc&0xFF))

	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}

func TestDecodeShortFrame(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	_, err := c.Decode([]byte{0xCC})
	if err == nil {
		t.Error("expected error for short frame")
	}
}
