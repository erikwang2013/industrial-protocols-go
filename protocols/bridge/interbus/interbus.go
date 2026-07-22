package interbus

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type InterbusProtocol struct{}

func New() *InterbusProtocol { return &InterbusProtocol{} }

func (p *InterbusProtocol) Name() string        { return "interbus" }
func (p *InterbusProtocol) Variants() []string   { return []string{"serial"} }
func (p *InterbusProtocol) DefaultPort() int     { return 0 }
func (p *InterbusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
