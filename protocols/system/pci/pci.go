package pci

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// PCIProtocol implements kernel.Protocol for PCI/PCIe bus access via sysfs.
type PCIProtocol struct{}

// New creates a new PCIProtocol.
func New() *PCIProtocol { return &PCIProtocol{} }

// Name returns the protocol identifier.
func (p *PCIProtocol) Name() string { return "pci" }

// Variants returns supported transport variants.
func (p *PCIProtocol) Variants() []string { return []string{"pci"} }

// DefaultPort returns the default port (0 for memory-mapped bus).
func (p *PCIProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec that passes through raw bytes to the PCI config space.
func (p *PCIProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return &pciCodec{}, nil
}

// pciCodec passes raw bytes through for PCI config space access.
// PCI config space reads/writes happen via the driver's Transport,
// which maps to /sys/bus/pci/devices/<bdf>/config.
type pciCodec struct{}

// Encode passes the request data through as-is.
func (c *pciCodec) Encode(req *kernel.Request) ([]byte, error) {
	return req.Data, nil
}

// Decode wraps the raw bytes into a Response.
func (c *pciCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
