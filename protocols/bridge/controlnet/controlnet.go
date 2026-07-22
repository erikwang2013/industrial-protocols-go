// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package controlnet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// ControlNetProtocol implements kernel.Protocol for ControlNet via the 1784-pcic-cli
// command-line utility.
type ControlNetProtocol struct{}

// New creates a new ControlNetProtocol.
func New() *ControlNetProtocol { return &ControlNetProtocol{} }

// Name returns the protocol identifier.
func (p *ControlNetProtocol) Name() string { return "controlnet" }

// Variants returns supported transport variants.
func (p *ControlNetProtocol) Variants() []string { return []string{"coax", "cmd"} }

// DefaultPort returns the default port (0).
func (p *ControlNetProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *ControlNetProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "coax", "cmd":
		return &controlnetCodec{}, nil
	default:
		return nil, fmt.Errorf("controlnet: unsupported variant %q", v)
	}
}

// controlnetCodec encodes/decodes ControlNet commands via 1784-pcic-cli CLI.
//
// Commands:
//   - read <addr> <count>     -- read from ControlNet device
//   - write <addr> <hex>      -- write to ControlNet device
//   - status                  -- query PCIC card status
type controlnetCodec struct{}

// Encode serializes a ControlNet request into a 1784-pcic-cli command.
func (c *controlnetCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "read":
		addr := req.Address
		if addr == "" {
			addr = "0x00"
		}
		count := req.Count
		if count == 0 {
			count = 1
		}
		return []byte(fmt.Sprintf("read %s %d\n", addr, count)), nil

	case "write":
		addr := req.Address
		if addr == "" {
			addr = "0x00"
		}
		hexData := fmt.Sprintf("%X", req.Data)
		return []byte(fmt.Sprintf("write %s %s\n", addr, hexData)), nil

	case "status":
		return []byte("status\n"), nil

	default:
		return nil, fmt.Errorf("controlnet: unknown function %q", req.Function)
	}
}

// Decode parses 1784-pcic-cli output into a kernel.Response.
func (c *controlnetCodec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.HasPrefix(line, "SLOT:") || strings.HasPrefix(line, "NODE:") || strings.HasPrefix(line, "STATUS:") {
		return &kernel.Response{
			Metadata: map[string]any{"state": line, "raw": line},
		}, nil
	}

	parsed, err := parseHex(line)
	if err == nil && len(parsed) > 0 {
		return &kernel.Response{
			Data:     parsed,
			Metadata: map[string]any{"raw": line},
		}, nil
	}

	return &kernel.Response{
		Data:     data,
		Metadata: map[string]any{"raw": line},
	}, nil
}

func parseHex(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	s = strings.ReplaceAll(s, " ", "")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	if len(s) == 0 {
		return nil, fmt.Errorf("empty")
	}
	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		b, err := strconv.ParseUint(s[i:i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		result[i/2] = byte(b)
	}
	return result, nil
}
