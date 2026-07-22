package sercos

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type SercosProtocol struct{}

func New() *SercosProtocol { return &SercosProtocol{} }

func (p *SercosProtocol) Name() string        { return "sercos" }
func (p *SercosProtocol) Variants() []string   { return []string{"fiber"} }
func (p *SercosProtocol) DefaultPort() int     { return 0 }
func (p *SercosProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
