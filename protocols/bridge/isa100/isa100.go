package isa100

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type Isa100Protocol struct{}

func New() *Isa100Protocol { return &Isa100Protocol{} }

func (p *Isa100Protocol) Name() string        { return "isa100" }
func (p *Isa100Protocol) Variants() []string   { return []string{"wireless"} }
func (p *Isa100Protocol) DefaultPort() int     { return 0 }
func (p *Isa100Protocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
