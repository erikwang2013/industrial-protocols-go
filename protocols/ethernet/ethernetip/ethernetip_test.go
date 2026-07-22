package ethernetip

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "ethernetip" {
		t.Errorf("expected ethernetip, got %s", p.Name())
	}
	if p.DefaultPort() != 44818 {
		t.Errorf("expected 44818, got %d", p.DefaultPort())
	}
}

func TestEncodeRegisterSession(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "register_session"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 28 {
		t.Fatalf("expected 28 bytes, got %d", len(data))
	}
	cmd := int(data[0]) | int(data[1])<<8
	if cmd != CmdRegisterSession {
		t.Errorf("expected 0x0065, got 0x%04X", cmd)
	}
}

func TestDecodeRegisterSessionResponse(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// Simulated RegisterSession response: success, session handle = 0x01020304
	resp := make([]byte, 28)
	// Use encoding/binary for LittleEndian
	resp[0] = 0x65
	resp[1] = 0x00 // CmdRegisterSession
	resp[2] = 0x04
	resp[3] = 0x00 // length: 4
	resp[4] = 0x01
	resp[5] = 0x02
	resp[6] = 0x03
	resp[7] = 0x04 // session handle

	result, err := c.Decode(resp)
	if err != nil {
		t.Fatal(err)
	}
	if result.Metadata["status"].(uint32) != 0 {
		t.Error("expected success status")
	}
	if result.Metadata["session"].(uint32) == 0 {
		t.Error("expected non-zero session handle")
	}
}

func TestEncodeReadTag(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{
		Function: "read_tag",
		Address:  "MyTag",
		Metadata: map[string]any{"session": float64(0x04030201)},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 28 {
		t.Errorf("expected at least 28 bytes, got %d", len(data))
	}
}
