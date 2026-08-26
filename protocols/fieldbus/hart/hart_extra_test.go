// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package hart

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraEncodeFrameLayout(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	raw, err := c.Encode(&kernel.Request{Function: "read_unique_id"})
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != 10 {
		t.Fatalf("len = %d, want 10", len(raw))
	}
	for i := 0; i < 5; i++ {
		if raw[i] != 0xFF {
			t.Errorf("preamble[%d] = 0x%02X, want 0xFF", i, raw[i])
		}
	}
	if raw[5] != 0x02 || raw[6] != 0x80 || raw[7] != 0x00 || raw[8] != 0x00 {
		t.Errorf("header = % X, want 02 80 00 00", raw[5:9])
	}
	// Checksum covers delimiter..byte-count (5 bytes).
	if want := xorChecksum(raw[5:9]); raw[9] != want {
		t.Errorf("checksum = 0x%02X, want 0x%02X", raw[9], want)
	}
}

func TestExtraEncodeCommand3(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	raw, err := c.Encode(&kernel.Request{Function: "read_dynamic_variables"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[7] != 3 {
		t.Errorf("cmd = %d, want 3", raw[7])
	}
}

func TestExtraEncodeUnknownFunctionDefaultsCmd0(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	raw, err := c.Encode(&kernel.Request{Function: "bogus"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[7] != 0 {
		t.Errorf("cmd = %d, want 0", raw[7])
	}
}

func TestExtraEncodePollingAddr(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	raw, _ := c.Encode(&kernel.Request{
		Function: "read_unique_id",
		Metadata: map[string]any{"polling_addr": byte(0x21)},
	})
	if raw[6] != 0xA1 {
		t.Errorf("addr = 0x%02X, want 0xA1 (polling 0x21 | master 0x80)", raw[6])
	}
}

func TestExtraDecodeCommand3PVUnits(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	// Short frame reply: delim 0x06, addr 0x80, cmd 3, byte count 1, PV units 0x05.
	frame := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x06, 0x80, 0x03, 0x01, 0x05, 0x81}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["pv_units"] != byte(0x05) {
		t.Errorf("pv_units = %v, want 5", resp.Metadata["pv_units"])
	}
	if resp.Metadata["command"] != 3 {
		t.Errorf("command = %v, want 3", resp.Metadata["command"])
	}
}

func TestExtraDecodeLongFrameCommand3(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	// Long frame: delim 0x86, 5-byte address, cmd 3, byte count 1, PV units 0x09.
	payload := []byte{0x86, 0x01, 0x02, 0x03, 0x04, 0x05, 0x03, 0x01, 0x09}
	payload = append(payload, xorChecksum(payload))
	frame := append([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, payload...)
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["long_frame"] != true {
		t.Error("long_frame = false, want true")
	}
	if resp.Metadata["pv_units"] != byte(0x09) {
		t.Errorf("pv_units = %v, want 9", resp.Metadata["pv_units"])
	}
}

func TestExtraDecodeByteCountMismatch(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	// Byte count 3 but only 1 data byte follows the header.
	frame := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x06, 0x80, 0x03, 0x03, 0x01}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for byte count mismatch")
	}
}

func TestExtraDecodeChecksumMismatch(t *testing.T) {
	c, _ := New().NewCodec("fsk")
	frame := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x06, 0x80, 0x03, 0x01, 0x05, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for checksum mismatch")
	}
}
