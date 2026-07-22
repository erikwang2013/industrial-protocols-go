package vme

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// VMEProtocol implements kernel.Protocol for VME/VPX bus access via procfs.
type VMEProtocol struct{}

// New creates a new VMEProtocol.
func New() *VMEProtocol { return &VMEProtocol{} }

// Name returns the protocol identifier.
func (p *VMEProtocol) Name() string { return "vme" }

// Variants returns supported transport variants.
func (p *VMEProtocol) Variants() []string { return []string{"vme"} }

// DefaultPort returns the default port (0 for memory-mapped bus).
func (p *VMEProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec that passes through raw bytes to the VME bus.
func (p *VMEProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return &vmeCodec{}, nil
}

// vmeCodec passes raw bytes through for VME bus access.
// VME bus reads/writes happen via the driver's Transport,
// which maps to /proc/vme/<slot>.
type vmeCodec struct{}

// Encode passes the request data through as-is.
func (c *vmeCodec) Encode(req *kernel.Request) ([]byte, error) {
	return req.Data, nil
}

// Decode wraps the raw bytes into a Response.
func (c *vmeCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
