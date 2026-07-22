package canopen

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "canopen" {
		t.Errorf("expected canopen, got %s", p.Name())
	}
	if len(p.Variants()) == 0 || p.Variants()[0] != "can" {
		t.Errorf("expected first variant 'can', got %v", p.Variants())
	}
	if p.DefaultPort() != 0 {
		t.Errorf("expected default port 0, got %d", p.DefaultPort())
	}
	c, err := p.NewCodec("can")
	if err != nil {
		t.Fatalf("NewCodec(can) returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewCodec(can) returned nil codec")
	}
}

func TestProtocolRejectsUnknownVariant(t *testing.T) {
	p := New()
	_, err := p.NewCodec("serial")
	if err == nil {
		t.Error("expected error for unknown variant, got nil")
	}
}

func TestEncodeSDORead(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "sdo_read",
		Metadata: map[string]any{
			"index": float64(0x1000),
			"sub":   float64(0),
		},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// First byte of SDO read is 0x40 (initiate domain upload)
	if data[4] != 0x40 {
		t.Errorf("expected SDO read command 0x40, got 0x%02X", data[4])
	}
}

func TestEncodeSDOWrite(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	req := &kernel.Request{
		Function: "sdo_write",
		Data:     []byte{0x11, 0x22, 0x33, 0x44},
		Metadata: map[string]any{
			"index": float64(0x1000),
			"sub":   float64(0),
		},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// First byte of SDO write expedited (4 bytes) is 0x23
	if data[4] != 0x23 {
		t.Errorf("expected SDO write command 0x23, got 0x%02X", data[4])
	}
}

func TestEncodeNMTStart(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "nmt_start"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	// NMT start: first data byte is 0x01
	if data[4] != 0x01 {
		t.Errorf("expected NMT start 0x01, got 0x%02X", data[4])
	}
}

func TestEncodeNMTStop(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "nmt_stop"})
	if err != nil {
		t.Fatal(err)
	}
	if data[4] != 0x02 {
		t.Errorf("expected NMT stop 0x02, got 0x%02X", data[4])
	}
}

func TestEncodeNMTReset(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "nmt_reset"})
	if err != nil {
		t.Fatal(err)
	}
	if data[4] != 0x82 {
		t.Errorf("expected NMT reset 0x82, got 0x%02X", data[4])
	}
}

func TestEncodeHeartbeat(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Encode(&kernel.Request{Function: "heartbeat"})
	if err != nil {
		t.Fatal(err)
	}
	if data[4] != 0x05 {
		t.Errorf("expected heartbeat state 0x05, got 0x%02X", data[4])
	}
}

func TestEncodeUnknownFunction(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Encode(&kernel.Request{Function: "bogus"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestDecodeBootup(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	f := marshalCAN(bridge.CANFrame{ID: 0x700 + 1, Data: []byte{0x00}})
	resp, err := c.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "bootup" {
		t.Errorf("expected bootup, got %v", resp.Metadata["type"])
	}
}

func TestDecodeSDOAbort(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	f := marshalCAN(bridge.CANFrame{
		ID:   0x580 + 1,
		Data: []byte{0x80, 0x10, 0x10, 0x01, 0x00, 0x00, 0x06, 0x06},
	})
	_, err = c.Decode(f)
	if err == nil {
		t.Error("expected ProtocolError for SDO abort")
	}
	if pe, ok := err.(*kernel.ProtocolError); ok {
		if pe.Message != "SDO abort" {
			t.Errorf("expected 'SDO abort', got %q", pe.Message)
		}
	} else {
		t.Errorf("expected *kernel.ProtocolError, got %T", err)
	}
}

func TestDecodeTooShort(t *testing.T) {
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Decode([]byte{0x00, 0x01})
	if err == nil {
		t.Error("expected error for short frame")
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	orig := bridge.CANFrame{
		ID:   0x701,
		Data: []byte{0xAB, 0xCD, 0xEF, 0x01},
		Ext:  false,
	}
	buf := marshalCAN(orig)
	restored := unmarshalCAN(buf)
	if restored.ID != orig.ID {
		t.Errorf("ID mismatch: %x vs %x", restored.ID, orig.ID)
	}
	if restored.Ext != orig.Ext {
		t.Error("Ext flag mismatch")
	}
	for i := 0; i < 4; i++ {
		if restored.Data[i] != orig.Data[i] {
			t.Errorf("data byte %d mismatch: %x vs %x", i, restored.Data[i], orig.Data[i])
		}
	}
}

func TestMarshalRoundTripExtended(t *testing.T) {
	orig := bridge.CANFrame{
		ID:   0x1FFFFFFF,
		Data: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		Ext:  true,
	}
	buf := marshalCAN(orig)
	restored := unmarshalCAN(buf)
	if restored.ID != orig.ID&0x7FFFFFFF {
		t.Errorf("ID mismatch: %x vs %x", restored.ID, orig.ID&0x7FFFFFFF)
	}
	if !restored.Ext {
		t.Error("expected extended frame flag")
	}
}
