package worldfip

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type WorldFIPProtocol struct{}

func New() *WorldFIPProtocol { return &WorldFIPProtocol{} }

func (p *WorldFIPProtocol) Name() string        { return "worldfip" }
func (p *WorldFIPProtocol) Variants() []string   { return []string{"serial"} }
func (p *WorldFIPProtocol) DefaultPort() int     { return 0 }
func (p *WorldFIPProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
