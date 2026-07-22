package cpci

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type CPCIProtocol struct{}

func New() *CPCIProtocol { return &CPCIProtocol{} }

func (p *CPCIProtocol) Name() string        { return "cpci" }
func (p *CPCIProtocol) Variants() []string   { return []string{"cpci"} }
func (p *CPCIProtocol) DefaultPort() int     { return 0 }
func (p *CPCIProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
