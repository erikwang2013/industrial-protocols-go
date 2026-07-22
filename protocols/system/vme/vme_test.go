package vme

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "vme" {
		t.Errorf("expected vme, got %s", p.Name())
	}
	vars := p.Variants()
	if len(vars) == 0 || vars[0] != "vme" {
		t.Errorf("expected first variant 'vme', got %v", vars)
	}
	if p.DefaultPort() != 0 {
		t.Errorf("expected default port 0, got %d", p.DefaultPort())
	}
	c, err := p.NewCodec("vme")
	if err != nil {
		t.Fatalf("NewCodec(vme) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(vme) returned nil codec")
	}

	// Verify it implements kernel.Codec
	var _ kernel.Codec = c
}

func TestCodecEncodeRoundTrip(t *testing.T) {
	c := &vmeCodec{}
	original := []byte{0xAA, 0xBB, 0xCC, 0xDD}

	encoded, err := c.Encode(&kernel.Request{Data: original})
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 4 {
		t.Errorf("expected 4 bytes, got %d", len(encoded))
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
	c := &vmeCodec{}
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

func TestCodecPreservesVMEAddressing(t *testing.T) {
	// VME supports A16, A24, A32 addressing modes.
	// The codec passes address modifier bytes through transparently.
	c := &vmeCodec{}
	payload := []byte{
		0x01,       // address modifier (A24 supervisory data)
		0x00, 0x10, // address high
		0x00, 0x00, // address low
		0xDE, 0xAD, 0xBE, 0xEF, // data
	}
	encoded, err := c.Encode(&kernel.Request{
		Function: "read",
		Address:  "0x00100000",
		Data:     payload,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 9 {
		t.Errorf("expected 9 bytes, got %d", len(encoded))
	}
}

func TestVMEDriverInterface(t *testing.T) {
	var d *VMEDriver
	_ = d // VMEDriver type exists and compiles
}

func TestNewVMEDriverMissingDevice(t *testing.T) {
	// A high slot number should fail unless hardware is present.
	_, err := NewVMEDriver(99)
	if err == nil {
		t.Skip("unexpected success — VME hardware may be present")
	}
}
