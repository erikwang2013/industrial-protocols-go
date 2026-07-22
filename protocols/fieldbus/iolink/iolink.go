package iolink

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type IOLinkProtocol struct{}

func New() *IOLinkProtocol { return &IOLinkProtocol{} }

func (p *IOLinkProtocol) Name() string        { return "iolink" }
func (p *IOLinkProtocol) Variants() []string   { return []string{"serial"} }
func (p *IOLinkProtocol) DefaultPort() int     { return 0 }
func (p *IOLinkProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
