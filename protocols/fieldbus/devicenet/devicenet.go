package devicenet

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type DeviceNetProtocol struct{}

func New() *DeviceNetProtocol { return &DeviceNetProtocol{} }

func (p *DeviceNetProtocol) Name() string        { return "devicenet" }
func (p *DeviceNetProtocol) Variants() []string   { return []string{"can"} }
func (p *DeviceNetProtocol) DefaultPort() int     { return 0 }
func (p *DeviceNetProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
