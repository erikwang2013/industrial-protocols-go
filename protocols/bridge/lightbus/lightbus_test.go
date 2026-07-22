package lightbus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "lightbus" {
		t.Errorf("expected lightbus, got %s", p.Name())
	}
	if p.DefaultPort() != 2008 {
		t.Errorf("expected 2008, got %d", p.DefaultPort())
	}
}

func TestNewCodecInvalid(t *testing.T) {
	_, err := New().NewCodec("fiber")
	if err == nil {
		t.Error("expected error for unsupported variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "0x10"})
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
	// Build minimal valid frame: magic(4C42) | cmd(01) | module(01) | channel(0000) | length(0000) | checksum
	hdr := []byte{0x4C, 0x42, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00}
	checksum := checksum8XOR(hdr)
	frame := append(hdr, checksum)

	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}
