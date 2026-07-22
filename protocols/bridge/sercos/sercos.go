// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package sercos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// SercosProtocol implements kernel.Protocol for SERCOS III via the netx_cli
// command-line utility.
type SercosProtocol struct{}

// New creates a new SercosProtocol.
func New() *SercosProtocol { return &SercosProtocol{} }

// Name returns the protocol identifier.
func (p *SercosProtocol) Name() string { return "sercos" }

// Variants returns supported transport variants.
func (p *SercosProtocol) Variants() []string { return []string{"fiber", "cmd"} }

// DefaultPort returns the default port (0).
func (p *SercosProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *SercosProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "fiber", "cmd":
		return &sercosCodec{}, nil
	default:
		return nil, fmt.Errorf("sercos: unsupported variant %q", v)
	}
}

// sercosCodec encodes/decodes SERCOS III commands via netx_cli CLI.
//
// Commands:
//   - read <addr> <count>   -- read IDN/S parameter
//   - write <addr> <hex>    -- write IDN/S parameter
//   - phase <N>             -- set communication phase
type sercosCodec struct{}

// Encode serializes a SERCOS III request into a netx_cli command.
func (c *sercosCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "read":
		addr := req.Address
		if addr == "" {
			addr = "S-0-0"
		}
		count := req.Count
		if count == 0 {
			count = 1
		}
		return []byte(fmt.Sprintf("read %s %d\n", addr, count)), nil

	case "write":
		addr := req.Address
		if addr == "" {
			addr = "S-0-0"
		}
		hexData := fmt.Sprintf("%X", req.Data)
		return []byte(fmt.Sprintf("write %s %s\n", addr, hexData)), nil

	case "phase":
		phase := "0"
		if req.Address != "" {
			phase = req.Address
		}
		return []byte(fmt.Sprintf("phase %s\n", phase)), nil

	default:
		return nil, fmt.Errorf("sercos: unknown function %q", req.Function)
	}
}

// Decode parses netx_cli output into a kernel.Response.
func (c *sercosCodec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.HasPrefix(line, "PHASE:") || strings.HasPrefix(line, "STATE:") {
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
