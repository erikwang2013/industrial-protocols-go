// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package most

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// MostProtocol implements kernel.Protocol for MOST (Media Oriented Systems Transport)
// over fiber optic with a serial adapter using AT-command interface.
type MostProtocol struct{}

// New creates a new MostProtocol.
func New() *MostProtocol { return &MostProtocol{} }

// Name returns the protocol identifier.
func (p *MostProtocol) Name() string { return "most" }

// Variants returns supported transport variants.
func (p *MostProtocol) Variants() []string { return []string{"optical", "serial"} }

// DefaultPort returns the default port (0 for serial).
func (p *MostProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *MostProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "optical", "serial":
		return &mostCodec{}, nil
	default:
		return nil, fmt.Errorf("most: unsupported variant %q", v)
	}
}

// mostCodec encodes/decodes MOST frames through AT-command serial interface.
//
// AT commands:
//   - AT+READ=<addr>,<len>  -- read from a MOST node
//   - AT+WRITE=<addr>,<hex> -- write data to a MOST node
//   - AT+STATUS             -- query ring status
//
// Response format:
//
//	+OK:<hex_data>
//	+ERR:<code>
type mostCodec struct{}

// Encode serializes a MOST request into an AT command string.
func (c *mostCodec) Encode(req *kernel.Request) ([]byte, error) {
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
		return []byte(fmt.Sprintf("AT+READ=%s,%d\r\n", addr, count)), nil

	case "write":
		addr := req.Address
		if addr == "" {
			addr = "0x0000"
		}
		hexData := fmt.Sprintf("%X", req.Data)
		return []byte(fmt.Sprintf("AT+WRITE=%s,%s\r\n", addr, hexData)), nil

	case "status":
		return []byte("AT+STATUS\r\n"), nil

	default:
		return nil, fmt.Errorf("most: unknown function %q", req.Function)
	}
}

// Decode parses AT command response text into a kernel.Response.
func (c *mostCodec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.HasPrefix(line, "+OK:") {
		payload := strings.TrimPrefix(line, "+OK:")
		// Parse hex data from response
		hexStr := strings.TrimSpace(payload)
		parsed, err := parseHex(hexStr)
		if err != nil {
			return &kernel.Response{
				Metadata: map[string]any{"raw": hexStr},
			}, nil
		}
		return &kernel.Response{
			Data:     parsed,
			Metadata: map[string]any{"raw": hexStr},
		}, nil
	}

	if strings.HasPrefix(line, "+ERR:") {
		code := strings.TrimPrefix(line, "+ERR:")
		return nil, fmt.Errorf("most: adapter error %s", strings.TrimSpace(code))
	}

	// Unrecognized response -- return raw
	return &kernel.Response{
		Data:     data,
		Metadata: map[string]any{"raw": line},
	}, nil
}

// parseHex converts a hex string (e.g. "48 65 6C 6C 6F" or "48656C6C6F") to bytes.
func parseHex(s string) ([]byte, error) {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "0x", "")
	s = strings.ReplaceAll(s, "0X", "")

	// Must be even length
	if len(s)%2 != 0 {
		s = "0" + s
	}

	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		b, err := strconv.ParseUint(s[i:i+2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("most: invalid hex byte %q: %w", s[i:i+2], err)
		}
		result[i/2] = byte(b)
	}
	return result, nil
}
