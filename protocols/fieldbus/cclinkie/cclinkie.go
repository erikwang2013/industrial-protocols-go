// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cclinkie

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	cclinkIEMagic    = 0xCC1E
	cclinkIECmdRead  = 0x01
	cclinkIECmdWrite = 0x02
)

// CCLinkIEProtocol implements the kernel.Protocol interface for CC-Link IE
// Field gateways (Mitsubishi MELSEC gateway).
type CCLinkIEProtocol struct{}

func New() *CCLinkIEProtocol { return &CCLinkIEProtocol{} }

func (p *CCLinkIEProtocol) Name() string       { return "cclinkie" }
func (p *CCLinkIEProtocol) Variants() []string { return []string{"gateway"} }
func (p *CCLinkIEProtocol) DefaultPort() int   { return 2005 }

func (p *CCLinkIEProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("cclinkie: unsupported variant %q", v)
	}
	return &cclinkieCodec{}, nil
}

// cclinkieCodec encodes/decodes CC-Link IE Field binary frames.
type cclinkieCodec struct{}

// Encode builds a binary CC-Link IE header + payload frame.
// Frame: magic(2) | cmd(1) | station(2) | length(2) | payload(n) | crc(2)
func (c *cclinkieCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := byte(cclinkIECmdRead)
	if req.Function == "write" {
		cmd = cclinkIECmdWrite
	}

	station := uint16(1)
	if v, ok := req.Metadata["station"].(float64); ok {
		station = uint16(v)
	}
	if v, ok := req.Metadata["station"].(uint16); ok {
		station = v
	}

	hdr := make([]byte, 7)
	binary.BigEndian.PutUint16(hdr[0:2], cclinkIEMagic)
	hdr[2] = cmd
	binary.BigEndian.PutUint16(hdr[3:5], station)
	binary.BigEndian.PutUint16(hdr[5:7], uint16(len(req.Data)))

	frame := append(hdr, req.Data...)
	crc := crc16CCITT(frame)
	crcEnd := make([]byte, 2)
	binary.BigEndian.PutUint16(crcEnd, crc)
	frame = append(frame, crcEnd...)

	return frame, nil
}

// Decode parses a binary CC-Link IE response frame.
func (c *cclinkieCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("cclinkie: frame too short (%d bytes)", len(data))
	}
	magic := binary.BigEndian.Uint16(data[0:2])
	if magic != cclinkIEMagic {
		return nil, fmt.Errorf("cclinkie: invalid magic %04X", magic)
	}

	cmd := data[2]
	station := binary.BigEndian.Uint16(data[3:5])
	length := binary.BigEndian.Uint16(data[5:7])

	payloadEnd := 7 + int(length)
	if payloadEnd+2 > len(data) {
		return nil, fmt.Errorf("cclinkie: data length mismatch")
	}
	payload := data[7:payloadEnd]

	expectedCRC := binary.BigEndian.Uint16(data[payloadEnd : payloadEnd+2])
	computedCRC := crc16CCITT(data[:payloadEnd])
	if expectedCRC != computedCRC {
		return nil, fmt.Errorf("cclinkie: CRC mismatch")
	}

	return &kernel.Response{
		Data: payload,
		Metadata: map[string]any{
			"station": int(station),
			"command": int(cmd),
		},
	}, nil
}

func crc16CCITT(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}
