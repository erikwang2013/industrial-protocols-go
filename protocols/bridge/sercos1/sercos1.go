package sercos1

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type Sercos1Protocol struct{}

func New() *Sercos1Protocol { return &Sercos1Protocol{} }

func (p *Sercos1Protocol) Name() string        { return "sercos1" }
func (p *Sercos1Protocol) Variants() []string   { return []string{"fiber"} }
func (p *Sercos1Protocol) DefaultPort() int     { return 0 }
func (p *Sercos1Protocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
