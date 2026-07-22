// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cclink

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	StartDelim = 0x7E
	CtrlRead   = 0x01
	CtrlWrite  = 0x02
)

type CCLinkProtocol struct{}

func New() *CCLinkProtocol                       { return &CCLinkProtocol{} }
func (p *CCLinkProtocol) Name() string           { return "cclink" }
func (p *CCLinkProtocol) Variants() []string     { return []string{"rs485"} }
func (p *CCLinkProtocol) DefaultPort() int       { return 0 }

func (p *CCLinkProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "rs485" {
		return nil, fmt.Errorf("cclink: unsupported variant %q", v)
	}
	return &cclinkCodec{}, nil
}

type cclinkCodec struct{}

func (c *cclinkCodec) Encode(req *kernel.Request) ([]byte, error) {
	station := byte(1)
	if v, ok := req.Metadata["station"].(byte); ok {
		station = v
	}
	if v, ok := req.Metadata["station"].(float64); ok {
		station = byte(v)
	}

	ctrl := byte(CtrlRead)
	if req.Function == "write" {
		ctrl = CtrlWrite
	}

	dataLen := byte(len(req.Data))

	frame := []byte{StartDelim, station, ctrl, dataLen}
	frame = append(frame, req.Data...)
	crc := crc16XMODEM(frame[1:]) // CRC over all bytes between start and end delimiter
	crcEnd := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcEnd, crc)
	frame = append(frame, crcEnd...)
	frame = append(frame, StartDelim)

	return frame, nil
}

func (c *cclinkCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 7 {
		return nil, fmt.Errorf("cclink: frame too short (%d bytes)", len(data))
	}
	if data[0] != StartDelim {
		return nil, fmt.Errorf("cclink: missing start delimiter")
	}

	station := data[1]
	ctrl := data[2]
	dataLen := data[3]

	// Extract payload (before CRC and end delimiter)
	payloadEnd := 4 + int(dataLen)
	if payloadEnd+3 > len(data) {
		return nil, fmt.Errorf("cclink: data length mismatch")
	}
	payload := data[4:payloadEnd]

	// Verify CRC
	expectedCRC := binary.LittleEndian.Uint16(data[payloadEnd : payloadEnd+2])
	computedCRC := crc16XMODEM(data[1:payloadEnd])
	if expectedCRC != computedCRC {
		return nil, fmt.Errorf("cclink: CRC mismatch (expected %04X, computed %04X)", expectedCRC, computedCRC)
	}

	return &kernel.Response{
		Address: fmt.Sprintf("%d", station),
		Data:    payload,
		Metadata: map[string]any{
			"station": int(station),
			"control": int(ctrl),
		},
	}, nil
}

func crc16XMODEM(data []byte) uint16 {
	var crc uint16 = 0x0000
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
