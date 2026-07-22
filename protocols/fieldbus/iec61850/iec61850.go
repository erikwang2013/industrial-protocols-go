// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package iec61850

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	ServiceInitiate = 0xA8
	ServiceConclude = 0xA9
	ServiceRead     = 0xA4
	ServiceWrite    = 0xA5

	TagBool          = 0x83
	TagInteger       = 0x85
	TagUnsigned      = 0x86
	TagFloat         = 0x87
	TagVisibleString = 0x8A
	TagOctetString   = 0x89
	TagStruct        = 0xA2
)

type Iec61850Protocol struct{}

func New() *Iec61850Protocol                    { return &Iec61850Protocol{} }
func (p *Iec61850Protocol) Name() string        { return "iec61850" }
func (p *Iec61850Protocol) Variants() []string  { return []string{"mms"} }
func (p *Iec61850Protocol) DefaultPort() int    { return 102 }

func (p *Iec61850Protocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "mms" {
		return nil, fmt.Errorf("iec61850: unsupported variant %q", v)
	}
	return &iecCodec{}, nil
}

type iecCodec struct{ invokeID uint32 }

func (c *iecCodec) Encode(req *kernel.Request) ([]byte, error) {
	c.invokeID++
	switch req.Function {
	case "initiate":
		return c.encodeInitiate()
	case "conclude":
		return c.encodeConclude()
	case "read":
		return c.encodeRead(req.Address)
	case "write":
		return c.encodeWrite(req)
	default:
		return nil, fmt.Errorf("iec61850: unknown function %q", req.Function)
	}
}

func (c *iecCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("iec61850: frame too short")
	}

	tag := data[0]
	switch tag {
	case 0xA2: // Initiate Response
		return &kernel.Response{Function: "initiate_response", Data: data}, nil
	case 0xA3: // Conclude Response
		return &kernel.Response{Function: "conclude_response", Data: data}, nil
	case 0xA4: // Read Response
		value, err := decodeValue(data[2:])
		if err != nil {
			return nil, err
		}
		return &kernel.Response{Function: "read_response", Data: value}, nil
	default:
		return &kernel.Response{Data: data, Metadata: map[string]any{"tag": int(tag)}}, nil
	}
}

func (c *iecCodec) encodeInitiate() ([]byte, error) {
	// InitiateRequestPDU = 0xA8 (initiate-RequestPDU) + length + InitRequestDetail
	detail := []byte{0x80, 0x01, 50, 0x81, 0x01, 50, 0x82, 0x01, 50} // proposed MMS version
	pdu := []byte{0xA8, byte(2 + len(detail)), 0x02, byte(len(detail))}
	pdu = append(pdu, detail...)
	return pdu, nil
}

func (c *iecCodec) encodeConclude() ([]byte, error) {
	return []byte{0xA9, 0x02, 0x80, 0x00}, nil
}

func (c *iecCodec) encodeRead(name string) ([]byte, error) {
	// Read Request = 0xA4 + len + variableSpec
	invokeID := byte(c.invokeID)
	nameBytes := encodeVisibleString(name)
	// VariableSpec = 0x0C (listOfAccessResult) + len + name
	varSpec := []byte{0x0C, byte(1 + 1 + len(nameBytes)), 0x02, byte(len(nameBytes))}
	varSpec = append(varSpec, nameBytes...)
	pdu := []byte{0xA4, byte(len(varSpec)), invokeID}
	pdu = append(pdu, varSpec...)
	return pdu, nil
}

func (c *iecCodec) encodeWrite(req *kernel.Request) ([]byte, error) {
	// Write Request = 0xA5 + len + variableSpec + tag + value
	invokeID := byte(c.invokeID)
	nameBytes := encodeVisibleString(req.Address)
	varSpec := []byte{0x02, byte(len(nameBytes))}
	varSpec = append(varSpec, nameBytes...)

	value := encodeValue(req.Data)
	innerLen := byte(len(varSpec) + len(value))

	pdu := []byte{0xA5, 0, invokeID}
	pdu[1] = innerLen
	pdu = append(pdu, varSpec...)
	pdu = append(pdu, value...)
	return pdu, nil
}

func encodeVisibleString(s string) []byte {
	result := []byte{TagVisibleString, byte(len(s))}
	result = append(result, []byte(s)...)
	return result
}

func encodeValue(data []byte) []byte {
	if len(data) == 0 {
		return []byte{TagBool, 1, 0}
	}
	if len(data) == 1 {
		return []byte{TagInteger, 1, data[0]}
	}
	if len(data) <= 4 {
		result := []byte{TagUnsigned, byte(len(data))}
		return append(result, data...)
	}
	result := []byte{TagOctetString, byte(len(data))}
	return append(result, data...)
}

func decodeValue(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("iec61850: value too short")
	}
	tag := data[0]
	length := int(data[1])
	if len(data) < 2+length {
		return nil, fmt.Errorf("iec61850: value truncated")
	}
	switch tag {
	case TagBool, TagInteger, TagUnsigned, TagOctetString, TagVisibleString:
		return data[2 : 2+length], nil
	case TagFloat:
		if length >= 4 {
			buf := make([]byte, 4)
			binary.BigEndian.PutUint32(buf, binary.BigEndian.Uint32(data[2:6]))
			return buf, nil
		}
	}
	return data[2 : 2+length], nil
}
