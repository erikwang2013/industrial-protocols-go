package worldfip

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	worldFIPMagic    = 0x57F1
	worldFIPProduce  = 0x01
	worldFIPConsume  = 0x02
)

// WorldFIPProtocol implements the keinel.Protocol interface for WorldFIP
// gateways (FIP gateway / FIPIO agent).
type WorldFIPProtocol struct{}

func New() *WorldFIPProtocol { return &WorldFIPProtocol{} }

func (p *WorldFIPProtocol) Name() string       { return "worldfip" }
func (p *WorldFIPProtocol) Variants() []string { return []string{"gateway"} }
func (p *WorldFIPProtocol) DefaultPort() int   { return 2007 }

func (p *WorldFIPProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("worldfip: unsupported variant %q", v)
	}
	return &worldfipCodec{}, nil
}

// worldfipCodec encodes/decodes WorldFIP produce/consume binary frames.
type worldfipCodec struct{}

// Encode builds a WorldFIP produce/consume binary frame.
// Frame: magic(2) | op(1) | tag(2) | length(2) | payload(n) | crc(2)
func (c *worldfipCodec) Encode(req *kernel.Request) ([]byte, error) {
	op := byte(worldFIPConsume)
	if req.Function == "write" || req.Function == "produce" {
		op = worldFIPProduce
	}

	tag := uint16(0)
	if req.Address != "" {
		tag = uint16(len(req.Address))
	}
	if v, ok := req.Metadata["tag"].(float64); ok {
		tag = uint16(v)
	}

	hdr := make([]byte, 7)
	binary.BigEndian.PutUint16(hdr[0:2], worldFIPMagic)
	hdr[2] = op
	binary.BigEndian.PutUint16(hdr[3:5], tag)
	binary.BigEndian.PutUint16(hdr[5:7], uint16(len(req.Data)))

	frame := append(hdr, req.Data...)
	crc := xorCRC16(frame)
	crcEnd := make([]byte, 2)
	binary.BigEndian.PutUint16(crcEnd, crc)
	frame = append(frame, crcEnd...)

	return frame, nil
}

// Decode parses a WorldFIP binary response frame.
func (c *worldfipCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("worldfip: frame too short (%d bytes)", len(data))
	}
	magic := binary.BigEndian.Uint16(data[0:2])
	if magic != worldFIPMagic {
		return nil, fmt.Errorf("worldfip: invalid magic %04X", magic)
	}

	op := data[2]
	tag := binary.BigEndian.Uint16(data[3:5])
	length := binary.BigEndian.Uint16(data[5:7])

	payloadEnd := 7 + int(length)
	if payloadEnd+2 > len(data) {
		return nil, fmt.Errorf("worldfip: data length mismatch")
	}
	payload := data[7:payloadEnd]

	expectedCRC := binary.BigEndian.Uint16(data[payloadEnd : payloadEnd+2])
	computedCRC := xorCRC16(data[:payloadEnd])
	if expectedCRC != computedCRC {
		return nil, fmt.Errorf("worldfip: CRC mismatch")
	}

	return &kernel.Response{
		Data: payload,
		Metadata: map[string]any{
			"tag":      int(tag),
			"operation": int(op),
		},
	}, nil
}

func xorCRC16(data []byte) uint16 {
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
