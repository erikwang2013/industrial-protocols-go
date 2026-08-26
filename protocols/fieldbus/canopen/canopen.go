// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package canopen

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// CanOpenProtocol implements kernel.Protocol for CANopen over CAN bus.
type CanOpenProtocol struct{}

// New creates a new CanOpenProtocol.
func New() *CanOpenProtocol { return &CanOpenProtocol{} }

// Name returns the protocol identifier.
func (p *CanOpenProtocol) Name() string { return "canopen" }

// Variants returns supported transport variants.
func (p *CanOpenProtocol) Variants() []string { return []string{"can"} }

// DefaultPort returns the default port (0 for CAN bus).
func (p *CanOpenProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant. Only "can" is supported.
// The codec uses Node-ID 1 by default; use driver.NewSocketCANDriver or set manually.
func (p *CanOpenProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "can" {
		return nil, fmt.Errorf("canopen: unsupported variant %q", v)
	}
	return &canopenCodec{nodeID: 1}, nil
}

// canopenCodec encodes/decodes CANopen frames to/from CAN bus wire format.
// CAN ID layout (standard 11-bit):
//
//	NMT     = 0x000
//	SYNC    = 0x080
//	SDO tx  = 0x580 + NodeID
//	SDO rx  = 0x600 + NodeID
//	PDO1 tx = 0x180 + NodeID
//	Heartbeat = 0x700 + NodeID
type canopenCodec struct {
	nodeID byte
}

// Encode serializes a kernel.Request into a CAN frame.
func (c *canopenCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "sdo_read":
		index, sub := c.getSDOAddr(req)
		f := bridge.CANFrame{
			ID:   0x600 + uint32(c.nodeID),
			Data: []byte{0x40, byte(index), byte(index >> 8), sub, 0, 0, 0, 0},
		}
		return marshalCAN(f), nil

	case "sdo_write":
		index, sub := c.getSDOAddr(req)
		data := req.Data
		if len(data) > 4 {
			data = data[:4]
		}
		f := bridge.CANFrame{
			ID:   0x600 + uint32(c.nodeID),
			Data: []byte{0x23, byte(index), byte(index >> 8), sub, 0, 0, 0, 0},
		}
		copy(f.Data[4:], data)
		return marshalCAN(f), nil

	case "nmt_start":
		return marshalCAN(bridge.CANFrame{
			ID: 0x000, Data: []byte{0x01, c.nodeID},
		}), nil

	case "nmt_stop":
		return marshalCAN(bridge.CANFrame{
			ID: 0x000, Data: []byte{0x02, c.nodeID},
		}), nil

	case "nmt_reset":
		return marshalCAN(bridge.CANFrame{
			ID: 0x000, Data: []byte{0x82, c.nodeID},
		}), nil

	case "heartbeat":
		return marshalCAN(bridge.CANFrame{
			ID: 0x700 + uint32(c.nodeID), Data: []byte{0x05},
		}), nil

	default:
		return nil, fmt.Errorf("canopen: unknown function %q", req.Function)
	}
}

// Decode deserializes a CAN frame into a kernel.Response.
func (c *canopenCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 8 { // 线格式固定 8 字节：4 字节 ID + 4 字节数据
		return nil, fmt.Errorf("canopen: frame too short (%d bytes)", len(data))
	}
	f := unmarshalCAN(data)
	meta := map[string]any{
		"can_id": f.ID,
		"ext":    f.Ext,
	}
	if f.ID == 0x700+uint32(c.nodeID) {
		meta["type"] = "bootup"
	}
	if len(f.Data) > 0 && f.Data[0] == 0x80 {
		code := fmt.Sprintf("%X", f.Data)
		if len(f.Data) > 4 {
			code = fmt.Sprintf("%X", f.Data[4:])
		}
		return nil, &kernel.ProtocolError{
			Code:    code,
			Message: "SDO abort",
			Raw:     f.Data,
		}
	}
	return &kernel.Response{Data: f.Data, Metadata: meta}, nil
}

// getSDOAddr extracts the object dictionary index and sub-index from request metadata.
func (c *canopenCodec) getSDOAddr(req *kernel.Request) (uint16, byte) {
	var index uint16 = 0x1000
	var sub byte
	if v, ok := req.Metadata["index"].(float64); ok {
		index = uint16(v)
	}
	if v, ok := req.Metadata["sub"].(float64); ok {
		sub = byte(v)
	}
	return index, sub
}

// marshalCAN serializes a CANFrame into an 8-byte wire format:
// bytes 0-3: little-endian ID (bit 7 of byte 3 set if extended frame)
// bytes 4-7: data (up to 4 bytes for now)
func marshalCAN(f bridge.CANFrame) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], f.ID)
	if f.Ext {
		buf[3] |= 0x80
	}
	copy(buf[4:], f.Data)
	return buf
}

// unmarshalCAN deserializes an 8-byte wire format into a CANFrame.
func unmarshalCAN(data []byte) bridge.CANFrame {
	id := binary.LittleEndian.Uint32(data[0:4])
	ext := data[3]&0x80 != 0
	id &= 0x7FFFFFFF
	return bridge.CANFrame{ID: id, Data: data[4:8], Ext: ext}
}
