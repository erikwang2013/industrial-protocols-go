// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package modbus

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// ModbusProtocol implements the kernel.Protocol interface for Modbus.
type ModbusProtocol struct{}

// NewProtocol returns a new ModbusProtocol.
func NewProtocol() *ModbusProtocol { return &ModbusProtocol{} }

func (p *ModbusProtocol) Name() string      { return "modbus" }
func (p *ModbusProtocol) Variants() []string { return []string{"tcp", "rtu", "ascii"} }
func (p *ModbusProtocol) DefaultPort() int   { return 502 }

func (p *ModbusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	switch variant {
	case "tcp", "rtu":
		return &modbusCodec{variant: variant}, nil
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", variant)
	}
}

// modbusCodec is a kernel.Codec for Modbus TCP and RTU.
type modbusCodec struct {
	variant     string
	transaction uint16
}

// NewCodec returns a modbusCodec for the given variant.
func NewCodec(variant string) (*modbusCodec, error) {
	switch variant {
	case "tcp", "rtu":
		return &modbusCodec{variant: variant}, nil
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", variant)
	}
}

// Encode serializes a kernel.Request into a Modbus frame.
func (c *modbusCodec) Encode(req *kernel.Request) ([]byte, error) {
	mr := c.toModbusReq(req)
	switch c.variant {
	case "tcp":
		return c.encodeTCP(mr), nil
	case "rtu":
		return c.encodeRTU(mr), nil
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", c.variant)
	}
}

// Decode parses a Modbus frame into a kernel.Response.
func (c *modbusCodec) Decode(data []byte) (*kernel.Response, error) {
	switch c.variant {
	case "tcp":
		return c.decodeTCP(data)
	case "rtu":
		return c.decodeRTU(data)
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", c.variant)
	}
}

// modbusRequest is the internal representation of a Modbus request.
type modbusRequest struct {
	function byte
	address  uint16
	count    uint16
	data     []byte
	unitID   byte
}

// toModbusReq converts a kernel.Request to a modbusRequest.
func (c *modbusCodec) toModbusReq(req *kernel.Request) modbusRequest {
	mr := modbusRequest{
		address: parseAddr(req.Address),
		count:   uint16(req.Count),
		data:    req.Data,
	}
	if uid, ok := req.Metadata["unit_id"].(byte); ok {
		mr.unitID = uid
	} else if uid, ok := req.Metadata["unit_id"].(float64); ok {
		mr.unitID = byte(uid)
	}
	switch req.Function {
	case "read_coils":
		mr.function = 1
	case "read_discrete_inputs":
		mr.function = 2
	case "read_holding_registers":
		mr.function = 3
	case "read_input_registers":
		mr.function = 4
	case "write_single_coil":
		mr.function = 5
	case "write_single_register":
		mr.function = 6
	case "write_multiple_registers":
		mr.function = 16
	default:
		mr.function = 3
	}
	return mr
}

// parseAddr parses a string address to uint16.
func parseAddr(s string) uint16 {
	var addr uint32
	fmt.Sscanf(s, "%d", &addr)
	return uint16(addr)
}

// encodeTCP builds a Modbus TCP frame (MBAP header + PDU).
func (c *modbusCodec) encodeTCP(mr modbusRequest) []byte {
	c.transaction++
	pdu := c.encodePDU(mr)
	buf := make([]byte, 7+len(pdu))
	binary.BigEndian.PutUint16(buf[0:2], c.transaction)
	binary.BigEndian.PutUint16(buf[2:4], 0)                       // protocol ID
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(pdu)+1))      // length (UID + PDU)
	buf[6] = mr.unitID
	copy(buf[7:], pdu)
	return buf
}

