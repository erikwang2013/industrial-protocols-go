// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pci

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "pci" {
		t.Errorf("expected pci, got %s", p.Name())
	}
	vars := p.Variants()
	if len(vars) == 0 || vars[0] != "pci" {
		t.Errorf("expected first variant 'pci', got %v", vars)
	}
	if p.DefaultPort() != 0 {
		t.Errorf("expected default port 0, got %d", p.DefaultPort())
	}
	c, err := p.NewCodec("pci")
	if err != nil {
		t.Fatalf("NewCodec(pci) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(pci) returned nil codec")
	}

	// Verify it implements kernel.Codec
	var _ kernel.Codec = c
}

func TestCodecEncodeRoundTrip(t *testing.T) {
	c := &pciCodec{}
	original := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}

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
	c := &pciCodec{}
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

func TestCodecEncodePreservesData(t *testing.T) {
	c := &pciCodec{}
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i)
	}
	encoded, err := c.Encode(&kernel.Request{Function: "config_read", Data: payload})
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 256 {
		t.Errorf("expected 256 bytes, got %d", len(encoded))
	}
}

func TestProtocolAllowsAnyVariant(t *testing.T) {
	// The passthrough codec accepts any variant string since it just
	// forwards raw bytes to the PCI config space.
	p := New()
	for _, v := range []string{"", "pci", "pcie", "any"} {
		c, err := p.NewCodec(v)
		if err != nil {
			t.Errorf("NewCodec(%q) should not error: %v", v, err)
		}
		if c == nil {
			t.Errorf("NewCodec(%q) returned nil codec", v)
		}
	}
}

func TestPCIDriverInterface(t *testing.T) {
	var d *PCIDriver
	_ = d // PCIDriver type exists and compiles
}

func TestNewPCIDriverMissingDevice(t *testing.T) {
	// A non-existent BDF should produce an error.
	_, err := NewPCIDriver("0000:ff:1f.7")
	if err == nil {
		t.Skip("unexpected success — device may exist on this system")
	}
}
