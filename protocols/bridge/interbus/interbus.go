// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package interbus

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	interbusMagic   = 0x1B55
	interbusCmdRead = 0x01
	interbusCmdWrite = 0x02
)

// InterbusProtocol implements the keinel.Protocol interface for Interbus
// gateways (Phoenix Contact IBS).
type InterbusProtocol struct{}

func New() *InterbusProtocol { return &InterbusProtocol{} }

func (p *InterbusProtocol) Name() string       { return "interbus" }
func (p *InterbusProtocol) Variants() []string { return []string{"gateway"} }
func (p *InterbusProtocol) DefaultPort() int   { return 2006 }

func (p *InterbusProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("interbus: unsupported variant %q", v)
	}
	return &interbusCodec{}, nil
}

// interbusCodec encodes/decodes Interbus binary command/response frames.
type interbusCodec struct{}

// Encode builds an Interbus binary command frame.
// Frame: magic(2) | cmd(1) | slave(1) | length(2) | payload(n) | crc(2)
func (c *interbusCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := byte(interbusCmdRead)
	if req.Function == "write" {
		cmd = interbusCmdWrite
	}

	slave := byte(1)
	if v, ok := req.Metadata["slave"].(float64); ok {
		slave = byte(v)
	}

	hdr := make([]byte, 6)
	binary.BigEndian.PutUint16(hdr[0:2], interbusMagic)
	hdr[2] = cmd
	hdr[3] = slave
	binary.BigEndian.PutUint16(hdr[4:6], uint16(len(req.Data)))

	frame := append(hdr, req.Data...)
	crc := checksum8(frame)
	frame = append(frame, crc)

	return frame, nil
}

// Decode parses an Interbus binary response frame.
func (c *interbusCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 7 {
		return nil, fmt.Errorf("interbus: frame too short (%d bytes)", len(data))
	}
	magic := binary.BigEndian.Uint16(data[0:2])
	if magic != interbusMagic {
		return nil, fmt.Errorf("interbus: invalid magic %04X", magic)
	}

	cmd := data[2]
	slave := data[3]
	length := binary.BigEndian.Uint16(data[4:6])

	payloadEnd := 6 + int(length)
	if payloadEnd+1 > len(data) {
		return nil, fmt.Errorf("interbus: data length mismatch")
	}
	payload := data[6:payloadEnd]

	expectedCRC := data[payloadEnd]
	computedCRC := checksum8(data[:payloadEnd])
	if expectedCRC != computedCRC {
		return nil, fmt.Errorf("interbus: checksum mismatch")
	}

	return &kernel.Response{
		Data: payload,
		Metadata: map[string]any{
			"slave":   int(slave),
			"command": int(cmd),
		},
	}, nil
}

func checksum8(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum ^= b
	}
	return sum
}
