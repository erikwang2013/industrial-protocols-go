// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package interbus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "interbus" {
		t.Errorf("expected interbus, got %s", p.Name())
	}
	if p.DefaultPort() != 2006 {
		t.Errorf("expected 2006, got %d", p.DefaultPort())
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
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "0x100"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty output")
	}
}

func TestDecode(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	// Build minimal valid frame: magic(1B55) | cmd(01) | slave(01) | length(0000) | checksum
	hdr := []byte{0x1B, 0x55, 0x01, 0x01, 0x00, 0x00}
	checksum := checksum8(hdr)
	frame := append(hdr, checksum)

	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}
