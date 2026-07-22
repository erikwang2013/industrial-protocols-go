package isa100

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// Isa100Protocol implements kernel.Protocol for ISA100.11a wireless via the
// yfgw410_cli (Yokogawa YFGW410 field wireless gateway) command-line utility.
type Isa100Protocol struct{}

// New creates a new Isa100Protocol.
func New() *Isa100Protocol { return &Isa100Protocol{} }

// Name returns the protocol identifier.
func (p *Isa100Protocol) Name() string { return "isa100" }

// Variants returns supported transport variants.
func (p *Isa100Protocol) Variants() []string { return []string{"wireless", "cmd"} }

// DefaultPort returns the default port (0).
func (p *Isa100Protocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *Isa100Protocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "wireless", "cmd":
		return &isa100Codec{}, nil
	default:
		return nil, fmt.Errorf("isa100: unsupported variant %q", v)
	}
}

// isa100Codec encodes/decodes ISA100.11a wireless device commands via yfgw410_cli.
//
// Commands:
//   - read <tag> <attribute>     -- read device attribute
//   - write <tag> <attr> <hex>   -- write device attribute
//   - list                        -- list provisioned devices
type isa100Codec struct{}

// Encode serializes an ISA100 request into a yfgw410_cli command.
func (c *isa100Codec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "read":
		tag := req.Address
		if tag == "" {
			tag = "DEV001"
		}
		return []byte(fmt.Sprintf("read %s\n", tag)), nil

	case "write":
		tag := req.Address
		if tag == "" {
			tag = "DEV001"
		}
		hexData := fmt.Sprintf("%X", req.Data)
		return []byte(fmt.Sprintf("write %s %s\n", tag, hexData)), nil

	case "list":
		return []byte("list\n"), nil

	default:
		return nil, fmt.Errorf("isa100: unknown function %q", req.Function)
	}
}

// Decode parses yfgw410_cli output into a kernel.Response.
func (c *isa100Codec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.Contains(line, "DEVICE:") || strings.Contains(line, "TAG:") {
		return &kernel.Response{
			Metadata: map[string]any{"device": line, "raw": line},
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
