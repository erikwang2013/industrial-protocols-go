package safej1850

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type SafeJ1850Protocol struct{}

func New() *SafeJ1850Protocol { return &SafeJ1850Protocol{} }

func (p *SafeJ1850Protocol) Name() string        { return "safej1850" }
func (p *SafeJ1850Protocol) Variants() []string   { return []string{"serial"} }
func (p *SafeJ1850Protocol) DefaultPort() int     { return 0 }
func (p *SafeJ1850Protocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
