package most

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type MostProtocol struct{}

func New() *MostProtocol { return &MostProtocol{} }

func (p *MostProtocol) Name() string        { return "most" }
func (p *MostProtocol) Variants() []string   { return []string{"optical"} }
func (p *MostProtocol) DefaultPort() int     { return 0 }
func (p *MostProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
