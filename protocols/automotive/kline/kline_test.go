package kline

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "kline" {
		t.Errorf("expected kline, got %s", p.Name())
	}
	if p.DefaultPort() != DefaultBaud {
		t.Error("wrong default baud")
	}
}

func TestEncodeCurrentData(t *testing.T) {
	c, _ := New().NewCodec("serial")
	req := &kernel.Request{
		Function: "current_data",
		Data:     []byte{0x00, 0x0C}, // PID 0x0C = RPM
		Metadata: map[string]any{"target": byte(0x10), "source": byte(0xF1)},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if data[0] != ReqFormat {
		t.Errorf("expected format 0x68, got 0x%02X", data[0])
	}
	if data[3] != SIDCurrentData {
		t.Errorf("expected SID 0x01, got 0x%02X", data[3])
	}

	cs := checksum(data[:len(data)-1])
	if cs != data[len(data)-1] {
		t.Error("checksum mismatch in encoded frame")
	}
}

func TestDecodeResponse(t *testing.T) {
	c, _ := New().NewCodec("serial")
	// Simulated SID 0x41 (Current Data response) for RPM PID: 0x00 0x0C 0x1A 0xF8 -> RPM = 1722
	payload := []byte{RespFormat, 0xF1, 0x10, 0x41, 0x00, 0x0C, 0x1A, 0xF8}
	payload = append(payload, checksum(payload))
	resp, err := c.Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	sid, _ := resp.Metadata["sid"].(byte)
	if sid != 0x01 {
		t.Errorf("expected response SID 0x01, got 0x%02X", sid)
	}
}

func TestChecksumMismatch(t *testing.T) {
	c, _ := New().NewCodec("serial")
	data := []byte{0x48, 0xF1, 0x10, 0x41, 0x00, 0x00}
	_, err := c.Decode(data)
	if err == nil {
		t.Error("expected checksum mismatch error")
	}
}

func TestKeyVerification(t *testing.T) {
	if VerifyKeyResponse([]byte{byte(0xFF ^ KeyByte1), byte(0xFF ^ KeyByte2)}) != true {
		t.Error("expected valid key response")
	}
	if VerifyKeyResponse([]byte{0x00, 0x00}) != false {
		t.Error("expected invalid key response")
	}
}
