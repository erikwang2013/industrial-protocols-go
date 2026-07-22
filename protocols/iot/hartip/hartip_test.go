package hartip

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "hartip" {
		t.Errorf("expected hartip, got %s", p.Name())
	}
	if p.DefaultPort() != 5094 {
		t.Errorf("expected 5094, got %d", p.DefaultPort())
	}
	if len(p.Variants()) != 2 {
		t.Errorf("expected 2 variants, got %d", len(p.Variants()))
	}

	// Valid variants
	for _, v := range []string{"tcp", "udp"} {
		c, err := p.NewCodec(v)
		if err != nil {
			t.Errorf("variant %q should be valid: %v", v, err)
		}
		if c == nil {
			t.Errorf("expected non-nil codec for variant %q", v)
		}
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "read_unique_id"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 8 {
		t.Errorf("too short: %d bytes", len(data))
	}
	if data[0] != 1 {
		t.Error("expected version 1")
	}

	resp, err := c.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if cmd, ok := resp.Metadata["command"].(int); !ok || cmd != 0 {
		t.Errorf("expected command 0, got %v", resp.Metadata["command"])
	}
}

func TestEncodeDecodeCommand3(t *testing.T) {
	c, _ := New().NewCodec("udp")
	req := &kernel.Request{Function: "read_dynamic_variables"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if cmd, ok := resp.Metadata["command"].(int); !ok || cmd != 3 {
		t.Errorf("expected command 3, got %v", resp.Metadata["command"])
	}
}

func TestFrameTooShort(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	_, err := c.Decode([]byte{0x01, 0x00})
	if err == nil {
		t.Error("expected error for short frame")
	}
}

func TestLengthMismatch(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// Header says length 100 but only 10 bytes provided
	data := make([]byte, 8+10)
	data[0] = 1
	data[6] = 0
	data[7] = 100
	_, err := c.Decode(data)
	if err == nil {
		t.Error("expected error for length mismatch")
	}
}
