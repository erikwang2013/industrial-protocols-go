// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package hart

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "hart" {
		t.Errorf("expected hart, got %s", p.Name())
	}
	if len(p.Variants()) != 1 {
		t.Error("expected 1 variant")
	}
	if p.DefaultPort() != 0 {
		t.Error("expected port 0")
	}

	// Verify valid variant creates a codec
	c, err := p.NewCodec("fsk")
	if err != nil {
		t.Fatalf("fsk variant should be valid: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil codec")
	}

	// Verify invalid variant returns error
	_, err = p.NewCodec("invalid")
	if err == nil {
		t.Error("expected error for invalid variant")
	}
}

func TestEncodeCommand0(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	req := &kernel.Request{Function: "read_unique_id"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 10 {
		t.Fatalf("expected 10 bytes, got %d", len(data))
	}
	if data[5] != 0x02 {
		t.Error("expected delimiter 0x02 (master→slave, short frame)")
	}
	if data[7] != 0 {
		t.Error("expected command 0")
	}
}

func TestEncodeCommand3(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	req := &kernel.Request{Function: "read_dynamic_variables"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 10 {
		t.Fatalf("expected 10 bytes, got %d", len(data))
	}
	if data[7] != 3 {
		t.Error("expected command 3")
	}
}

func TestEncodeWithPollingAddr(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	req := &kernel.Request{
		Function: "read_unique_id",
		Metadata: map[string]any{"polling_addr": byte(5)},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// addr should be 5 | 0x80 = 0x85
	if data[6] != 0x85 {
		t.Errorf("expected addr 0x85, got 0x%02X", data[6])
	}
}

func TestDecodeCommand0Response(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	// Short frame: 5 preambles + delim 0x06 + addr + cmd 0 + byteCount 11 + 11 data bytes
	// Build manually then set checksum
	resp := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x06, 0x80, 0, 11, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	// Index 20 will be checksum
	resp = append(resp, 0x00)
	resp[len(resp)-1] = xorChecksum(resp[5 : len(resp)-1])

	result, err := c.Decode(resp)
	if err != nil {
		t.Fatal(err)
	}
	cmd, ok := result.Metadata["command"].(int)
	if !ok || cmd != 0 {
		t.Errorf("expected command 0, got %v", cmd)
	}

	longFrame, _ := result.Metadata["long_frame"].(bool)
	if longFrame {
		t.Error("expected short frame")
	}

	mfr, ok := result.Metadata["manufacturer"].(byte)
	if !ok || mfr != 1 {
		t.Errorf("expected manufacturer 1, got %v", mfr)
	}

	if len(result.Data) != 11 {
		t.Errorf("expected 11 data bytes, got %d", len(result.Data))
	}
}

func TestDecodeLongFrame(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	// Long frame: delim 0x86 + 5-byte addr + cmd + byteCount + data
	// Build: [preamble*5, delim, addr[5], cmd, byteCount, data[5]]
	resp := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x86, 0x80, 0x01, 0x02, 0x03, 0x04, 0, 5, 10, 20, 30, 40, 50}
	// Append space for checksum and compute it
	resp = append(resp, 0x00)
	resp[len(resp)-1] = xorChecksum(resp[5 : len(resp)-1])

	result, err := c.Decode(resp)
	if err != nil {
		t.Fatal(err)
	}
	longFrame, _ := result.Metadata["long_frame"].(bool)
	if !longFrame {
		t.Error("expected long frame")
	}
	cmd, _ := result.Metadata["command"].(int)
	if cmd != 0 {
		t.Errorf("expected command 0, got %d", cmd)
	}
	if len(result.Data) != 5 {
		t.Errorf("expected 5 data bytes, got %d", len(result.Data))
	}
}

func TestChecksumMismatch(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	resp := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x06, 0x80, 0, 0, 0x00}
	_, err := c.Decode(resp)
	if err == nil {
		t.Error("expected checksum mismatch error")
	}
}

func TestFrameTooShort(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	_, err := c.Decode([]byte{0xFF, 0xFF})
	if err == nil {
		t.Error("expected error for short frame")
	}
}

func TestXorChecksum(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	cs := xorChecksum(data)
	if cs != 0x00 {
		t.Errorf("expected 0x00, got 0x%02X", cs)
	}
}
