// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package lonworks

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// LonWorksProtocol implements the keinel.Protocol interface for LonWorks
// gateways (Echelon U60/U70 network interface).
type LonWorksProtocol struct{}

func New() *LonWorksProtocol { return &LonWorksProtocol{} }

func (p *LonWorksProtocol) Name() string       { return "lonworks" }
func (p *LonWorksProtocol) Variants() []string { return []string{"gateway"} }
func (p *LonWorksProtocol) DefaultPort() int   { return 2009 }

func (p *LonWorksProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("lonworks: unsupported variant %q", v)
	}
	return &lonworksCodec{}, nil
}

// lonworksCodec encodes/decodes LonWorks ASCII network variable read/write commands.
type lonworksCodec struct{}

// Encode builds a LonWorks text request: <function> <nvi_name>\r\n
func (c *lonworksCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := fmt.Sprintf("%s %s\r\n", req.Function, req.Address)
	if len(req.Data) > 0 {
		cmd = fmt.Sprintf("%s %s %x\r\n", req.Function, req.Address, req.Data)
	}
	return []byte(cmd), nil
}

// Decode returns the raw gateway response as a keinel.Response.
func (c *lonworksCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
