// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package bacnet

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestEncodeWhoIsExact(t *testing.T) {
	c, _ := New().NewCodec("ip")
	raw, err := c.Encode(&kernel.Request{Function: "who_is"})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x82, 0x0B, 0x00, 0x04, 0x01, 0x00, 0x10, 0x03}
	if !bytes.Equal(raw, want) {
		t.Errorf("who_is = % X, want % X", raw, want)
	}
}

func TestEncodeUnknownFunctionDefaultsWhoIs(t *testing.T) {
	c, _ := New().NewCodec("ip")
	raw, err := c.Encode(&kernel.Request{Function: "bogus"})
	if err != nil {
		t.Fatal(err)
	}
	if raw[0] != BVLCTypeOrigBroadcast {
		t.Errorf("bvlc type = 0x%02X, want broadcast", raw[0])
	}
}

func TestEncodeReadPropertyExact(t *testing.T) {
	c, _ := New().NewCodec("ip")
	raw, err := c.Encode(&kernel.Request{
		Function: "read_property",
		Metadata: map[string]any{
			"object_type":      byte(0x02),
			"object_instance":  float64(100),
			"property_id":      float64(85),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != 18 {
		t.Fatalf("len = %d, want 18", len(raw))
	}
	if raw[0] != BVLCTypeOrigUnicast {
		t.Errorf("bvlc type = 0x%02X, want unicast", raw[0])
	}
	if binary.BigEndian.Uint16(raw[2:4]) != 14 {
		t.Errorf("bvlc length = %d, want 14", binary.BigEndian.Uint16(raw[2:4]))
	}
	if raw[6] != APDUTypeConfirmedReq || raw[9] != ServiceReadProperty {
		t.Errorf("apdu = %02X/%02X, want 0x00/0x0C", raw[6], raw[9])
	}
	// Object ID: 10-bit type << 22 | 22-bit instance.
	objID := binary.BigEndian.Uint32(raw[10:14])
	wantObjID := (uint32(0x02) & 0x3FF) << 22 | 100
	if objID != wantObjID {
		t.Errorf("object id = 0x%08X, want 0x%08X", objID, wantObjID)
	}
	if binary.BigEndian.Uint32(raw[14:18]) != 85 {
		t.Errorf("property id = %d, want 85", binary.BigEndian.Uint32(raw[14:18]))
	}
}

func TestEncodeReadPropertyDefaults(t *testing.T) {
	c, _ := New().NewCodec("ip")
	raw, err := c.Encode(&kernel.Request{Function: "read_property"})
	if err != nil {
		t.Fatal(err)
	}
	// Defaults: object type 0, instance 0, property 8 (present value).
	if binary.BigEndian.Uint32(raw[10:14]) != 0 {
		t.Error("default object id should be 0")
	}
	if binary.BigEndian.Uint32(raw[14:18]) != 8 {
		t.Error("default property id should be 8")
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("ip")
	if _, err := c.Decode([]byte{0x81, 0x0B, 0x00}); err == nil {
		t.Fatal("expected error for short frame")
	}
}

func TestDecodeUnsupportedBVLC(t *testing.T) {
	c, _ := New().NewCodec("ip")
	if _, err := c.Decode([]byte{0x80, 0x0B, 0x00, 0x04, 0x01, 0x00, 0x10, 0x03}); err == nil {
		t.Fatal("expected error for unsupported BVLC type")
	}
}

func TestDecodeMetadata(t *testing.T) {
	c, _ := New().NewCodec("ip")
	// Who-Is request frame decoded back.
	raw, _ := c.Encode(&kernel.Request{Function: "who_is"})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["bvlc_type"] != byte(0x82) {
		t.Errorf("bvlc_type = %v", resp.Metadata["bvlc_type"])
	}
	if resp.Metadata["apdu_type"] != 1 {
		t.Errorf("apdu_type = %v, want 1 (unconfirmed)", resp.Metadata["apdu_type"])
	}
	if resp.Metadata["service"] != ServiceWhoIs {
		t.Errorf("service = %v, want 3", resp.Metadata["service"])
	}
}

func TestDecodeShortNPDU(t *testing.T) {
	c, _ := New().NewCodec("ip")
	// BVLC says 4 bytes of NPDU+APDU but only 1 byte follows.
	if _, err := c.Decode([]byte{0x82, 0x0B, 0x00, 0x04, 0x01}); err == nil {
		t.Fatal("expected error for short NPDU")
	}
}
