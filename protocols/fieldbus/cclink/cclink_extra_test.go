// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cclink

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraCRC16XMODEMKnownVector(t *testing.T) {
	// Published CRC-16/XMODEM check value.
	if got := crc16XMODEM([]byte("123456789")); got != 0x31C3 {
		t.Errorf("crc16XMODEM(123456789) = 0x%04X, want 0x31C3", got)
	}
}

func TestExtraEncodeWriteLayout(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	raw, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0xAB},
		Metadata: map[string]any{"station": float64(2)},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{StartDelim, 0x02, CtrlWrite, 0x01, 0xAB, 0xB8, 0xB4, StartDelim}
	if !bytes.Equal(raw, want) {
		t.Errorf("Encode = % X, want % X", raw, want)
	}
}

func TestExtraEncodeReadDefault(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	raw, err := c.Encode(&kernel.Request{Function: "read", Data: []byte{0x01}})
	if err != nil {
		t.Fatal(err)
	}
	if raw[0] != StartDelim || raw[1] != 0x01 || raw[2] != CtrlRead {
		t.Errorf("header = % X, want 7E 01 01", raw[:3])
	}
}

func TestExtraDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	raw, _ := c.Encode(&kernel.Request{Function: "write", Data: []byte{0xDE, 0xAD}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Data, []byte{0xDE, 0xAD}) {
		t.Errorf("data = % X, want DE AD", resp.Data)
	}
	if resp.Metadata["station"] != 1 || resp.Metadata["control"] != int(CtrlWrite) {
		t.Errorf("metadata = %v", resp.Metadata)
	}
	if resp.Address != "1" {
		t.Errorf("address = %q, want 1", resp.Address)
	}
}

func TestExtraDecodeErrors(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	if _, err := c.Decode(make([]byte, 6)); err == nil {
		t.Error("expected error for 6-byte frame")
	}
	if _, err := c.Decode([]byte{0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x7E}); err == nil {
		t.Error("expected error for missing start delimiter")
	}
	// dataLen = 5 but only 1 payload byte present.
	if _, err := c.Decode([]byte{0x7E, 0x01, 0x01, 0x05, 0xAB, 0x00, 0x00, 0x7E}); err == nil {
		t.Error("expected error for data length mismatch")
	}
	// Valid shape but wrong CRC.
	raw, _ := c.Encode(&kernel.Request{Function: "read", Data: []byte{0x01}})
	raw[len(raw)-2] ^= 0xFF
	if _, err := c.Decode(raw); err == nil {
		t.Error("expected error for CRC mismatch")
	}
}
