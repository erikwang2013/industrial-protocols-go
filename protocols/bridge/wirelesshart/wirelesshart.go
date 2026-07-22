package wirelesshart

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type WirelessHARTProtocol struct{}

func New() *WirelessHARTProtocol { return &WirelessHARTProtocol{} }

func (p *WirelessHARTProtocol) Name() string        { return "wirelesshart" }
func (p *WirelessHARTProtocol) Variants() []string   { return []string{"wireless"} }
func (p *WirelessHARTProtocol) DefaultPort() int     { return 0 }
func (p *WirelessHARTProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
