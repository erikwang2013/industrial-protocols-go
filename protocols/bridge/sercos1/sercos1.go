// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package sercos1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// Sercos1Protocol implements kernel.Protocol for SERCOS I/II (legacy serial fiber)
// via the sercos_cli command-line utility.
type Sercos1Protocol struct{}

// New creates a new Sercos1Protocol.
func New() *Sercos1Protocol { return &Sercos1Protocol{} }

// Name returns the protocol identifier.
func (p *Sercos1Protocol) Name() string { return "sercos1" }

// Variants returns supported transport variants.
func (p *Sercos1Protocol) Variants() []string { return []string{"fiber", "cmd"} }

// DefaultPort returns the default port (0).
func (p *Sercos1Protocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *Sercos1Protocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "fiber", "cmd":
		return &sercos1Codec{}, nil
	default:
		return nil, fmt.Errorf("sercos1: unsupported variant %q", v)
	}
}

// sercos1Codec encodes/decodes SERCOS I/II commands via sercos_cli CLI.
//
// Commands:
//   - read <addr> <count>   -- read IDN
//   - write <addr> <hex>    -- write IDN
//   - status                -- query drive status
type sercos1Codec struct{}

// Encode serializes a SERCOS I/II request into a sercos_cli command.
func (c *sercos1Codec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "read":
		addr := req.Address
		if addr == "" {
			addr = "0x0000"
		}
		count := req.Count
		if count == 0 {
			count = 1
		}
		return []byte(fmt.Sprintf("read %s %d\n", addr, count)), nil

	case "write":
		addr := req.Address
		if addr == "" {
			addr = "0x0000"
		}
		hexData := fmt.Sprintf("%X", req.Data)
		return []byte(fmt.Sprintf("write %s %s\n", addr, hexData)), nil

	case "status":
		return []byte("status\n"), nil

	default:
		return nil, fmt.Errorf("sercos1: unknown function %q", req.Function)
	}
}

// Decode parses sercos_cli output into a kernel.Response.
func (c *sercos1Codec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.HasPrefix(line, "STATUS:") || strings.HasPrefix(line, "DRIVE:") {
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
