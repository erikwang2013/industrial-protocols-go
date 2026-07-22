package modbusplus

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	mbpMagic    = 0x4D42
	mbpCmdRead  = 0x01
	mbpCmdWrite = 0x02
)

// ModbusPlusProtocol implements the keinel.Protocol interface for Modbus Plus
// gateways (Schneider SA85/BM85 bridge).
type ModbusPlusProtocol struct{}

func New() *ModbusPlusProtocol { return &ModbusPlusProtocol{} }

func (p *ModbusPlusProtocol) Name() string       { return "modbusplus" }
func (p *ModbusPlusProtocol) Variants() []string { return []string{"gateway", "cmd"} }
func (p *ModbusPlusProtocol) DefaultPort() int   { return 2010 }

func (p *ModbusPlusProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "gateway", "cmd":
		return &modbusplusCodec{}, nil
	default:
		return nil, fmt.Errorf("modbusplus: unsupported variant %q", v)
	}
}

// modbusplusCodec encodes/decodes Modbus Plus token-passing binary frames.
type modbusplusCodec struct{}

// Encode builds a Modbus Plus binary frame.
// Frame: magic(2) | dest(1) | cmd(1) | length(2) | payload(n) | crc(2)
func (c *modbusplusCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := byte(mbpCmdRead)
	if req.Function == "write" {
		cmd = mbpCmdWrite
	}

	dest := byte(1)
	if v, ok := req.Metadata["dest"].(float64); ok {
		dest = byte(v)
	}

	hdr := make([]byte, 6)
	binary.BigEndian.PutUint16(hdr[0:2], mbpMagic)
	hdr[2] = dest
	hdr[3] = cmd
	binary.BigEndian.PutUint16(hdr[4:6], uint16(len(req.Data)))

	frame := append(hdr, req.Data...)
	crc := modbusCRC16(frame)
	crcEnd := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcEnd, crc)
	frame = append(frame, crcEnd...)

	return frame, nil
}

// Decode parses a Modbus Plus binary response frame.
func (c *modbusplusCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("modbusplus: frame too short (%d bytes)", len(data))
	}
	magic := binary.BigEndian.Uint16(data[0:2])
	if magic != mbpMagic {
		return nil, fmt.Errorf("modbusplus: invalid magic %04X", magic)
	}

	dest := data[2]
	cmd := data[3]
	length := binary.BigEndian.Uint16(data[4:6])

	payloadEnd := 6 + int(length)
	if payloadEnd+2 > len(data) {
		return nil, fmt.Errorf("modbusplus: data length mismatch")
	}
	payload := data[6:payloadEnd]

	expectedCRC := binary.LittleEndian.Uint16(data[payloadEnd : payloadEnd+2])
	computedCRC := modbusCRC16(data[:payloadEnd])
	if expectedCRC != computedCRC {
		return nil, fmt.Errorf("modbusplus: CRC mismatch")
	}

	return &kernel.Response{
		Data: payload,
		Metadata: map[string]any{
			"dest":    int(dest),
			"command": int(cmd),
		},
	}, nil
}

func modbusCRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
