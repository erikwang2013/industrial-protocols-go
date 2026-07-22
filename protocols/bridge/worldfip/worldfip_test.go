// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package worldfip

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "worldfip" {
		t.Errorf("expected worldfip, got %s", p.Name())
	}
	if p.DefaultPort() != 2007 {
		t.Errorf("expected 2007, got %d", p.DefaultPort())
	}
}

func TestNewCodecInvalid(t *testing.T) {
	_, err := New().NewCodec("serial")
	if err == nil {
		t.Error("expected error for unsupported variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "DAT1"})
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
		Address:  "DAT1",
		Data:     []byte{0x01, 0x02},
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
	// Build minimal valid frame: magic(57F1) | op(01) | tag(0000) | length(0000) | crc
	hdr := []byte{0x57, 0xF1, 0x01, 0x00, 0x00, 0x00, 0x00}
	crc := xorCRC16(hdr)
	frame := append(hdr, byte(crc>>8), byte(crc&0xFF))

	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}
