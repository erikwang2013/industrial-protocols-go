package modbus

import (
	"errors"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestTCPCodec_ReadCoils(t *testing.T) {
	c, _ := NewCodec("tcp")
	req := &kernel.Request{Function: "read_coils", Address: "0", Count: 8, Metadata: map[string]any{"unit_id": byte(1)}}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// MBAP header (7) + PDU (5) = 12
	if len(enc) != 12 {
		t.Fatalf("expected 12 bytes (7 MBAP + 5 PDU), got %d", len(enc))
	}
	// Function code is at byte 7 (after MBAP header)
	if enc[7] != 1 {
		t.Errorf("expected function code 1, got %d", enc[7])
	}
}

func TestTCPCodec_WriteSingleCoil_ON(t *testing.T) {
	c, _ := NewCodec("tcp")
	req := &kernel.Request{Function: "write_single_coil", Address: "10", Data: []byte{0xFF}, Metadata: map[string]any{"unit_id": byte(1)}}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// MBAP (7) + PDU (5) = 12; PDU bytes 3-4 should be 0xFF 0x00
	if enc[10] != 0xFF || enc[11] != 0x00 {
		t.Errorf("expected coil ON value 0xFF00, got 0x%02X%02X", enc[10], enc[11])
	}
}

func TestRTUCodec_ReadCoils(t *testing.T) {
	c, _ := NewCodec("rtu")
	req := &kernel.Request{Function: "read_coils", Address: "0", Count: 8, Metadata: map[string]any{"unit_id": byte(1)}}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// address (1) + PDU (5) + CRC (2) = 8
	if len(enc) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(enc))
	}
}

func TestRTUCodec_WriteSingleRegister(t *testing.T) {
	c, _ := NewCodec("rtu")
	req := &kernel.Request{Function: "write_single_register", Address: "100", Data: []byte{0x12, 0x34}, Metadata: map[string]any{"unit_id": byte(1)}}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// address (1) + PDU (5) + CRC (2) = 8
	if len(enc) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(enc))
	}
	// PDU bytes 3-4 should be 0x12 0x34 (register value)
	if enc[4] != 0x12 || enc[5] != 0x34 {
		t.Errorf("expected register value 0x1234, got 0x%02X%02X", enc[4], enc[5])
	}
}

func TestDecodeReadHoldingRegisters(t *testing.T) {
	c, _ := NewCodec("tcp")
	// TCP frame: TID=1, PID=0, LEN=7 (UID 1 + PDU 6), UID=1, FC=3, byte-count=4, data=[0,10,0,20]
	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x07, 0x01, 0x03, 0x04, 0x00, 0x0A, 0x00, 0x14}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 4 {
		t.Fatalf("expected 4 data bytes, got %d", len(resp.Data))
	}
}

func TestDecodeModbusException(t *testing.T) {
	c, _ := NewCodec("tcp")
	// TCP frame with exception: FC=0x81 (error), exception code=0x02
	_, err := c.Decode([]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x81, 0x02})
	if err == nil {
		t.Fatal("expected ProtocolError for exception response")
	}
	var pe *kernel.ProtocolError
	if !errors.As(err, &pe) {
		t.Errorf("expected ProtocolError, got %T", err)
	}
}

func TestDecodeTCPFrameTooShort(t *testing.T) {
	c, _ := NewCodec("tcp")
	_, err := c.Decode([]byte{0x00, 0x01})
	if err == nil {
		t.Fatal("expected error for short frame")
	}
}

func TestDecodeRTUFrameTooShort(t *testing.T) {
	c, _ := NewCodec("rtu")
	_, err := c.Decode([]byte{0x01, 0x03})
	if err == nil {
		t.Fatal("expected error for short RTU frame")
	}
}

func TestProtocolInterface(t *testing.T) {
	p := NewProtocol()
	if p.Name() != "modbus" {
		t.Errorf("expected modbus, got %s", p.Name())
	}
	if p.DefaultPort() != 502 {
		t.Errorf("expected 502, got %d", p.DefaultPort())
	}
	variants := p.Variants()
	if len(variants) != 3 {
		t.Errorf("expected 3 variants, got %d", len(variants))
	}
}

func TestProtocolNewCodec(t *testing.T) {
	p := NewProtocol()
	codec, err := p.NewCodec("tcp")
	if err != nil {
		t.Fatalf("expected no error for tcp variant, got %v", err)
	}
	if codec == nil {
		t.Fatal("expected non-nil codec")
	}
	_, err = p.NewCodec("invalid")
	if err == nil {
		t.Fatal("expected error for unsupported variant")
	}
}

func TestCRC16(t *testing.T) {
	// Known vector: 0x01 0x03 0x00 0x00 0x00 0x01 -> CRC = 0x840A
	data := []byte{0x01, 0x03, 0x00, 0x00, 0x00, 0x01}
	crc := crc16(data)
	if crc != 0x0A84 {
		t.Errorf("expected CRC 0x0A84, got 0x%04X", crc)
	}
}

func TestUnsupportedVariant(t *testing.T) {
	_, err := NewCodec("ascii")
	if err == nil {
		t.Fatal("expected error for unsupported variant ascii")
	}
}

func TestModbusExceptionMessages(t *testing.T) {
	tests := []struct {
		code byte
		want string
	}{
		{1, "illegal function"},
		{2, "illegal data address"},
		{3, "illegal data value"},
		{4, "device failure"},
		{6, "device busy"},
		{99, "unknown exception"},
	}
	for _, tt := range tests {
		if got := modbusException(tt.code); got != tt.want {
			t.Errorf("modbusException(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestEncodeWriteMultipleRegisters(t *testing.T) {
	c, _ := NewCodec("tcp")
	req := &kernel.Request{
		Function: "write_multiple_registers",
		Address:  "10",
		Count:    2,
		Data:     []byte{0x00, 0x0A, 0x01, 0x02},
		Metadata: map[string]any{"unit_id": byte(1)},
	}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	// Function code is at byte 7
	if enc[7] != 16 {
		t.Errorf("expected function code 16, got %d", enc[7])
	}
	// Address at PDU bytes 1-2
	addr := uint16(enc[8])<<8 | uint16(enc[9])
	if addr != 10 {
		t.Errorf("expected address 10, got %d", addr)
	}
}
