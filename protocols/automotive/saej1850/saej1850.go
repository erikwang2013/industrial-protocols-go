package saej1850

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type SaeJ1850Protocol struct{}

func New() *SaeJ1850Protocol { return &SaeJ1850Protocol{} }

func (p *SaeJ1850Protocol) Name() string        { return "saej1850" }
func (p *SaeJ1850Protocol) Variants() []string   { return []string{"serial"} }
func (p *SaeJ1850Protocol) DefaultPort() int     { return 0 }
func (p *SaeJ1850Protocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
