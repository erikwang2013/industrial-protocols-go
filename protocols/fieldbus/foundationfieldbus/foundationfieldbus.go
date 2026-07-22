package foundationfieldbus

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type FoundationFieldbusProtocol struct{}

func New() *FoundationFieldbusProtocol { return &FoundationFieldbusProtocol{} }

func (p *FoundationFieldbusProtocol) Name() string        { return "foundationfieldbus" }
func (p *FoundationFieldbusProtocol) Variants() []string   { return []string{"h1","hse"} }
func (p *FoundationFieldbusProtocol) DefaultPort() int     { return 0 }
func (p *FoundationFieldbusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
