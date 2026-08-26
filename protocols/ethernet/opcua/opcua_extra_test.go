// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package opcua

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraHELMessageSizeAtOffset4(t *testing.T) {
	// BUG: the OPC UA HEL message size must sit at offset 4 (spec), but
	// encodeHello writes the placeholder at offset 4 and then patches the
	// size into offset 8 (the protocol-version slot), leaving bytes 4:8 = 0.
	c, _ := New().NewCodec("binary")
	raw, err := c.Encode(&kernel.Request{Function: "hel"})
	if err != nil {
		t.Fatal(err)
	}
	if got := binary.LittleEndian.Uint32(raw[4:8]); got != uint32(len(raw)) {
		t.Errorf("message size at offset 4 = %d, want %d", got, len(raw))
	}
}

func TestExtraHELDefaultLayout(t *testing.T) {
	c, _ := New().NewCodec("binary")
	raw, err := c.Encode(&kernel.Request{Function: "hel"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw[:4]) != "HELF" {
		t.Errorf("type = %q, want HELF", raw[:4])
	}
	// Size field at offset 4 is asserted by TestExtraHELMessageSizeAtOffset4;
	// bytes 8:12 hold the protocol version (0).
	if got := binary.LittleEndian.Uint32(raw[8:12]); got != 0 {
		t.Errorf("protocol version at offset 8 = %d, want 0", got)
	}
	const endpoint = "opc.tcp://localhost:4840"
	if len(raw) != 36+len(endpoint) {
		t.Fatalf("len = %d, want %d", len(raw), 36+len(endpoint))
	}
	if string(raw[36:]) != endpoint {
		t.Errorf("endpoint = %q, want %q", raw[36:], endpoint)
	}
}

func TestExtraHELCustomEndpoint(t *testing.T) {
	c, _ := New().NewCodec("binary")
	raw, err := c.Encode(&kernel.Request{
		Function: "hel",
		Metadata: map[string]any{"endpoint": "opc.tcp://server:4840"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw[36:]) != "opc.tcp://server:4840" {
		t.Errorf("endpoint = %q", raw[36:])
	}
}

func TestExtraOPNLayout(t *testing.T) {
	c, _ := New().NewCodec("binary")
	raw, err := c.Encode(&kernel.Request{Function: "open_secure_channel"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw[:3]) != "OPN" {
		t.Errorf("type = %q, want OPN", raw[:3])
	}
	if len(raw) != 40 {
		t.Errorf("len = %d, want 40", len(raw))
	}
	// BUG: the OPN size field at offset 4 is written as a placeholder but
	// never patched (unlike encodeHello), leaving it at 0.
	if got := binary.LittleEndian.Uint32(raw[4:8]); got != 40 {
		t.Errorf("size = %d, want 40", got)
	}
}

func TestExtraDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("binary")
	if _, err := c.Decode([]byte{0x41, 0x43, 0x4B}); err == nil {
		t.Fatal("expected error for 3-byte frame")
	}
}

func TestExtraDecodeUnknownMessageType(t *testing.T) {
	c, _ := New().NewCodec("binary")
	resp, err := c.Decode([]byte("XYZ01234567"))
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != "XYZ01234567" {
		t.Errorf("data = %q, want passthrough", resp.Data)
	}
}
