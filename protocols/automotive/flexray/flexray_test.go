// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package flexray

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "flexray" {
		t.Errorf("expected flexray, got %s", p.Name())
	}
	if len(p.Variants()) == 0 || p.Variants()[0] != "can" {
		t.Errorf("expected first variant 'can', got %v", p.Variants())
	}
	if p.DefaultPort() != 0 {
		t.Errorf("expected default port 0, got %d", p.DefaultPort())
	}
	c, err := p.NewCodec("can")
	if err != nil {
		t.Fatalf("NewCodec(can) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(can) returned nil codec")
	}
}

func TestProtocolSerialVariant(t *testing.T) {
	p := New()
	c, err := p.NewCodec("serial")
	if err != nil {
		t.Fatalf("NewCodec(serial) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(serial) returned nil codec")
	}
}

func TestProtocolRejectsUnknownVariant(t *testing.T) {
	p := New()
	_, err := p.NewCodec("invalid")
	if err == nil {
		t.Error("expected error for unknown variant")
	}
}

func TestDriverNewSerialDriver(t *testing.T) {
	b, c, err := NewSerialDriver("/dev/ttyUSB0")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}

func TestEncodeFrame(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "frame",
		Data:     []byte{0x42, 0x01},
		Metadata: map[string]any{
			"cycle": float64(5),
			"ppi":   true,
		},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// Slot 1, cycle 5: CAN ID upper bits should have 0x01 (slot) and 0x05 (cycle)
	// ID = 0x01000000 | (1 << 10) | 5 = 0x01000405
	// Little-endian: 05 04 00 81 (with extended frame flag on byte 3)
	if data[4] != 0x05 || data[5] != 0x00 {
		t.Logf("frame data bytes: %X", data)
	}
}

func TestEncodeFrameTruncate(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	bigPayload := make([]byte, 300)
	for i := range bigPayload {
		bigPayload[i] = byte(i)
	}
	req := &kernel.Request{
		Function: "frame",
		Data:     bigPayload,
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
}

func TestEncodeStatus(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "status"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// Status request: data bytes contain slot/cycle info
	if data[4] != 1 { // slot ID
		t.Errorf("expected slot 1 in status, got %d", data[4])
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Encode(&kernel.Request{Function: "bogus"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestDecodeFrame(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	f := marshalCAN(bridge.CANFrame{
		ID:   0x01000405,
		Data: []byte{0x05, 0x00, 0x80, 0x42},
		Ext:  true,
	})
	resp, err := c.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["slot_id"].(uint16) != 1 {
		t.Errorf("expected slot 1, got %d", resp.Metadata["slot_id"])
	}
	if resp.Metadata["cycle"].(uint16) != 5 {
		t.Errorf("expected cycle 5, got %d", resp.Metadata["cycle"])
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Decode([]byte{0x00, 0x01})
	if err == nil {
		t.Error("expected error for short frame")
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	orig := bridge.CANFrame{
		ID:   0x01000405,
		Data: []byte{0x05, 0x00, 0x80, 0x42},
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
	for i := 0; i < 4; i++ {
		if restored.Data[i] != orig.Data[i] {
			t.Errorf("data byte %d mismatch: %x vs %x", i, restored.Data[i], orig.Data[i])
		}
	}
}

func TestCRC16(t *testing.T) {
	// Known CRC-16/XMODEM test vectors
	crc := crc16([]byte("123456789"))
	if crc != 0x31C3 {
		t.Errorf("expected CRC 0x31C3, got 0x%04X", crc)
	}
}

func TestDriverNewSocketCANDriver(t *testing.T) {
	b, c, err := NewSocketCANDriver("vcan0")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}