// decodeTCP parses a Modbus TCP frame.
func (c *modbusCodec) decodeTCP(data []byte) (*kernel.Response, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("modbus: tcp frame too short (%d bytes)", len(data))
	}
	length := binary.BigEndian.Uint16(data[4:6])
	// length counts unit ID + PDU bytes; extract just the PDU.
	if int(length) < 2 || len(data) < 6+int(length) {
		return nil, fmt.Errorf("modbus: tcp frame length mismatch")
	}
	pdu := data[7 : 7+length-1]
	return c.decodePDU(pdu)
}

// encodeRTU builds a Modbus RTU frame (address + PDU + CRC).
func (c *modbusCodec) encodeRTU(mr modbusRequest) []byte {
	pdu := c.encodePDU(mr)
	buf := make([]byte, 1+len(pdu)+2)
	buf[0] = mr.unitID
	copy(buf[1:], pdu)
	crc := crc16(buf[:1+len(pdu)])
	binary.LittleEndian.PutUint16(buf[1+len(pdu):], crc)
	return buf
}

// decodeRTU parses a Modbus RTU frame.
func (c *modbusCodec) decodeRTU(data []byte) (*kernel.Response, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("modbus: rtu frame too short (%d bytes)", len(data))
	}
	pdu := data[1 : len(data)-2]
	return c.decodePDU(pdu)
}

// encodePDU builds the Protocol Data Unit for a Modbus request.
func (c *modbusCodec) encodePDU(mr modbusRequest) []byte {
	switch mr.function {
	case 1, 2, 3, 4: // read coils, discrete inputs, holding regs, input regs
		pdu := make([]byte, 5)
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		binary.BigEndian.PutUint16(pdu[3:5], mr.count)
		return pdu
	case 5: // write single coil
		pdu := make([]byte, 5)
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		if len(mr.data) > 0 && mr.data[0] != 0 {
			pdu[3] = 0xFF
			pdu[4] = 0x00
		}
		return pdu
	case 6: // write single register
		pdu := make([]byte, 5)
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		if len(mr.data) >= 2 {
			copy(pdu[3:5], mr.data[:2])
		}
		return pdu
	case 16: // write multiple registers
		pdu := make([]byte, 6+len(mr.data))
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		binary.BigEndian.PutUint16(pdu[3:5], uint16(len(mr.data)/2))
		pdu[5] = byte(len(mr.data))
		copy(pdu[6:], mr.data)
		return pdu
	default:
		return []byte{mr.function}
	}
}

// decodePDU parses a Modbus response PDU into a kernel.Response.
func (c *modbusCodec) decodePDU(pdu []byte) (*kernel.Response, error) {
	if len(pdu) == 0 {
		return nil, fmt.Errorf("modbus: empty PDU")
	}
	fn := pdu[0]
	if fn&0x80 != 0 {
		if len(pdu) < 2 {
			return nil, fmt.Errorf("modbus: exception PDU too short")
		}
		return nil, &kernel.ProtocolError{
			Code:    fmt.Sprintf("%02X", pdu[1]),
			Message: modbusException(pdu[1]),
			Raw:     pdu,
		}
	}
	switch fn {
	case 1, 2, 3, 4:
		if len(pdu) < 3 {
			return nil, fmt.Errorf("modbus: PDU too short")
		}
		n := int(pdu[1])
		if 2+n > len(pdu) {
			return nil, fmt.Errorf("modbus: byte count %d exceeds PDU length %d", n, len(pdu))
		}
		return &kernel.Response{Address: "0", Data: pdu[2 : 2+n]}, nil
	case 5, 6:
		return &kernel.Response{Address: "0"}, nil
	case 16:
		return &kernel.Response{Address: "0"}, nil
	default:
		return &kernel.Response{Data: pdu}, nil
	}
}

// modbusException returns a human-readable message for a Modbus exception code.
func modbusException(code byte) string {
	switch code {
	case 1:
		return "illegal function"
	case 2:
		return "illegal data address"
	case 3:
		return "illegal data value"
	case 4:
		return "device failure"
	case 6:
		return "device busy"
	default:
		return "unknown exception"
	}
}

// crc16 computes the Modbus CRC-16 over the given data.
func crc16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
