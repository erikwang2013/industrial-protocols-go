package lin

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "lin" {
		t.Errorf("expected lin, got %s", p.Name())
	}
}

func TestCalcPID(t *testing.T) {
	pid := calcPID(0x0C) // ID 12
	if pid&0x3F != 0x0C {
		t.Errorf("expected ID 12, got %d", pid&0x3F)
	}
}

func TestEncodeMasterHeader(t *testing.T) {
	c, _ := New().NewCodec("uart")
	req := &kernel.Request{Function: "read", Metadata: map[string]any{"id": byte(0x0C)}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if data[0] != 0x55 {
		t.Error("expected sync byte 0x55")
	}
	if data[1]&0x3F != 0x0C {
		t.Error("expected ID 0x0C")
	}
}

func TestEncodeWithData(t *testing.T) {
	c, _ := New().NewCodec("uart")
	req := &kernel.Request{Function: "write", Metadata: map[string]any{"id": byte(0x10)}, Data: []byte{0x01, 0x02, 0x03}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 6 {
		t.Fatalf("expected 6 bytes (sync + pid + 3 data + checksum), got %d", len(data))
	}
}

func TestDecodeValidFrame(t *testing.T) {
	c, _ := New().NewCodec("uart")
	// sync(0x55) + pid(0x0C = ID 0x0C, no parity for test) + data(0x01,0x02) + checksum
	pid := calcPID(0x0C)
	frame := []byte{0x55, pid, 0x01, 0x02}
	frame = append(frame, calcEnhancedChecksum(pid, []byte{0x01, 0x02}))
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := resp.Metadata["id"].(int)
	if id != 0x0C {
		t.Errorf("expected ID 0x0C, got 0x%02X", id)
	}
}

func TestDecodeChecksumMismatch(t *testing.T) {
	c, _ := New().NewCodec("uart")
	frame := []byte{0x55, 0x00, 0x01, 0x02, 0xFF}
	_, err := c.Decode(frame)
	if err == nil {
		t.Error("expected checksum mismatch error")
	}
}
