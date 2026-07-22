// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package dali

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "dali" {
		t.Errorf("expected dali, got %s", p.Name())
	}
	if p.DefaultPort() != 0 {
		t.Error("expected port 0")
	}
}

func TestEncodeOff(t *testing.T) {
	c, _ := New().NewCodec("serial")
	req := &kernel.Request{Function: "off", Metadata: map[string]any{"address": byte(0x01)}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 bytes, got %d", len(data))
	}
	if data[0] != 0x01 {
		t.Errorf("expected addr 0x01, got 0x%02X", data[0])
	}
	if data[1] != CmdOff {
		t.Errorf("expected CmdOff 0x00, got 0x%02X", data[1])
	}
}

func TestEncodeBroadcast(t *testing.T) {
	c, _ := New().NewCodec("serial")
	req := &kernel.Request{Function: "max"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if data[0] != BroadcastAddr {
		t.Errorf("expected broadcast addr 0xFE, got 0x%02X", data[0])
	}
	if data[1] != CmdMax {
		t.Errorf("expected CmdMax, got 0x%02X", data[1])
	}
}

func TestEncodeDirectArc(t *testing.T) {
	c, _ := New().NewCodec("serial")
	req := &kernel.Request{Function: "direct_arc", Data: []byte{128}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if data[1] != 128 {
		t.Errorf("expected 128, got %d", data[1])
	}
}

func TestDecodeStatus(t *testing.T) {
	c, _ := New().NewCodec("serial")
	resp, err := c.Decode([]byte{0x03})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["lamp_failure"] != true {
		t.Error("expected lamp_failure=true")
	}
	if resp.Metadata["lamp_on"] != true {
		t.Error("expected lamp_on=true")
	}
}
