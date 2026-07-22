package dali

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	CmdOff         = 0x00
	CmdMax         = 0x05
	CmdRecallMin   = 0x06
	CmdDimUp       = 0x08
	CmdDimDown     = 0x09
	CmdQueryStatus = 0xA1
	BroadcastAddr  = 0xFE
)

type DaliProtocol struct{}

func New() *DaliProtocol                         { return &DaliProtocol{} }
func (p *DaliProtocol) Name() string             { return "dali" }
func (p *DaliProtocol) Variants() []string       { return []string{"serial"} }
func (p *DaliProtocol) DefaultPort() int         { return 0 }

func (p *DaliProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "serial" {
		return nil, fmt.Errorf("dali: unsupported variant %q", v)
	}
	return &daliCodec{}, nil
}

type daliCodec struct{}

func (c *daliCodec) Encode(req *kernel.Request) ([]byte, error) {
	var addr, cmd byte

	if v, ok := req.Metadata["address"].(byte); ok {
		addr = v
	} else {
		addr = BroadcastAddr
	}

	switch req.Function {
	case "off":
		cmd = CmdOff
	case "max":
		cmd = CmdMax
	case "recall_min":
		cmd = CmdRecallMin
	case "dim_up":
		cmd = CmdDimUp
	case "dim_down":
		cmd = CmdDimDown
	case "query_status":
		cmd = CmdQueryStatus
	case "direct_arc":
		if len(req.Data) > 0 {
			cmd = req.Data[0]
		}
	default:
		cmd = CmdOff
	}

	// Forward frame: address byte + command byte = 16 bits
	return []byte{addr, cmd}, nil
}

func (c *daliCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) == 0 {
		return &kernel.Response{}, nil
	}
	// Backward frame: 8-bit status
	status := data[0]
	return &kernel.Response{
		Data: data,
		Metadata: map[string]any{
			"status":       status,
			"lamp_failure": status&0x01 != 0,
			"lamp_on":      status&0x02 != 0,
		},
	}, nil
}
