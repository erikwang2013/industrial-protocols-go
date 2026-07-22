package powerlink

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type PowerlinkProtocol struct{}

func New() *PowerlinkProtocol { return &PowerlinkProtocol{} }

func (p *PowerlinkProtocol) Name() string        { return "powerlink" }
func (p *PowerlinkProtocol) Variants() []string   { return []string{"ethernet"} }
func (p *PowerlinkProtocol) DefaultPort() int     { return 0 }
func (p *PowerlinkProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
