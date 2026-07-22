package vme

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type VMEProtocol struct{}

func New() *VMEProtocol { return &VMEProtocol{} }

func (p *VMEProtocol) Name() string        { return "vme" }
func (p *VMEProtocol) Variants() []string   { return []string{"vme"} }
func (p *VMEProtocol) DefaultPort() int     { return 0 }
func (p *VMEProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
