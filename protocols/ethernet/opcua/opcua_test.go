package opcua

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "opcua" {
		t.Errorf("expected opcua, got %s", p.Name())
	}
	if p.DefaultPort() != 4840 {
		t.Errorf("expected 4840, got %d", p.DefaultPort())
	}
}

func TestEncodeHello(t *testing.T) {
	c, _ := New().NewCodec("binary")
	req := &kernel.Request{Function: "hel", Metadata: map[string]any{"endpoint": "opc.tcp://192.168.1.1:4840"}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 32 {
		t.Errorf("HELLO frame too short: %d bytes", len(data))
	}
	if string(data[:3]) != "HEL" {
		t.Errorf("expected HEL, got %s", string(data[:3]))
	}
}

func TestEncodeOpenSecureChannel(t *testing.T) {
	c, _ := New().NewCodec("binary")
	req := &kernel.Request{Function: "open_secure_channel"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 8 {
		t.Errorf("OPN frame too short: %d bytes", len(data))
	}
	if string(data[:3]) != "OPN" {
		t.Errorf("expected OPN, got %s", string(data[:3]))
	}
}

func TestDecodeAck(t *testing.T) {
	c, _ := New().NewCodec("binary")
	ack := make([]byte, 28)
	copy(ack, "ACK")
	ack[3] = 'F'
	resp, err := c.Decode(ack)
	if err != nil {
		t.Fatal(err)
	}
	fn, ok := resp.Metadata["function"].(string)
	if !ok || fn != "ack" {
		t.Errorf("expected ack, got %v", resp.Metadata["function"])
	}
}

func TestDecodeOpenResponse(t *testing.T) {
	c, _ := New().NewCodec("binary")
	opn := make([]byte, 8)
	copy(opn, "OPN")
	opn[3] = 'F'
	resp, err := c.Decode(opn)
	if err != nil {
		t.Fatal(err)
	}
	fn, ok := resp.Metadata["function"].(string)
	if !ok || fn != "open_response" {
		t.Errorf("expected open_response, got %v", resp.Metadata["function"])
	}
}

func TestUnknownFunction(t *testing.T) {
	c, _ := New().NewCodec("binary")
	_, err := c.Encode(&kernel.Request{Function: "unknown"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}
