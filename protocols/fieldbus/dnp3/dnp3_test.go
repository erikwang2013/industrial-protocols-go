package dnp3

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "dnp3" {
		t.Errorf("expected dnp3, got %s", p.Name())
	}
	if p.DefaultPort() != 20000 {
		t.Errorf("expected 20000, got %d", p.DefaultPort())
	}
}

func TestEncodeRead(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	data, err := c.Encode(&kernel.Request{Function: "read"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 12 {
		t.Fatalf("expected at least 12 bytes, got %d", len(data))
	}
	if data[0] != TPDUStart1 || data[1] != TPDUStart2 {
		t.Error("missing TPDU start bytes")
	}
}

func TestEncodeClass0Poll(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	data, err := c.Encode(&kernel.Request{Function: "class0_poll"})
	if err != nil {
		t.Fatal(err)
	}
	control := data[3]
	if control&FIR == 0 || control&FIN == 0 {
		t.Error("class 0 poll should set FIR and FIN")
	}
}

func TestDecodeResponse(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// Build a minimal TPDU-wrapped response
	apdu := []byte{0xC1, FuncResponse}
	tpdu := make([]byte, 10+len(apdu)+2)
	tpdu[0], tpdu[1] = TPDUStart1, TPDUStart2
	tpdu[2] = byte(len(apdu) + 5)
	tpdu[3] = 0xC1
	tpdu[4], tpdu[5] = 1, 0
	tpdu[6], tpdu[7] = 2, 0
	hdrCRC := crc16DNP(tpdu[3:8])
	tpdu[8], tpdu[9] = byte(hdrCRC), byte(hdrCRC>>8)
	copy(tpdu[10:], apdu)
	apduCRC := crc16DNP(apdu)
	tpdu[10+len(apdu)], tpdu[10+len(apdu)+1] = byte(apduCRC), byte(apduCRC>>8)

	resp, err := c.Decode(tpdu)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["function"].(int) != FuncResponse {
		t.Error("expected response function")
	}
}

func TestCRC16DNP(t *testing.T) {
	// Known vector: DNP3 CRC of empty data = 0x0000 before final XOR;
	// final XOR (^crc) inverts to 0xFFFF.
	crc := crc16DNP([]byte{})
	if crc != 0xFFFF {
		t.Errorf("expected 0xFFFF for empty data, got 0x%04X", crc)
	}
}
