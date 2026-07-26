// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package ethercat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// EtherCatProtocol implements kernel.Protocol for EtherCAT (Ethernet for Control
// Automation Technology) via the ethercat command-line tool.
type EtherCatProtocol struct{}

// New creates a new EtherCatProtocol.
func New() *EtherCatProtocol { return &EtherCatProtocol{} }

// Name returns the protocol identifier.
func (p *EtherCatProtocol) Name() string { return "ethercat" }

// Variants returns supported transport variants.
func (p *EtherCatProtocol) Variants() []string { return []string{"ethercat", "cmd"} }

// DefaultPort returns the default port (0).
func (p *EtherCatProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *EtherCatProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "ethercat", "cmd":
		return &ethercatCodec{}, nil
	default:
		return nil, fmt.Errorf("ethercat: unsupported variant %q", v)
	}
}

// ethercatCodec encodes/decodes EtherCAT commands via the ethercat CLI tool.
//
// Commands:
//   - upload <addr> <size>   -- read SDO
//   - download <addr> <hex>  -- write SDO
//   - slaves                 -- list slaves
type ethercatCodec struct{}

// Encode serializes an EtherCAT request into an ethercat CLI command.
func (c *ethercatCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "upload":
		addr := req.Address
		if addr == "" {
			addr = "0x0000"
		}
		count := req.Count
		if count == 0 {
			count = 1
		}
		return []byte(fmt.Sprintf("upload %s %d\n", addr, count)), nil

	case "download":
		addr := req.Address
		if addr == "" {
			addr = "0x0000"
		}
		hexData := fmt.Sprintf("0x%X", req.Data)
		return []byte(fmt.Sprintf("download %s %s\n", addr, hexData)), nil

	case "slaves":
		return []byte("slaves\n"), nil

	default:
		return nil, fmt.Errorf("ethercat: unknown function %q", req.Function)
	}
}

// Decode parses ethercat CLI output into a kernel.Response.
func (c *ethercatCodec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

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
