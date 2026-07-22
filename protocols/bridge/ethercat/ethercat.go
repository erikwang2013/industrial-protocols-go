package ethercat

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type EtherCatProtocol struct{}

func New() *EtherCatProtocol { return &EtherCatProtocol{} }

func (p *EtherCatProtocol) Name() string        { return "ethercat" }
func (p *EtherCatProtocol) Variants() []string   { return []string{"ethercat"} }
func (p *EtherCatProtocol) DefaultPort() int     { return 0 }
func (p *EtherCatProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
