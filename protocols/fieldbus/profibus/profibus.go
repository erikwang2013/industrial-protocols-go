// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package profibus

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// ProfibusProtocol implements the keinel.Protocol interface for PROFIBUS
// gateways (Anybus Communicator / Siemens CP 5611 proxy).
type ProfibusProtocol struct{}

func New() *ProfibusProtocol { return &ProfibusProtocol{} }

func (p *ProfibusProtocol) Name() string       { return "profibus" }
func (p *ProfibusProtocol) Variants() []string { return []string{"gateway"} }
func (p *ProfibusProtocol) DefaultPort() int   { return 2000 }

func (p *ProfibusProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "gateway" {
		return nil, fmt.Errorf("profibus: unsupported variant %q", v)
	}
	return &profibusCodec{}, nil
}

// profibusCodec encodes/decodes PROFIBUS ASCII frames over a GatewayBridge.
type profibusCodec struct{}

// Encode builds a PROFIBUS text request: <function> <address>\r\n
func (c *profibusCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := fmt.Sprintf("%s %s\r\n", req.Function, req.Address)
	data := []byte(cmd)
	return data, nil
}

// Decode returns the raw gateway response as a keinel.Response.
func (c *profibusCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}
