// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package devicenet

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "devicenet" {
		t.Errorf("expected devicenet, got %s", p.Name())
	}
	if len(p.Variants()) < 2 {
		t.Errorf("expected at least 2 variants, got %d", len(p.Variants()))
	}
	if p.DefaultPort() != 2001 {
		t.Errorf("expected default port 2001, got %d", p.DefaultPort())
	}
}

func TestNewCodecCAN(t *testing.T) {
	p := New()
	c, err := p.NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("expected non-nil codec")
	}
}

func TestNewCodecGateway(t *testing.T) {
	p := New()
	c, err := p.NewCodec("gateway")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("expected non-nil gateway codec")
	}
}

func TestNewCodecUnknownVariant(t *testing.T) {
	p := New()
	_, err := p.NewCodec("serial")
	if err == nil {
		t.Error("expected error for unknown variant")
	}
}

func TestEncodeOpen(t *testing.T) {
	c, _ := New().NewCodec("can")
	data, err := c.Encode(&kernel.Request{Function: "open"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// Open command: first data byte is 0x4B
	if data[4] != 0x4B {
		t.Errorf("expected open command 0x4B, got 0x%02X", data[4])
	}
}

func TestEncodePoll(t *testing.T) {
	c, _ := New().NewCodec("can")
	req := &kernel.Request{
		Function: "poll",
		Data:     []byte{0xA5, 0x5A},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	if data[4] != 0xA5 || data[5] != 0x5A {
		t.Errorf("expected poll data A5 5A, got %02X %02X", data[4], data[5])
	}
}

func TestEncodePollTruncate(t *testing.T) {
	c, _ := New().NewCodec("can")
	req := &kernel.Request{
		Function: "poll",
		Data:     []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, _ := New().NewCodec("can")
	_, err := c.Encode(&kernel.Request{Function: "bogus"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestDecodeCANFrame(t *testing.T) {
	c, _ := New().NewCodec("can")
	f := marshalCAN(bridge.CANFrame{
		ID:   0x401,
		Data: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0x00, 0x00, 0x00, 0x00},
	})
	resp, err := c.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data[0] != 0xAA {
		t.Errorf("expected first byte 0xAA, got 0x%02X", resp.Data[0])
	}
	if resp.Metadata["can_id"].(uint32) != 0x401 {
		t.Errorf("expected can_id 0x401, got 0x%X", resp.Metadata["can_id"])
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("can")
	_, err := c.Decode([]byte{0x00, 0x01})
	if err == nil {
		t.Error("expected error for short frame")
	}
}

func TestGatewayCodecEncode(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	req := &kernel.Request{
		Function: "poll",
		Data:     []byte{0xAB, 0xCD},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	expected := "poll abcd\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestGatewayCodecDecode(t *testing.T) {
	c, _ := New().NewCodec("gateway")
	resp, err := c.Decode([]byte("OK 01\n"))
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != "OK 01\n" {
		t.Errorf("expected 'OK 01\\n', got %q", string(resp.Data))
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	orig := bridge.CANFrame{
		ID:   0x401,
		Data: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		Ext:  false,
	}
	buf := marshalCAN(orig)
	restored := unmarshalCAN(buf)
	if restored.ID != orig.ID {
		t.Errorf("ID mismatch: %x vs %x", restored.ID, orig.ID)
	}
	if restored.Ext != orig.Ext {
		t.Error("Ext flag mismatch")
	}
	for i := 0; i < 4; i++ {
		if restored.Data[i] != orig.Data[i] {
			t.Errorf("data byte %d mismatch: %x vs %x", i, restored.Data[i], orig.Data[i])
		}
	}
}

func TestMarshalRoundTripExtended(t *testing.T) {
	orig := bridge.CANFrame{
		ID:   0x80000000,
		Data: []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x00, 0x00, 0x00},
		Ext:  true,
	}
	buf := marshalCAN(orig)
	restored := unmarshalCAN(buf)
	if restored.ID != orig.ID&0x7FFFFFFF {
		t.Errorf("ID mismatch: %x vs %x", restored.ID, orig.ID&0x7FFFFFFF)
	}
	if !restored.Ext {
		t.Error("expected extended frame flag")
	}
}

func TestDriverNewSocketCANDriver(t *testing.T) {
	b, c, err := NewSocketCANDriver("vcan0", 0)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}

func TestDriverNewGatewayDriver(t *testing.T) {
	b, c, err := NewGatewayDriver("localhost:2001")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}
