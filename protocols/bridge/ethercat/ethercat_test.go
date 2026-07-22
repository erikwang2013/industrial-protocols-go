// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package ethercat

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "ethercat" {
		t.Errorf("expected ethercat, got %s", p.Name())
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
	for _, v := range []string{"ethercat", "cmd"} {
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

func TestEncodeUpload(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "upload",
		Address:  "0x1000",
		Count:    4,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "upload 0x1000 4\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeUploadDefaults(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "upload"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "upload 0x0000 1\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeDownload(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{
		Function: "download",
		Address:  "0x2000",
		Data:     []byte{0xCA, 0xFE},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := "download 0x2000 0xCAFE\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeSlaves(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "slaves"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "slaves\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Encode(&kernel.Request{Function: "bogus"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestDecodeHexOutput(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("0x1000: CA FE BA BE"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["raw"] == nil {
		t.Error("expected raw metadata")
	}
}

func TestDecodePlainHex(t *testing.T) {
	c, err := New().NewCodec("cmd")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Decode([]byte("CAFE"))
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) < 2 {
		t.Error("expected parsed hex data")
	}
	expected := []byte{0xCA, 0xFE}
	if string(resp.Data) != string(expected) {
		t.Errorf("expected %X, got %X", expected, resp.Data)
	}
}

func TestDriverNewCmdDriver(t *testing.T) {
	b, c, err := NewCmdDriver("/usr/bin/ethercat")
	if err != nil {
		t.Fatal(err)
	}
	if b == nil || c == nil {
		t.Error("expected non-nil bridge and codec")
	}
}
