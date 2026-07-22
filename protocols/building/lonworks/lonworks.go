package lonworks

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type LonWorksProtocol struct{}

func New() *LonWorksProtocol { return &LonWorksProtocol{} }

func (p *LonWorksProtocol) Name() string        { return "lonworks" }
func (p *LonWorksProtocol) Variants() []string   { return []string{"serial"} }
func (p *LonWorksProtocol) DefaultPort() int     { return 0 }
func (p *LonWorksProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
