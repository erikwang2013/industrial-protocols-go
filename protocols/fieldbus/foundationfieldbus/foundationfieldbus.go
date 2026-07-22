package foundationfieldbus

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// FoundationFieldbusProtocol implements the keinel.Protocol interface for
// Foundation Fieldbus (H1/HSE) via a gateway bridge (NI USB-8486 / Softing FFusb).
type FoundationFieldbusProtocol struct{}

func New() *FoundationFieldbusProtocol { return &FoundationFieldbusProtocol{} }

func (p *FoundationFieldbusProtocol) Name() string       { return "foundationfieldbus" }
func (p *FoundationFieldbusProtocol) Variants() []string { return []string{"gateway"} }
func (p *FoundationFieldbusProtocol) DefaultPort() int   { return 2002 }

func (p *FoundationFieldbusProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("foundationfieldbus: unsupported variant %q", v)
	}
	return &foundationfieldbusCodec{}, nil
}

// foundationfieldbusCodec encodes/decodes Foundation Fieldbus ASCII/hex command frames.
type foundationfieldbusCodec struct{}

// Encode builds a Foundation Fieldbus text request: <function> <tag>\r\n
func (c *foundationfieldbusCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := fmt.Sprintf("%s %s\r\n", req.Function, req.Address)
	if len(req.Data) > 0 {
		cmd = fmt.Sprintf("%s %s %x\r\n", req.Function, req.Address, req.Data)
	}
	return []byte(cmd), nil
}

// Decode returns the raw gateway response as a keinel.Response.
func (c *foundationfieldbusCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
