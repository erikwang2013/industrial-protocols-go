package profibus

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type ProfibusProtocol struct{}

func New() *ProfibusProtocol { return &ProfibusProtocol{} }

func (p *ProfibusProtocol) Name() string        { return "profibus" }
func (p *ProfibusProtocol) Variants() []string   { return []string{"dp"} }
func (p *ProfibusProtocol) DefaultPort() int     { return 0 }
func (p *ProfibusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
