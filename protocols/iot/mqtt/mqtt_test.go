// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package mqtt

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "mqtt" {
		t.Errorf("expected mqtt, got %s", p.Name())
	}
	if p.DefaultPort() != 1883 {
		t.Errorf("expected 1883, got %d", p.DefaultPort())
	}
	c, _ := p.NewCodec("tcp")
	if c == nil {
		t.Error("NewCodec returned nil")
	}
}

func TestEncodeConnect(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "connect", Metadata: map[string]any{"client_id": "test"}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 10 {
		t.Errorf("connect frame too short: %d bytes", len(data))
	}
	if data[0] != 0x10 {
		t.Errorf("expected CONNECT type 0x10, got 0x%02X", data[0])
	}
}

func TestEncodePublish(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "publish", Address: "sensor/temp", Data: []byte("25.5")}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 5 {
		t.Errorf("publish frame too short: %d bytes", len(data))
	}
}

func TestEncodeSubscribe(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "subscribe", Address: "sensor/#"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 6 {
		t.Errorf("subscribe frame too short: %d bytes", len(data))
	}
}

func TestEncodePingreq(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "pingreq"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 2 || data[0] != 0xC0 || data[1] != 0x00 {
		t.Errorf("expected PINGREQ [0xC0, 0x00], got %v", data)
	}
}

func TestDecodeConnack(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0x20, 0x02, 0x00, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "connack" {
		t.Errorf("expected connack, got %s", resp.Function)
	}
}

func TestDecodePublish(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0x30, 0x06, 0x00, 0x04, 't', 'e', 's', 't', 'h', 'i'})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "publish" {
		t.Errorf("expected publish, got %s", resp.Function)
	}
	if resp.Address != "test" {
		t.Errorf("expected topic 'test', got '%s'", resp.Address)
	}
}

func TestDecodeSuback(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0x90, 0x03, 0x00, 0x01, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "suback" {
		t.Errorf("expected suback, got %s", resp.Function)
	}
}

func TestDecodePingresp(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0xD0, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "pingresp" {
		t.Errorf("expected pingresp, got %s", resp.Function)
	}
}
