package canopen

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type CanOpenProtocol struct{}

func New() *CanOpenProtocol { return &CanOpenProtocol{} }

func (p *CanOpenProtocol) Name() string        { return "canopen" }
func (p *CanOpenProtocol) Variants() []string   { return []string{"can"} }
func (p *CanOpenProtocol) DefaultPort() int     { return 0 }
func (p *CanOpenProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
