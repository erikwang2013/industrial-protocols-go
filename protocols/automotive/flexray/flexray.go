package flexray

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type FlexRayProtocol struct{}

func New() *FlexRayProtocol { return &FlexRayProtocol{} }

func (p *FlexRayProtocol) Name() string        { return "flexray" }
func (p *FlexRayProtocol) Variants() []string   { return []string{"serial"} }
func (p *FlexRayProtocol) DefaultPort() int     { return 0 }
func (p *FlexRayProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
