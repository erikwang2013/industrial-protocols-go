package wirelesshart

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// WirelessHARTProtocol implements kernel.Protocol for WirelessHART via the
// emerson_1410_cli (Emerson 1410/1420 Wireless Gateway) command-line utility.
type WirelessHARTProtocol struct{}

// New creates a new WirelessHARTProtocol.
func New() *WirelessHARTProtocol { return &WirelessHARTProtocol{} }

// Name returns the protocol identifier.
func (p *WirelessHARTProtocol) Name() string { return "wirelesshart" }

// Variants returns supported transport variants.
func (p *WirelessHARTProtocol) Variants() []string { return []string{"wireless", "cmd"} }

// DefaultPort returns the default port (0).
func (p *WirelessHARTProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *WirelessHARTProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "wireless", "cmd":
		return &wirelesshartCodec{}, nil
	default:
		return nil, fmt.Errorf("wirelesshart: unsupported variant %q", v)
	}
}

// wirelesshartCodec encodes/decodes WirelessHART commands via emerson_1410_cli.
//
// Commands:
//   - read <tag> <param>      -- read device parameter
//   - write <tag> <param> <v> -- write device parameter
//   - scan                     -- scan for wireless devices
type wirelesshartCodec struct{}

// Encode serializes a WirelessHART request into an emerson_1410_cli command.
func (c *wirelesshartCodec) Encode(req *kernel.Request) ([]byte, error) {
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

	case "scan":
		return []byte("scan\n"), nil

	default:
		return nil, fmt.Errorf("wirelesshart: unknown function %q", req.Function)
	}
}

// Decode parses emerson_1410_cli output into a kernel.Response.
func (c *wirelesshartCodec) Decode(data []byte) (*kernel.Response, error) {
	line := strings.TrimSpace(string(data))

	if strings.HasPrefix(line, "DEVICE:") || strings.HasPrefix(line, "TAG:") || strings.HasPrefix(line, "NETWORK:") {
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
