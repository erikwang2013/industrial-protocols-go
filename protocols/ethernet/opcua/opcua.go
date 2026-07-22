// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package opcua

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type OpcuaProtocol struct{}

func New() *OpcuaProtocol                           { return &OpcuaProtocol{} }
func (p *OpcuaProtocol) Name() string               { return "opcua" }
func (p *OpcuaProtocol) Variants() []string         { return []string{"binary"} }
func (p *OpcuaProtocol) DefaultPort() int           { return 4840 }

func (p *OpcuaProtocol) NewCodec(v string) (kernel.Codec, error) {
	return &opcuaCodec{}, nil
}

type opcuaCodec struct{}

func (c *opcuaCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "hel":
		return c.encodeHello(req)
	case "open_secure_channel":
		return c.encodeOpenSecureChannel()
	default:
		return nil, fmt.Errorf("opcua: unknown function %q", req.Function)
	}
}

func (c *opcuaCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("opcua: frame too short (%d bytes)", len(data))
	}
	msgType := string(data[:3])
	switch msgType {
	case "ACK":
		return &kernel.Response{Data: data, Metadata: map[string]any{"function": "ack"}}, nil
	case "OPN":
		return &kernel.Response{Data: data, Metadata: map[string]any{"function": "open_response"}}, nil
	default:
		return &kernel.Response{Data: data}, nil
	}
}

func (c *opcuaCodec) encodeHello(req *kernel.Request) ([]byte, error) {
	endpoint := "opc.tcp://localhost:4840"
	if v, ok := req.Metadata["endpoint"].(string); ok {
		endpoint = v
	}
	buf := &bytes.Buffer{}
	buf.WriteString("HEL")
	buf.WriteByte('F')
	binary.Write(buf, binary.LittleEndian, uint32(0))
	sizePos := buf.Len()
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(65536))
	binary.Write(buf, binary.LittleEndian, uint32(65536))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, int32(len(endpoint)))
	buf.WriteString(endpoint)
	raw := buf.Bytes()
	binary.LittleEndian.PutUint32(raw[sizePos:], uint32(len(raw)))
	return raw, nil
}

func (c *opcuaCodec) encodeOpenSecureChannel() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("OPN")
	buf.WriteByte('F')
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(44))
	for i := 0; i < 7; i++ {
		binary.Write(buf, binary.LittleEndian, uint32(0))
	}
	return buf.Bytes(), nil
}
