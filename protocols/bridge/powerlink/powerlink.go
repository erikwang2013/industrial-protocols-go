package powerlink

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// PowerlinkProtocol implements kernel.Protocol for POWERLINK (Ethernet POWERLINK)
// via the openPOWERLINK_demo command-line utility.
type PowerlinkProtocol struct{}

// New creates a new PowerlinkProtocol.
func New() *PowerlinkProtocol { return &PowerlinkProtocol{} }

// Name returns the protocol identifier.
func (p *PowerlinkProtocol) Name() string { return "powerlink" }

// Variants returns supported transport variants.
func (p *PowerlinkProtocol) Variants() []string { return []string{"ethernet", "cmd"} }

// DefaultPort returns the default port (0).
func (p *PowerlinkProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *PowerlinkProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "ethernet", "cmd":
		return &powerlinkCodec{}, nil
	default:
		return nil, fmt.Errorf("powerlink: unsupported variant %q", v)
	}
}

// powerlinkCodec encodes/decodes POWERLINK commands via openPOWERLINK_demo CLI.
//
// Commands:
//   - read <addr> <count>     -- read object dictionary entry
//   - write <addr> <hex>      -- write object dictionary entry
//   - status                  -- query node status
type powerlinkCodec struct{}

// Encode serializes a POWERLINK request into a CLI command.
func (c *powerlinkCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "read":
		addr := req.Address
		if addr == "" {
			addr = "0x1000"
		}
		count := req.Count
		if count == 0 {
			count = 1
		}
		return []byte(fmt.Sprintf("read %s %d\n", addr, count)), nil

	case "write":
		addr := req.Address
		if addr == "" {
			addr = "0x1000"
		}
		hexData := fmt.Sprintf("%X", req.Data)
		return []byte(fmt.Sprintf("write %s %s\n", addr, hexData)), nil

	case "status":
		return []byte("status\n"), nil

	default:
		return nil, fmt.Errorf("powerlink: unknown function %q", req.Function)
	}
}

// Decode parses openPOWERLINK CLI output into a kernel.Response.
func (c *powerlinkCodec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.HasPrefix(line, "NMT_") || strings.HasPrefix(line, "STATE:") {
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
