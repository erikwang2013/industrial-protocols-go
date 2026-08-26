// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package lin

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestCalcPIDVectors(t *testing.T) {
	// Protected ID parity bits: P0 at bit 6 (even parity of b0,b1,b2,b4),
	// P1 at bit 7 (inverse parity of b1,b3,b4,b5).
	tests := []struct {
		id   byte
		want byte
	}{
		{0x00, 0x80},
		{0x0C, 0x4C},
		{0x11, 0x11},
		{0x3C, 0x3C},
		{0x3D, 0x7D},
		{0x1A, 0x1A},
		{0x2A, 0x6A},
	}
	for _, tt := range tests {
		if got := calcPID(tt.id); got != tt.want {
			t.Errorf("calcPID(0x%02X) = 0x%02X, want 0x%02X", tt.id, got, tt.want)
		}
	}
}

func TestCalcPIDMasksID(t *testing.T) {
	// Bits above 5 are ignored.
	if got := calcPID(0xFF); got != calcPID(0x3F) {
		t.Errorf("calcPID(0xFF) = 0x%02X, want 0x%02X", got, calcPID(0x3F))
	}
}

func TestCalcClassicChecksum(t *testing.T) {
	// Classic checksum = 0xFF - sum(data) (inverted sum).
	if got := calcClassicChecksum([]byte{0x01, 0x02, 0x03}); got != 0xF9 {
		t.Errorf("classic checksum = 0x%02X, want 0xF9", got)
	}
	if got := calcClassicChecksum(nil); got != 0xFF {
		t.Errorf("classic checksum of empty = 0x%02X, want 0xFF", got)
	}
}

func TestCalcEnhancedChecksum(t *testing.T) {
	// Enhanced checksum includes the PID in the sum (with carry fold).
	if got := calcEnhancedChecksum(0x3C, []byte{0x01, 0x02, 0x03}); got != 0xBD {
		t.Errorf("enhanced checksum = 0x%02X, want 0xBD", got)
	}
	// Including the PID must differ from the classic checksum.
	if calcEnhancedChecksum(0x3C, []byte{0x01}) == calcClassicChecksum([]byte{0x01}) {
		t.Error("enhanced checksum should differ from classic when PID nonzero")
	}
}

func TestEncodeClassicFrame(t *testing.T) {
	c, _ := New().NewCodec("uart")
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0x11, 0x22},
		Metadata: map[string]any{"id": byte(0x0C)},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Sync (0x55) + PID(0x4C) + data + classic checksum.
	if data[0] != 0x55 {
		t.Errorf("sync byte = 0x%02X, want 0x55", data[0])
	}
	if data[1] != 0x4C {
		t.Errorf("PID = 0x%02X, want 0x4C", data[1])
	}
	want := calcClassicChecksum([]byte{0x11, 0x22})
	if data[len(data)-1] != want {
		t.Errorf("checksum = 0x%02X, want 0x%02X", data[len(data)-1], want)
	}
}

func TestEncodeEnhancedFrame(t *testing.T) {
	c, _ := New().NewCodec("uart")
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Data:     []byte{0x11},
		Metadata: map[string]any{"id": byte(0x0C), "enhanced_checksum": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := calcEnhancedChecksum(0x4C, []byte{0x11})
	if data[len(data)-1] != want {
		t.Errorf("checksum = 0x%02X, want 0x%02X", data[len(data)-1], want)
	}
}

func TestEncodeIDFromFloat64(t *testing.T) {
	c, _ := New().NewCodec("uart")
	data, err := c.Encode(&kernel.Request{
		Function: "write",
		Metadata: map[string]any{"id": float64(0x11)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if data[1] != calcPID(0x11) {
		t.Errorf("PID = 0x%02X, want 0x%02X", data[1], calcPID(0x11))
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	c, _ := New().NewCodec("uart")
	req := &kernel.Request{
		Function: "write",
		Data:     []byte{0xDE, 0xAD},
		Metadata: map[string]any{"id": byte(0x22)},
	}
	frame, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["id"] != int(0x22) {
		t.Errorf("id = %v, want 34", resp.Metadata["id"])
	}
	if len(resp.Data) != 2 || resp.Data[0] != 0xDE {
		t.Errorf("data = %v, want [0xDE 0xAD]", resp.Data)
	}
}

func TestDecodeSyncError(t *testing.T) {
	c, _ := New().NewCodec("uart")
	resp, err := c.Decode([]byte{0x00, 0x4C, 0x01, 0xFE})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["sync_error"] != true {
		t.Error("expected sync_error metadata for missing sync byte")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("uart")
	if _, err := c.Decode([]byte{0x55, 0x4C}); err == nil {
		t.Fatal("expected error for 2-byte frame")
	}
}
