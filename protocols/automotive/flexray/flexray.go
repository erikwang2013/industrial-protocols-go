package flexray

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// FlexRayProtocol implements kernel.Protocol for FlexRay over CAN.
type FlexRayProtocol struct{}

// New creates a new FlexRayProtocol.
func New() *FlexRayProtocol { return &FlexRayProtocol{} }

// Name returns the protocol identifier.
func (p *FlexRayProtocol) Name() string { return "flexray" }

// Variants returns supported transport variants.
func (p *FlexRayProtocol) Variants() []string { return []string{"can", "serial"} }

// DefaultPort returns the default port (0 for CAN bus).
func (p *FlexRayProtocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant.
func (p *FlexRayProtocol) NewCodec(v string) (kernel.Codec, error) {
	switch v {
	case "can", "serial":
		return &flexrayCodec{slotID: 1, baseCycle: 0}, nil
	default:
		return nil, fmt.Errorf("flexray: unsupported variant %q", v)
	}
}

// flexrayCodec encodes/decodes FlexRay frames over CAN bus.
//
// FlexRay cycle payload format:
//
//	Header (2 bytes): cycle number (little-endian)
//	Status (1 byte):  bit 7=PPI, bit 6=NFI, bit 5=SYF, bit 4=SUF
//	Data (N bytes):   payload (max 254 bytes)
//	CRC (2 bytes):    CRC-16 over header+status+data (little-endian)
//
// Each FlexRay frame maps to a CAN frame with 29-bit extended ID.
type flexrayCodec struct {
	slotID    uint16
	baseCycle uint16
}

// Encode serializes a FlexRay request into a CAN frame.
// The CAN ID encodes slot ID and cycle number.
func (c *flexrayCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "frame":
		cycle := c.baseCycle
		if v, ok := req.Metadata["cycle"].(float64); ok {
			cycle = uint16(v)
		}
		ppi := byte(0)
		if v, ok := req.Metadata["ppi"].(bool); ok && v {
			ppi = 1
		}
		payload := req.Data
		if len(payload) > 254 {
			payload = payload[:254]
		}

		// Build FlexRay cycle payload
		frPayload := make([]byte, 2+1+len(payload)+2)
		binary.LittleEndian.PutUint16(frPayload[0:2], cycle)
		frPayload[2] = ppi << 7 // PPI flag in bit 7 of status byte
		copy(frPayload[3:], payload)
		crc := crc16(frPayload[:3+len(payload)])
		binary.LittleEndian.PutUint16(frPayload[3+len(payload):], crc)

		// CAN ID: 29-bit with slot and cycle
		canID := uint32(0x01000000) | (uint32(c.slotID) << 10) | uint32(cycle)

		return marshalCAN(bridge.CANFrame{
			ID:   canID,
			Data: frPayload[:min(len(frPayload), 4)], // First 4 bytes in CAN frame
			Ext:  true,
		}), nil

	case "status":
		// Status request returns current slot configuration
		canID := uint32(0x02000000) | uint32(c.slotID)
		return marshalCAN(bridge.CANFrame{
			ID:   canID,
			Data: []byte{byte(c.slotID), byte(c.slotID >> 8), byte(c.baseCycle), byte(c.baseCycle >> 8)},
			Ext:  true,
		}), nil

	default:
		return nil, fmt.Errorf("flexray: unknown function %q", req.Function)
	}
}

// Decode deserializes a CAN frame into a kernel.Response.
func (c *flexrayCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("flexray: frame too short")
	}
	f := unmarshalCAN(data)
	meta := map[string]any{
		"can_id":  f.ID,
		"ext":     f.Ext,
		"slot_id": uint16((f.ID >> 10) & 0x3F),
		"cycle":   uint16(f.ID & 0x3F),
	}
	return &kernel.Response{Data: f.Data, Metadata: meta}, nil
}

// marshalCAN serializes a CANFrame into an 8-byte wire format.
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

// crc16 computes CRC-16/XMODEM over the given data.
func crc16(data []byte) uint16 {
	var crc uint16 = 0x0000
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}
