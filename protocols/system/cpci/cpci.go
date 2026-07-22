package cpci

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// CPCIProtocol implements kernel.Protocol for CompactPCI bus access via sysfs.
type CPCIProtocol struct{}

// New creates a new CPCIProtocol.
func New() *CPCIProtocol { return &CPCIProtocol{} }

// Name returns the protocol identifier.
func (p *CPCIProtocol) Name() string { return "cpci" }

// Variants returns supported transport variants.
func (p *CPCIProtocol) Variants() []string { return []string{"cpci"} }

// DefaultPort returns the default port (0 for memory-mapped bus).
func (p *CPCIProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec that passes through raw bytes to the CPCI config space.
func (p *CPCIProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return &cpciCodec{}, nil
}

// cpciCodec passes raw bytes through for CompactPCI config space access.
// CompactPCI uses the same sysfs PCI interface: /sys/bus/pci/devices/<bdf>/config.
// The codec is protocol-agnostic and forwards data transparently.
type cpciCodec struct{}

// Encode passes the request data through as-is.
func (c *cpciCodec) Encode(req *kernel.Request) ([]byte, error) {
	return req.Data, nil
}

// Decode wraps the raw bytes into a Response.
func (c *cpciCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
