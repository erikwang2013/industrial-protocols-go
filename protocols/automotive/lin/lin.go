// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package lin

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type LinProtocol struct{}

func New() *LinProtocol           { return &LinProtocol{} }
func (p *LinProtocol) Name() string               { return "lin" }
func (p *LinProtocol) Variants() []string         { return []string{"uart"} }
func (p *LinProtocol) DefaultPort() int           { return 0 }

func (p *LinProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "uart" {
		return nil, fmt.Errorf("lin: unsupported variant %q", v)
	}
	return &linCodec{}, nil
}

type linCodec struct{}

func (c *linCodec) Encode(req *kernel.Request) ([]byte, error) {
	id := byte(0)
	if v, ok := req.Metadata["id"].(byte); ok {
		id = v
	} else if v, ok := req.Metadata["id"].(float64); ok {
		id = byte(v)
	}

	pid := calcPID(id)

	enhanced := false
	if v, ok := req.Metadata["enhanced_checksum"].(bool); ok {
		enhanced = v
	}

	// Sync Break is not representable as bytes (it's a line break)
	// We encode: Sync Byte (0x55) + PID + Data + Checksum
	frame := []byte{0x55, pid}
	frame = append(frame, req.Data...)

	if enhanced {
		cs := calcEnhancedChecksum(pid, req.Data)
		frame = append(frame, cs)
	} else {
		cs := calcClassicChecksum(req.Data)
		frame = append(frame, cs)
	}

	return frame, nil
}

func (c *linCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("lin: frame too short (%d bytes)", len(data))
	}

	// Skip sync byte (0x55)
	if data[0] != 0x55 {
		return &kernel.Response{Metadata: map[string]any{"sync_error": true}}, nil
	}

	pid := data[1]
	id := pid & 0x3F // lower 6 bits

	payload := data[2 : len(data)-1]
	receivedCS := data[len(data)-1]

	// Verify checksum
	expectedCS := calcEnhancedChecksum(pid, payload)
	if expectedCS != receivedCS {
		// Try classic checksum
		expectedCS = calcClassicChecksum(payload)
		if expectedCS != receivedCS {
			return nil, fmt.Errorf("lin: checksum mismatch (got %02X, expected %02X or %02X)", receivedCS, expectedCS, calcClassicChecksum(payload))
		}
	}

	return &kernel.Response{
		Data: payload,
		Metadata: map[string]any{
			"id":  int(id),
			"pid": pid,
		},
	}, nil
}

func calcPID(id byte) byte {
	id = id & 0x3F
	p0 := ((id>>0)&1 ^ (id>>1)&1 ^ (id>>2)&1 ^ (id>>4)&1) << 6
	p1 := (^(id>>1)&1 ^ (id>>3)&1 ^ (id>>4)&1 ^ (id>>5)&1) << 7
	return id | byte(p0) | byte(p1)
}

func calcClassicChecksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	if sum == 0xFF {
		sum = 0xFF
	}
	return ^sum
}

func calcEnhancedChecksum(pid byte, data []byte) byte {
	var sum uint16 = uint16(pid)
	for _, b := range data {
		sum += uint16(b)
	}
	sum = (sum & 0xFF) + (sum >> 8)
	return ^byte(sum)
}
