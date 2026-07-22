// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package saej1850

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "saej1850" {
		t.Errorf("expected saej1850, got %s", p.Name())
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

func TestProtocolRejectsUnknownVariant(t *testing.T) {
	p := New()
	_, err := p.NewCodec("serial")
	if err == nil {
		t.Error("expected error for unknown variant")
	}
}

func TestEncodeMode01(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "mode01",
		Metadata: map[string]any{"pid": float64(0x0C)}, // PID $0C = RPM
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// Mode $01 + PID $0C
	// ISO 15765 single frame: byte 4=0x02 (length), byte 5=0x01 (mode), byte 6=0x0C (PID)
	if data[5] != 0x01 {
		t.Errorf("expected mode 0x01, got 0x%02X", data[5])
	}
	if data[6] != 0x0C {
		t.Errorf("expected PID 0x0C, got 0x%02X", data[6])
	}
}

func TestEncodeMode03(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "mode03"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	if data[5] != 0x03 {
		t.Errorf("expected mode 0x03, got 0x%02X", data[5])
	}
}

func TestEncodeMode0A(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "mode0A"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	if data[5] != 0x0A {
		t.Errorf("expected mode 0x0A, got 0x%02X", data[5])
	}
}

func TestEncodeDiagRequest(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "diag_request",
		Metadata: map[string]any{"mode": float64(0x09), "pid": float64(0x02)}, // Mode $09 PID $02 = VIN
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	if data[5] != 0x09 {
		t.Errorf("expected mode 0x09, got 0x%02X", data[5])
	}
	if data[6] != 0x02 {
		t.Errorf("expected PID 0x02 (VIN), got 0x%02X", data[6])
	}
}

func TestEncodeDiagResponse(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "diag_response",
		Data:     []byte{0x10, 0x14, 0x49, 0x02, 0x01, 0x57},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
}

func TestEncodeBroadcast(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "broadcast",
		Data:     []byte{0x01, 0x00},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
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

func TestDecodePhysicalRequest(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	j := c.(*j1850Codec)
	f := j.marshalFrame(j.buildPhysicalReqID(), []byte{0x02, 0x01, 0x00}, true)
	resp, err := c.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "physical_request" {
		t.Errorf("expected physical_request, got %v", resp.Metadata["type"])
	}
}

func TestDecodePhysicalResponse(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	j := c.(*j1850Codec)
	f := j.marshalFrame(j.buildPhysicalRespID(), []byte{0x06, 0x41, 0x00, 0xBE, 0x3F, 0xB8}, true)
	resp, err := c.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "physical_response" {
		t.Errorf("expected physical_response, got %v", resp.Metadata["type"])
	}
}

func TestDecodeFunctionalRequest(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	j := c.(*j1850Codec)
	f := j.marshalFrame(j.buildFunctionalID(), []byte{0x02, 0x01, 0x00}, true)
	resp, err := c.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "functional_request" {
		t.Errorf("expected functional_request, got %v", resp.Metadata["type"])
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
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	j := c.(*j1850Codec)
	id := j.buildPhysicalReqID()
	data := []byte{0x02, 0x01, 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00}

	f := j.marshalFrame(id, data, true)
	restored := j.unmarshalFrame(f)

	if restored.ID != id&0x7FFFFFFF {
		t.Errorf("ID mismatch: %x vs %x", restored.ID, id&0x7FFFFFFF)
	}
	if !restored.Ext {
		t.Error("expected extended frame flag")
	}
	for i := 0; i < 4; i++ {
		if restored.Data[i] != data[i] {
			t.Errorf("data byte %d mismatch: %x vs %x", i, restored.Data[i], data[i])
		}
	}
}

func TestPriority(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	// Default priority is 6
	j := c.(*j1850Codec)
	id := j.buildPhysicalReqID()
	pri := byte((id >> 26) & 0x07)
	if pri != 6 {
		t.Errorf("expected priority 6, got %d", pri)
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
