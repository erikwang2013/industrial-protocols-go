// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package asinterface

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// AsInterfaceProtocol implements the keinel.Protocol interface for AS-Interface
// gateways (Bihl+Wiedemann / Pepperl+Fuchs).
type AsInterfaceProtocol struct{}

func New() *AsInterfaceProtocol { return &AsInterfaceProtocol{} }

func (p *AsInterfaceProtocol) Name() string       { return "asinterface" }
func (p *AsInterfaceProtocol) Variants() []string { return []string{"gateway"} }
func (p *AsInterfaceProtocol) DefaultPort() int   { return 2003 }

func (p *AsInterfaceProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("asinterface: unsupported variant %q", v)
	}
	return &asinterfaceCodec{}, nil
}

// asinterfaceCodec encodes/decodes AS-Interface ASCII slave address commands.
type asinterfaceCodec struct{}

// Encode builds an AS-Interface text command: <function> <slave> <param>\r\n
func (c *asinterfaceCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := fmt.Sprintf("%s %s\r\n", req.Function, req.Address)
	if len(req.Data) > 0 {
		cmd = fmt.Sprintf("%s %s %x\r\n", req.Function, req.Address, req.Data)
	}
	return []byte(cmd), nil
}

// Decode returns the raw gateway response as a keinel.Response.
func (c *asinterfaceCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
