package modbusplus

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type ModbusPlusProtocol struct{}

func New() *ModbusPlusProtocol { return &ModbusPlusProtocol{} }

func (p *ModbusPlusProtocol) Name() string        { return "modbusplus" }
func (p *ModbusPlusProtocol) Variants() []string   { return []string{"serial"} }
func (p *ModbusPlusProtocol) DefaultPort() int     { return 0 }
func (p *ModbusPlusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
