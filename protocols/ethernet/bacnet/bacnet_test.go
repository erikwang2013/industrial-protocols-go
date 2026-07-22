package bacnet

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "bacnet" {
		t.Errorf("expected bacnet, got %s", p.Name())
	}
	if p.DefaultPort() != 47808 {
		t.Errorf("expected 47808, got %d", p.DefaultPort())
	}
}

func TestEncodeWhoIs(t *testing.T) {
	c, _ := New().NewCodec("ip")
	data, err := c.Encode(&kernel.Request{Function: "who_is"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(data))
	}
	if data[0] != BVLCTypeOrigBroadcast {
		t.Errorf("expected broadcast 0x81, got 0x%02X", data[0])
	}
	if data[7] != ServiceWhoIs {
		t.Errorf("expected Who-Is (3), got %d", data[7])
	}
}

func TestEncodeReadProperty(t *testing.T) {
	c, _ := New().NewCodec("ip")
	req := &kernel.Request{
		Function: "read_property",
		Metadata: map[string]any{
			"object_type": byte(0), "object_instance": float64(1), "property_id": float64(85),
		},
	}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 18 {
		t.Fatalf("expected 18 bytes, got %d", len(data))
	}
	if data[0] != BVLCTypeOrigUnicast {
		t.Error("expected unicast")
	}
	if data[9] != ServiceReadProperty {
		t.Errorf("expected ReadProperty (12), got %d", data[9])
	}
}

func TestDecodeIAm(t *testing.T) {
	// Simulated I-Am: BVLC(4) + NPDU(2) + APDU(0x10, 0) + ObjectID(4) + MaxAPDU(2) + Segmentation(1) + VendorID(2)
	iAm := []byte{0x81, 0x0B, 0x00, 0x0F, 0x01, 0x00, 0x10, 0x00,
		0xC0, 0x00, 0x00, 0x01, 0x01, 0xE0, 0x00, 0x00, 0x00, 0x00, 0x00}
	c, _ := New().NewCodec("ip")
	resp, err := c.Decode(iAm)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["bvlc_type"].(byte) != BVLCTypeOrigUnicast {
		t.Error("expected unicast")
	}
	if resp.Metadata["service"].(int) != ServiceIAm {
		t.Error("expected I-Am service")
	}
}
