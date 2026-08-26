// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package dali

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestEncodeCommandCodes(t *testing.T) {
	tests := []struct {
		fn   string
		want byte
	}{
		{"off", CmdOff},
		{"max", CmdMax},
		{"recall_min", CmdRecallMin},
		{"dim_up", CmdDimUp},
		{"dim_down", CmdDimDown},
		{"query_status", CmdQueryStatus},
	}
	c, _ := New().NewCodec("serial")
	for _, tt := range tests {
		raw, err := c.Encode(&kernel.Request{Function: tt.fn})
		if err != nil {
			t.Fatalf("Encode(%q) error: %v", tt.fn, err)
		}
		if raw[1] != tt.want {
			t.Errorf("Encode(%q) cmd = 0x%02X, want 0x%02X", tt.fn, raw[1], tt.want)
		}
	}
}

func TestEncodeAddressMetadata(t *testing.T) {
	c, _ := New().NewCodec("serial")
	req := &kernel.Request{
		Function: "off",
		Metadata: map[string]any{"address": byte(0x01)},
	}
	raw, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if raw[0] != 0x01 {
		t.Errorf("addr = 0x%02X, want 0x01", raw[0])
	}

	// float64 metadata is also accepted.
	req.Metadata["address"] = float64(0x02)
	raw, _ = c.Encode(req)
	if raw[0] != 0x02 {
		t.Errorf("addr = 0x%02X, want 0x02", raw[0])
	}
}

func TestEncodeBroadcastDefault(t *testing.T) {
	c, _ := New().NewCodec("serial")
	raw, err := c.Encode(&kernel.Request{Function: "max"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[0] != BroadcastAddr {
		t.Errorf("addr = 0x%02X, want 0xFE (broadcast)", raw[0])
	}
}

func TestEncodeDirectArcCmd(t *testing.T) {
	c, _ := New().NewCodec("serial")
	raw, err := c.Encode(&kernel.Request{Function: "direct_arc", Data: []byte{0x80}})
	if err != nil {
		t.Fatal(err)
	}
	if raw[1] != 0x80 {
		t.Errorf("cmd = 0x%02X, want 0x80", raw[1])
	}
}

func TestEncodeUnknownFunctionDefaultsOff(t *testing.T) {
	c, _ := New().NewCodec("serial")
	raw, err := c.Encode(&kernel.Request{Function: "bogus"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[1] != CmdOff {
		t.Errorf("cmd = 0x%02X, want 0x00", raw[1])
	}
}

func TestDecodeStatusBits(t *testing.T) {
	c, _ := New().NewCodec("serial")
	// Bit 0 = lamp failure, bit 1 = lamp on.
	resp, err := c.Decode([]byte{0x03})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["lamp_failure"] != true {
		t.Error("lamp_failure should be true")
	}
	if resp.Metadata["lamp_on"] != true {
		t.Error("lamp_on should be true")
	}

	resp, err = c.Decode([]byte{0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["lamp_failure"] != false || resp.Metadata["lamp_on"] != false {
		t.Error("bits 0/1 should be false for status 0x00")
	}
}

func TestDecodeEmpty(t *testing.T) {
	c, _ := New().NewCodec("serial")
	resp, err := c.Decode(nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data != nil {
		t.Errorf("data = %v, want nil", resp.Data)
	}
}
