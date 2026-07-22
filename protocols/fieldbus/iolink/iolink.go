// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package iolink

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// IOLinkProtocol implements the keinel.Protocol interface for IO-Link
// masters (ifm / Balluff gateway).
type IOLinkProtocol struct{}

func New() *IOLinkProtocol { return &IOLinkProtocol{} }

func (p *IOLinkProtocol) Name() string       { return "iolink" }
func (p *IOLinkProtocol) Variants() []string { return []string{"gateway"} }
func (p *IOLinkProtocol) DefaultPort() int   { return 2004 }

func (p *IOLinkProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("iolink: unsupported variant %q", v)
	}
	return &iolinkCodec{}, nil
}

// iolinkCodec encodes/decodes IO-Link ASCII device access commands.
type iolinkCodec struct{}

// Encode builds an IO-Link text request: <function> <port>/<subindex>\r\n
func (c *iolinkCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := fmt.Sprintf("%s %s\r\n", req.Function, req.Address)
	if len(req.Data) > 0 {
		cmd = fmt.Sprintf("%s %s %x\r\n", req.Function, req.Address, req.Data)
	}
	return []byte(cmd), nil
}

// Decode returns the raw gateway response as a keinel.Response.
func (c *iolinkCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
