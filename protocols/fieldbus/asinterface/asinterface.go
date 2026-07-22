package asinterface

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type AsInterfaceProtocol struct{}

func New() *AsInterfaceProtocol { return &AsInterfaceProtocol{} }

func (p *AsInterfaceProtocol) Name() string        { return "asinterface" }
func (p *AsInterfaceProtocol) Variants() []string   { return []string{"serial"} }
func (p *AsInterfaceProtocol) DefaultPort() int     { return 0 }
func (p *AsInterfaceProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
