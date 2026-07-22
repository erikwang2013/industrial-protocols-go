// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cclink

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "cclink" {
		t.Errorf("expected cclink, got %s", p.Name())
	}
	if len(p.Variants()) != 1 {
		t.Error("expected 1 variant")
	}
}

func TestEncodeRead(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	req := &kernel.Request{Function: "read", Metadata: map[string]any{"station": float64(1)}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// 7 bytes: start + station + ctrl + len(=0) + crc(2) + end
	if len(data) != 7 {
		t.Fatalf("expected 7 bytes, got %d", len(data))
	}
	if data[0] != StartDelim || data[len(data)-1] != StartDelim {
		t.Error("missing start/end delimiter")
	}
}

func TestEncodeWrite(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	req := &kernel.Request{
		Function: "write",
		Data:     []byte{0x01, 0x02, 0x03},
		Metadata: map[string]any{"station": float64(2)},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 10 {
		t.Fatalf("expected 10 bytes, got %d", len(data))
	}
}

func TestDecodeResponse(t *testing.T) {
	c, _ := New().NewCodec("rs485")
	// Station 1, control response, data len 2, data {0xAA, 0xBB}, CRC, end delimiter
	payload := []byte{StartDelim, 1, CtrlRead, 2, 0xAA, 0xBB}
	crc := crc16XMODEM(payload[1:])
	crcEnd := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcEnd, crc)
	frame := append(payload, crcEnd...)
	frame = append(frame, StartDelim)

	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	st, _ := resp.Metadata["station"].(int)
	if st != 1 {
		t.Errorf("expected station 1, got %d", st)
	}
}

func TestCRC16XMODEM(t *testing.T) {
	// Known vector: "123456789" -> 0x31C3
	result := crc16XMODEM([]byte("123456789"))
	if result != 0x31C3 {
		t.Errorf("expected CRC 0x31C3, got 0x%04X", result)
	}
}
