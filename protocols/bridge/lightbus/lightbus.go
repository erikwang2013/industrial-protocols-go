package lightbus

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type LightbusProtocol struct{}

func New() *LightbusProtocol { return &LightbusProtocol{} }

func (p *LightbusProtocol) Name() string        { return "lightbus" }
func (p *LightbusProtocol) Variants() []string   { return []string{"fiber"} }
func (p *LightbusProtocol) DefaultPort() int     { return 0 }
func (p *LightbusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
