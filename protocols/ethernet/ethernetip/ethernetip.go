// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package ethernetip

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	CmdRegisterSession   = 0x0065
	CmdUnRegisterSession = 0x0066
	CmdSendRRData        = 0x006F

	CIPServiceReadTag  = 0x4C
	CIPServiceWriteTag = 0x4D
)

type EthernetIPProtocol struct{}

func New() *EthernetIPProtocol { return &EthernetIPProtocol{} }

func (p *EthernetIPProtocol) Name() string       { return "ethernetip" }
func (p *EthernetIPProtocol) Variants() []string  { return []string{"tcp"} }
func (p *EthernetIPProtocol) DefaultPort() int    { return 44818 }

func (p *EthernetIPProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "tcp" {
		return nil, fmt.Errorf("ethernetip: unsupported variant %q", v)
	}
	return &enipCodec{invokeID: 0}, nil
}

type enipCodec struct{ invokeID uint32 }

func (c *enipCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "register_session":
		return c.encodeSession(CmdRegisterSession), nil
	case "unregister_session":
		session := uint32(0)
		if v, ok := req.Metadata["session"].(float64); ok {
			session = uint32(v)
		}
		return c.encodeSessionWithHandle(CmdUnRegisterSession, session), nil
	case "read_tag":
		session := uint32(0)
		if v, ok := req.Metadata["session"].(float64); ok {
			session = uint32(v)
		}
		tagName := req.Address
		return c.encodeReadTag(session, tagName)
	default:
		return nil, fmt.Errorf("ethernetip: unknown function %q", req.Function)
	}
}

func (c *enipCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 24 {
		return nil, fmt.Errorf("ethernetip: frame too short (%d bytes)", len(data))
	}
	cmd := binary.LittleEndian.Uint16(data[0:2])
	status := binary.LittleEndian.Uint32(data[8:12])
	session := binary.LittleEndian.Uint32(data[4:8])

	meta := map[string]any{"command": int(cmd), "status": status, "session": session}

	if cmd == CmdRegisterSession && status == 0 {
		return &kernel.Response{Metadata: meta}, nil
	}

	return &kernel.Response{Data: data[24:], Metadata: meta}, nil
}

func (c *enipCodec) encodeSession(cmd uint16) []byte {
	hdr := make([]byte, 24)
	binary.LittleEndian.PutUint16(hdr[0:2], cmd)
	binary.LittleEndian.PutUint16(hdr[2:4], 4) // length
	binary.LittleEndian.PutUint32(hdr[18:22], 1)
	payload := []byte{0x01, 0x00, 0x00, 0x00} // protocol version 1
	binary.LittleEndian.PutUint16(hdr[2:4], uint16(len(payload)))
	return append(hdr, payload...)
}

func (c *enipCodec) encodeSessionWithHandle(cmd uint16, session uint32) []byte {
	hdr := make([]byte, 24)
	binary.LittleEndian.PutUint16(hdr[0:2], cmd)
	binary.LittleEndian.PutUint16(hdr[2:4], 0)
	binary.LittleEndian.PutUint32(hdr[4:8], session)
	return hdr
}

func (c *enipCodec) encodeReadTag(session uint32, tagName string) ([]byte, error) {
	c.invokeID++
	// ENIP SendRRData header
	hdr := make([]byte, 24)
	binary.LittleEndian.PutUint16(hdr[0:2], CmdSendRRData)
	binary.LittleEndian.PutUint32(hdr[4:8], session)

	// CIP Read Tag request
	cipReq := make([]byte, 8+len(tagName))
	cipReq[0] = CIPServiceReadTag
	binary.LittleEndian.PutUint16(cipReq[2:4], uint16(c.invokeID))
	cipReq[4] = 0x02 // segments: 2
	cipReq[5] = 0x20 // symbolic segment (ANSI extended symbol)
	cipReq[6] = byte(len(tagName))
	copy(cipReq[7:], tagName)

	totalLen := 4 + len(cipReq) // 24 字节头部之后的 item 头(4) + CIP 请求
	binary.LittleEndian.PutUint16(hdr[2:4], uint16(totalLen))
	// Interface handle (4) + Timeout (2) + ItemCount (2) + ItemID (2) + ItemLength (2)
	binary.LittleEndian.PutUint32(hdr[12:16], 0) // interface handle
	binary.LittleEndian.PutUint16(hdr[16:18], 0) // timeout

	payload := make([]byte, 28+len(cipReq))
	copy(payload, hdr)
	binary.LittleEndian.PutUint32(payload[12:16], 0)     // interface handle
	binary.LittleEndian.PutUint16(payload[20:22], 0)      // timeout
	binary.LittleEndian.PutUint16(payload[22:24], 1)      // 1 item
	binary.LittleEndian.PutUint16(payload[24:26], 0x00B2) // unconnected data item
	binary.LittleEndian.PutUint16(payload[26:28], uint16(len(cipReq)))
	copy(payload[28:], cipReq)

	return payload, nil
}
