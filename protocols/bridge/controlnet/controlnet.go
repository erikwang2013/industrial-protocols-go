package controlnet

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type ControlNetProtocol struct{}

func New() *ControlNetProtocol { return &ControlNetProtocol{} }

func (p *ControlNetProtocol) Name() string        { return "controlnet" }
func (p *ControlNetProtocol) Variants() []string   { return []string{"coax"} }
func (p *ControlNetProtocol) DefaultPort() int     { return 0 }
func (p *ControlNetProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
