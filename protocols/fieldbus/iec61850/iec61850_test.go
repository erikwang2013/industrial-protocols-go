// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package iec61850

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "iec61850" {
		t.Errorf("expected iec61850, got %s", p.Name())
	}
	if p.DefaultPort() != 102 {
		t.Errorf("expected 102, got %d", p.DefaultPort())
	}
}

func TestEncodeInitiate(t *testing.T) {
	c, _ := New().NewCodec("mms")
	data, err := c.Encode(&kernel.Request{Function: "initiate"})
	if err != nil {
		t.Fatal(err)
	}
	if data[0] != 0xA8 {
		t.Errorf("expected initiate tag 0xA8, got 0x%02X", data[0])
	}
}

func TestEncodeConclude(t *testing.T) {
	c, _ := New().NewCodec("mms")
	data, err := c.Encode(&kernel.Request{Function: "conclude"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 4 {
		t.Fatalf("expected 4 bytes, got %d", len(data))
	}
}

func TestEncodeRead(t *testing.T) {
	c, _ := New().NewCodec("mms")
	data, err := c.Encode(&kernel.Request{Function: "read", Address: "TestVar"})
	if err != nil {
		t.Fatal(err)
	}
	if data[0] != 0xA4 {
		t.Errorf("expected read tag 0xA4, got 0x%02X", data[0])
	}
}

func TestEncodeWrite(t *testing.T) {
	c, _ := New().NewCodec("mms")
	data, err := c.Encode(&kernel.Request{Function: "write", Address: "Var", Data: []byte{42}})
	if err != nil {
		t.Fatal(err)
	}
	if data[0] != 0xA5 {
		t.Errorf("expected write tag 0xA5, got 0x%02X", data[0])
	}
}

func TestDecodeReadResponse(t *testing.T) {
	c, _ := New().NewCodec("mms")
	// Read response with integer value 42: 0xA4, len, 0x85, 1, 42
	resp, err := c.Decode([]byte{0xA4, 0x03, 0x85, 0x01, 42})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "read_response" {
		t.Errorf("expected read_response, got %s", resp.Function)
	}
}

func TestDecodeConcludeResponse(t *testing.T) {
	c, _ := New().NewCodec("mms")
	resp, err := c.Decode([]byte{0xA3, 0x02, 0x80, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "conclude_response" {
		t.Errorf("expected conclude_response, got %s", resp.Function)
	}
}
