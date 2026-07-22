package modbusplus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "modbusplus" {
		t.Errorf("expected modbusplus, got %s", p.Name())
	}
	if p.DefaultPort() != 2010 {
		t.Errorf("expected 2010, got %d", p.DefaultPort())
	}
}

func TestNewCodecCmd(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatalf("NewCodec(cmd) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(cmd) returned nil codec")
	}
}

func TestNewCodecInvalid(t *testing.T) {
	_, err := New().NewCodec("invalid")
	if err == nil {
		t.Error("expected error for unsupported variant")
	}
}

func TestDriverNewCmdDriver(t *testing.T) {
	b, c, err := NewCmdDriver("/usr/bin/sa85_cli")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "40001"})
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
		Address:  "40001",
		Data:     []byte{0x12, 0x34},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 8 {
		t.Error("frame too short for write")
	}
}

func TestDecode(t *testing.T) {
	c, err := New().NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	// Build minimal valid frame: magic(4D42) | dest(01) | cmd(01) | length(0000) | crc
	hdr := []byte{0x4D, 0x42, 0x01, 0x01, 0x00, 0x00}
	crc := modbusCRC16(hdr)
	frame := append(hdr, byte(crc&0xFF), byte(crc>>8))

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
	_, err := c.Decode([]byte{0x4D})
	if err == nil {
		t.Error("expected error for short frame")
	}
}
