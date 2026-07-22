package most

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "most" {
		t.Errorf("expected most, got %s", p.Name())
	}
	if len(p.Variants()) == 0 {
		t.Error("should have at least one variant")
	}
	if p.DefaultPort() < 0 {
		t.Error("default port must not be negative")
	}
}

func TestProtocol_NewCodec(t *testing.T) {
	p := New()
	for _, v := range []string{"optical", "serial"} {
		c, err := p.NewCodec(v)
		if err != nil {
			t.Errorf("NewCodec(%q) returned error: %v", v, err)
		}
		if c == nil {
			t.Errorf("NewCodec(%q) returned nil codec", v)
		}
	}
}

func TestProtocol_RejectsUnknownVariant(t *testing.T) {
	p := New()
	_, err := p.NewCodec("invalid")
	if err == nil {
		t.Error("expected error for unknown variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "read",
		Address:  "0x0100",
		Count:    4,
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	expected := "AT+READ=0x0100,4\r\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeReadDefaults(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "AT+READ=0x0000,1\r\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeWrite(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "write",
		Address:  "0x0200",
		Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	expected := "AT+WRITE=0x0200,DEADBEEF\r\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeStatus(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "status"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "AT+STATUS\r\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Encode(&kernel.Request{Function: "bogus"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestDecodeOK(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("+OK:48 65 6C 6C 6F"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("Hello")
	if string(resp.Data) != string(expected) {
		t.Errorf("expected %q, got %q", expected, resp.Data)
	}
}

func TestDecodeOKCompact(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("+OK:48656C6C6F"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("Hello")
	if string(resp.Data) != string(expected) {
		t.Errorf("expected %q, got %q", expected, resp.Data)
	}
}

func TestDecodeOKOddLength(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("+OK:A"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data[0] != 0x0A {
		t.Errorf("expected 0x0A, got 0x%02X", resp.Data[0])
	}
}

func TestDecodeError(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Decode([]byte("+ERR:02"))
	if err == nil {
		t.Error("expected error for +ERR response")
	}
}

func TestDecodeRaw(t *testing.T) {
	c, err := New().NewCodec("serial")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("SOME OTHER RESPONSE"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["raw"] != "SOME OTHER RESPONSE" {
		t.Errorf("expected raw metadata set")
	}
}

func TestDriverNewSerialDriver(t *testing.T) {
	b, c, err := NewSerialDriver("/dev/ttyUSB0")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}

func TestParseHex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"48 65 6C 6C 6F", "Hello"},
		{"48656C6C6F", "Hello"},
		{"0x48 0x65 0x6C 0x6C 0x6F", "Hello"},
		{"", ""},
		{"41", "A"},
		{"0x41", "A"},
	}
	for _, tt := range tests {
		result, err := parseHex(tt.input)
		if err != nil {
			t.Errorf("parseHex(%q) error: %v", tt.input, err)
			continue
		}
		if string(result) != tt.expected {
			t.Errorf("parseHex(%q): expected %q, got %q", tt.input, tt.expected, string(result))
		}
	}
}
