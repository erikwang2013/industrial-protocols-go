// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package hartip

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	hart "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart"
)

type HartIPProtocol struct{}

func New() *HartIPProtocol                { return &HartIPProtocol{} }
func (p *HartIPProtocol) Name() string    { return "hartip" }
func (p *HartIPProtocol) Variants() []string { return []string{"tcp", "udp"} }
func (p *HartIPProtocol) DefaultPort() int   { return 5094 }

func (p *HartIPProtocol) NewCodec(v string) (kernel.Codec, error) {
	return &hartIPCodec{hartProto: &hart.HartProtocol{}}, nil
}

type hartIPCodec struct {
	hartProto kernel.Protocol
}

func (c *hartIPCodec) Encode(req *kernel.Request) ([]byte, error) {
	hartCodec, err := c.hartProto.NewCodec("fsk")
	if err != nil {
		return nil, fmt.Errorf("hartip: %w", err)
	}
	hartFrame, err := hartCodec.Encode(req)
	if err != nil {
		return nil, fmt.Errorf("hartip: %w", err)
	}

	// Strip preamble
	delimIdx := 0
	for i, b := range hartFrame {
		if b != 0xFF {
			delimIdx = i
			break
		}
	}
	hartPayload := hartFrame[delimIdx:]

	// HART-IP header (6 bytes) + length (2 bytes)
	hdr := make([]byte, 8)
	hdr[0] = 1 // version
	hdr[1] = 0 // message type: request
	binary.BigEndian.PutUint16(hdr[6:8], uint16(len(hartPayload)))

	result := append(hdr, hartPayload...)
	return result, nil
}

func (c *hartIPCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("hartip: frame too short (%d bytes)", len(data))
	}
	length := int(binary.BigEndian.Uint16(data[6:8]))
	if len(data) < 8+length {
		return nil, fmt.Errorf("hartip: length mismatch")
	}

	hartCodec, err := c.hartProto.NewCodec("fsk")
	if err != nil {
		return nil, fmt.Errorf("hartip: %w", err)
	}
	// Add preamble for the HART decoder
	payload := append([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, data[8:8+length]...)
	return hartCodec.Decode(payload)
}
