package devicenet

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// DeviceNetProtocol implements kernel.Protocol for DeviceNet over CAN bus.
type DeviceNetProtocol struct{}

// New creates a new DeviceNetProtocol.
func New() *DeviceNetProtocol { return &DeviceNetProtocol{} }

// Name returns the protocol identifier.
func (p *DeviceNetProtocol) Name() string { return "devicenet" }

// Variants returns supported transport variants.
func (p *DeviceNetProtocol) Variants() []string { return []string{"can", "gateway"} }

// DefaultPort returns the default TCP port for gateway connections.
func (p *DeviceNetProtocol) DefaultPort() int { return 2001 }

// NewCodec creates a codec for the given variant.
func (p *DeviceNetProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "can":
		return &dnCodec{macID: 0}, nil
	case "gateway":
		return &dnGatewayCodec{}, nil
	default:
		return nil, fmt.Errorf("devicenet: unsupported variant %q", v)
	}
}

// dnCodec implements kernel.Codec for DeviceNet via CAN bus.
// Uses Group 2 messages (0x400 + NodeID) for predefined master/slave connections.
type dnCodec struct {
	macID byte
}

// Encode serializes a kernel.Request into a CAN frame or wire data.
func (c *dnCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "open":
		// Explicit connection open: Group 3, Message ID 6
		id := uint32(0x80000000 | (uint32(c.macID) << 16))
		return marshalCAN(bridge.CANFrame{
			ID:   id,
			Data: []byte{0x4B, 0x03, 0x01, 0x01, 0x01, 0x00},
			Ext:  true,
		}), nil

	case "poll":
		// Group 2 poll request (master to slave)
		data := req.Data
		if len(data) > 8 {
			data = data[:8]
		}
		// Pad to 8 bytes if shorter
		padded := make([]byte, 8)
		copy(padded, data)
		id := uint32(0x400 | uint32(c.macID))
		return marshalCAN(bridge.CANFrame{ID: id, Data: padded}), nil

	default:
		return nil, fmt.Errorf("devicenet: unknown function %q", req.Function)
	}
}

// Decode deserializes wire data into a kernel.Response.
func (c *dnCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("devicenet: frame too short")
	}
	f := unmarshalCAN(data)
	return &kernel.Response{
		Data:     f.Data,
		Metadata: map[string]any{"can_id": f.ID, "ext": f.Ext},
	}, nil
}

// dnGatewayCodec implements kernel.Codec for DeviceNet via TCP gateway.
// Commands are sent as plain-text: "function hexdata\n".
type dnGatewayCodec struct{}

// Encode serializes a request as a plain-text gateway command.
func (c *dnGatewayCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := fmt.Sprintf("%s %x\n", req.Function, req.Data)
	return []byte(cmd), nil
}

// Decode deserializes the gateway response.
func (c *dnGatewayCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}

// marshalCAN serializes a CANFrame into an 8-byte wire format:
// bytes 0-3: little-endian ID (bit 7 of byte 3 set if extended frame)
// bytes 4-7: data
func marshalCAN(f bridge.CANFrame) []byte {
	buf := make([]byte, 8)
	buf[0] = byte(f.ID)
	buf[1] = byte(f.ID >> 8)
	buf[2] = byte(f.ID >> 16)
	buf[3] = byte(f.ID >> 24)
	if f.Ext {
		buf[3] |= 0x80
	}
	copy(buf[4:], f.Data)
	return buf
}

// unmarshalCAN deserializes an 8-byte wire format into a CANFrame.
func unmarshalCAN(data []byte) bridge.CANFrame {
	id := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	ext := data[3]&0x80 != 0
	id &= 0x7FFFFFFF
	return bridge.CANFrame{ID: id, Data: data[4:8], Ext: ext}
}
