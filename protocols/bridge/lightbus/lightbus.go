package lightbus

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	lightbusMagic   = 0x4C42
	lightbusCmdRead = 0x01
	lightbusCmdWrite = 0x02
)

// LightbusProtocol implements the keinel.Protocol interface for Lightbus
// gateways (Beckhoff fiber optic ring).
type LightbusProtocol struct{}

func New() *LightbusProtocol { return &LightbusProtocol{} }

func (p *LightbusProtocol) Name() string       { return "lightbus" }
func (p *LightbusProtocol) Variants() []string { return []string{"gateway"} }
func (p *LightbusProtocol) DefaultPort() int   { return 2008 }

func (p *LightbusProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("lightbus: unsupported variant %q", v)
	}
	return &lightbusCodec{}, nil
}

// lightbusCodec encodes/decodes Lightbus fiber optic binary frames.
type lightbusCodec struct{}

// Encode builds a Lightbus fiber frame.
// Frame: magic(2) | cmd(1) | module(1) | channel(2) | length(2) | payload(n) | crc(2)
func (c *lightbusCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := byte(lightbusCmdRead)
	if req.Function == "write" {
		cmd = lightbusCmdWrite
	}

	module := byte(1)
	channel := uint16(0)
	if v, ok := req.Metadata["module"].(float64); ok {
		module = byte(v)
	}
	if v, ok := req.Metadata["channel"].(float64); ok {
		channel = uint16(v)
	}

	hdr := make([]byte, 8)
	binary.BigEndian.PutUint16(hdr[0:2], lightbusMagic)
	hdr[2] = cmd
	hdr[3] = module
	binary.BigEndian.PutUint16(hdr[4:6], channel)
	binary.BigEndian.PutUint16(hdr[6:8], uint16(len(req.Data)))

	frame := append(hdr, req.Data...)
	crc := checksum8XOR(frame)
	frame = append(frame, crc)

	return frame, nil
}

// Decode parses a Lightbus fiber response frame.
func (c *lightbusCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("lightbus: frame too short (%d bytes)", len(data))
	}
	magic := binary.BigEndian.Uint16(data[0:2])
	if magic != lightbusMagic {
		return nil, fmt.Errorf("lightbus: invalid magic %04X", magic)
	}

	cmd := data[2]
	module := data[3]
	channel := binary.BigEndian.Uint16(data[4:6])
	length := binary.BigEndian.Uint16(data[6:8])

	payloadEnd := 8 + int(length)
	if payloadEnd+1 > len(data) {
		return nil, fmt.Errorf("lightbus: data length mismatch")
	}
	payload := data[8:payloadEnd]

	expectedCRC := data[payloadEnd]
	computedCRC := checksum8XOR(data[:payloadEnd])
	if expectedCRC != computedCRC {
		return nil, fmt.Errorf("lightbus: checksum mismatch")
	}

	return &kernel.Response{
		Data: payload,
		Metadata: map[string]any{
			"module":  int(module),
			"channel": int(channel),
			"command": int(cmd),
		},
	}, nil
}

func checksum8XOR(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum ^= b
	}
	return sum
}
