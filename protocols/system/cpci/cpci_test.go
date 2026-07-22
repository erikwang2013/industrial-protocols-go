// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cpci

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "cpci" {
		t.Errorf("expected cpci, got %s", p.Name())
	}
	vars := p.Variants()
	if len(vars) == 0 || vars[0] != "cpci" {
		t.Errorf("expected first variant 'cpci', got %v", vars)
	}
	if p.DefaultPort() != 0 {
		t.Errorf("expected default port 0, got %d", p.DefaultPort())
	}
	c, err := p.NewCodec("cpci")
	if err != nil {
		t.Fatalf("NewCodec(cpci) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(cpci) returned nil codec")
	}

	// Verify it implements kernel.Codec
	var _ kernel.Codec = c
}

func TestCodecEncodeRoundTrip(t *testing.T) {
	c := &cpciCodec{}
	original := []byte{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}

	encoded, err := c.Encode(&kernel.Request{Data: original})
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 8 {
		t.Errorf("expected 8 bytes, got %d", len(encoded))
	}

	decoded, err := c.Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded.Data) != len(original) {
		t.Errorf("expected %d bytes, got %d", len(original), len(decoded.Data))
	}
	for i := range original {
		if decoded.Data[i] != original[i] {
			t.Errorf("byte %d: expected %02x, got %02x", i, original[i], decoded.Data[i])
		}
	}
}

func TestCodecDecodeEmpty(t *testing.T) {
	c := &cpciCodec{}
	resp, err := c.Decode([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected empty data, got %v", resp.Data)
	}
}

func TestCodecCPCIConfigSpaceAccess(t *testing.T) {
	// CompactPCI config space is 256 bytes for Type 0, 4096 for Type 1.
	// The codec passes through the full config space transparently.
	c := &cpciCodec{}

	// Simulate a config space read of vendor/device ID
	configRead := []byte{0x00, 0x00, 0x00, 0x00} // offset 0
	encoded, err := c.Encode(&kernel.Request{
		Function: "read",
		Address:  "0x00",
		Count:    4,
		Data:     configRead,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 4 {
		t.Errorf("expected 4 bytes, got %d", len(encoded))
	}
}

func TestCPCIDriverInterface(t *testing.T) {
	var d *CPCIDriver
	_ = d // CPCIDriver type exists and compiles
}

func TestNewCPCIDriverMissingDevice(t *testing.T) {
	// A non-existent BDF should produce an error.
	_, err := NewCPCIDriver("0000:ff:1f.7")
	if err == nil {
		t.Skip("unexpected success — device may exist on this system")
	}
}
