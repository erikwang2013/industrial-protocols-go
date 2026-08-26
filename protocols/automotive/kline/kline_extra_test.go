// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kline

import (
	"errors"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestFastInitSequence(t *testing.T) {
	seq := FastInitSequence()
	if len(seq) != 1 || seq[0] != SIDFastInit {
		t.Errorf("FastInitSequence() = %v, want [0x33]", seq)
	}
}

func TestKeyBytes(t *testing.T) {
	kb := KeyBytes()
	if len(kb) != 2 || kb[0] != KeyByte1 || kb[1] != KeyByte2 {
		t.Errorf("KeyBytes() = %v, want [0x08 0x08]", kb)
	}
}

func TestVerifyKeyResponse(t *testing.T) {
	// ECU responds with the bitwise complement of each key byte.
	if !VerifyKeyResponse([]byte{0xFF ^ KeyByte1, 0xFF ^ KeyByte2}) {
		t.Error("expected valid key response to verify")
	}
	if VerifyKeyResponse([]byte{0x00, 0x00}) {
		t.Error("expected wrong key response to fail")
	}
	if VerifyKeyResponse([]byte{0xF7}) {
		t.Error("expected short response to fail")
	}
	if VerifyKeyResponse(nil) {
		t.Error("expected empty response to fail")
	}
}

func TestChecksum(t *testing.T) {
	if got := checksum([]byte{0x68, 0x10, 0xF1, 0x01}); got != 0x6A { // 0x68+0x10+0xF1+0x01 = 362, 低 8 位 0x6A
		t.Errorf("checksum = 0x%02X, want 0x6A", got)
	}
}

func TestEncodeDefaultTargetSource(t *testing.T) {
	c, _ := New().NewCodec("serial")
	data, err := c.Encode(&kernel.Request{Function: "current_data", Data: []byte{0x0C}})
	if err != nil {
		t.Fatal(err)
	}
	// Format + default target 0x33 + default source 0xF1 + SID + data + CS
	if len(data) != 4+1+1 {
		t.Fatalf("len = %d, want 6", len(data))
	}
	if data[1] != 0x33 || data[2] != 0xF1 {
		t.Errorf("target/source = %02X/%02X, want 33/F1", data[1], data[2])
	}
	if data[len(data)-1] != checksum(data[:len(data)-1]) {
		t.Error("trailing byte is not the checksum")
	}
}

func TestEncodeUnknownFunctionDefaultsToCurrentData(t *testing.T) {
	c, _ := New().NewCodec("serial")
	data, err := c.Encode(&kernel.Request{Function: "bogus"})
	if err != nil {
		t.Fatal(err)
	}
	if data[3] != SIDCurrentData {
		t.Errorf("SID = 0x%02X, want 0x01", data[3])
	}
}

func TestEncodeRequestDTCs(t *testing.T) {
	c, _ := New().NewCodec("serial")
	data, err := c.Encode(&kernel.Request{Function: "request_dtcs"})
	if err != nil {
		t.Fatal(err)
	}
	if data[3] != SIDRequestDTCs {
		t.Errorf("SID = 0x%02X, want 0x03", data[3])
	}
}

func TestDecodeFastInitSync(t *testing.T) {
	c, _ := New().NewCodec("serial")
	resp, err := c.Decode([]byte{FastInitSync})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["fast_init_sync"] != true {
		t.Error("expected fast_init_sync metadata")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("serial")
	if _, err := c.Decode([]byte{0x68, 0x10}); err == nil {
		t.Fatal("expected error for 2-byte frame")
	}
}

func TestDecodeResponseMetadata(t *testing.T) {
	c, _ := New().NewCodec("serial")
	// Response SID 0x41 (0x01 + 0x40), source 0xF1, target 0x10.
	payload := []byte{RespFormat, 0xF1, 0x10, 0x41, 0x00, 0x0C}
	payload = append(payload, checksum(payload))
	resp, err := c.Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["sid"] != byte(0x01) {
		t.Errorf("sid = %v, want 0x01", resp.Metadata["sid"])
	}
	if resp.Metadata["source"] != byte(0xF1) {
		t.Errorf("source = %v, want 0xF1", resp.Metadata["source"])
	}
	if resp.Metadata["target"] != byte(0x10) {
		t.Errorf("target = %v, want 0x10", resp.Metadata["target"])
	}
	if len(resp.Data) != 2 || resp.Data[0] != 0x00 || resp.Data[1] != 0x0C {
		t.Errorf("data = %v, want [0x00 0x0C]", resp.Data)
	}
}

func TestNewCodecRejectsUnknownVariant(t *testing.T) {
	if _, err := New().NewCodec("tcp"); err == nil {
		t.Fatal("expected error for unsupported variant")
	}
}

func TestProtocolImplementsKernelInterface(t *testing.T) {
	var _ kernel.Protocol = New()
	// Compile-time check that the codec satisfies kernel.Codec.
	var _ kernel.Codec = &kLineCodec{}
}

func TestErrInvalidAddressNotUsed(t *testing.T) {
	// kLine rejects invalid variants with a descriptive error, not a sentinel.
	_, err := New().NewCodec("nope")
	if errors.Is(err, kernel.ErrInvalidAddress) {
		t.Error("unexpected sentinel error")
	}
}
