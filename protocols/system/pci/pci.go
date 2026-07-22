package pci

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type PCIProtocol struct{}

func New() *PCIProtocol { return &PCIProtocol{} }

func (p *PCIProtocol) Name() string        { return "pci" }
func (p *PCIProtocol) Variants() []string   { return []string{"pci"} }
func (p *PCIProtocol) DefaultPort() int     { return 0 }
func (p *PCIProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
