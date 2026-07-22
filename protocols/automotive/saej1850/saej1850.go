package saej1850

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// SaeJ1850Protocol implements kernel.Protocol for SAE J1850 OBD-II over CAN bus.
type SaeJ1850Protocol struct{}

// New creates a new SaeJ1850Protocol.
func New() *SaeJ1850Protocol { return &SaeJ1850Protocol{} }

// Name returns the protocol identifier.
func (p *SaeJ1850Protocol) Name() string { return "saej1850" }

// Variants returns supported transport variants.
func (p *SaeJ1850Protocol) Variants() []string { return []string{"can"} }

// DefaultPort returns the default port (0 for CAN bus).
func (p *SaeJ1850Protocol) DefaultPort() int { return 0 }

// NewCodec creates a codec for the given variant. Only "can" is supported.
func (p *SaeJ1850Protocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "can" {
		return nil, fmt.Errorf("saej1850: unsupported variant %q", v)
	}
	return &j1850Codec{
		targetAddr: 0x01, // Default: ECM (Engine Control Module)
		sourceAddr: 0xF1, // Default: scan tool
	}, nil
}

// j1850Codec implements kernel.Codec for SAE J1850 OBD-II over CAN.
//
// CAN ID format (29-bit):
//
//	Bits 28-26: Priority (0-7, default 6)
//	Bit  25:    Extended ID = 1
//	Bits 24-16: PF (Parameter Format / Header)
//	Bits 15-8:  PS (Parameter Specific / target or source)
//	Bits 7-0:   SA (Source Address)
//
// OBD-II functional request:  0x18DB33F1 (broadcast)
// OBD-II physical request:    0x18DAxxF1 (xx = target ECU address)
// OBD-II physical response:   0x18DAF1xx (xx = source ECU address)
type j1850Codec struct {
	targetAddr byte
	sourceAddr byte
	priority   byte
}

// Encode serializes a kernel.Request into a CAN frame for J1850 OBD-II.
func (c *j1850Codec) Encode(req *kernel.Request) ([]byte, error) {
	pri := c.priority
	if pri == 0 {
		pri = 6 // Default priority: normal message
	}

	switch req.Function {
	case "diag_request":
		// Physical diagnostic request to specific ECU
		canID := c.buildPhysicalReqID()
		data := c.buildDiagPayload(req)
		return c.marshalFrame(canID, data, true), nil

	case "diag_response":
		// Physical diagnostic response from ECU
		canID := c.buildPhysicalRespID()
		data := req.Data
		if len(data) > 8 {
			data = data[:8]
		}
		return c.marshalFrame(canID, data, true), nil

	case "mode01":
		// Mode $01: Request current powertrain diagnostic data
		pid := byte(0x00) // PID $00 = supported PIDs
		if v, ok := req.Metadata["pid"].(float64); ok {
			pid = byte(v)
		}
		canID := c.buildPhysicalReqID()
		return c.marshalFrame(canID, []byte{0x02, 0x01, pid, 0x00, 0x00, 0x00, 0x00, 0x00}, true), nil

	case "mode03":
		// Mode $03: Request emission-related DTCs
		canID := c.buildPhysicalReqID()
		return c.marshalFrame(canID, []byte{0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, true), nil

	case "mode0A":
		// Mode $0A: Request permanent DTCs
		canID := c.buildPhysicalReqID()
		return c.marshalFrame(canID, []byte{0x02, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, true), nil

	case "broadcast":
		// Functional broadcast (all ECUs)
		canID := c.buildFunctionalID()
		data := req.Data
		if len(data) > 8 {
			data = data[:8]
		}
		return c.marshalFrame(canID, data, true), nil

	default:
		return nil, fmt.Errorf("saej1850: unknown function %q", req.Function)
	}
}

// Decode deserializes a CAN frame into a kernel.Response.
func (c *j1850Codec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("saej1850: frame too short")
	}
	f := c.unmarshalFrame(data)
	id := f.ID

	// Extract J1850 fields from 29-bit ID
	priority := byte((id >> 26) & 0x07)
	pf := byte((id >> 16) & 0xFF)
	ps := byte((id >> 8) & 0xFF)
	sa := byte(id & 0xFF)

	meta := map[string]any{
		"can_id":   f.ID,
		"ext":      f.Ext,
		"priority": priority,
		"pf":       pf,
		"ps":       ps,
		"sa":       sa,
	}

	// Determine message type from ID
	if pf == 0xDA && sa == 0xF1 {
		meta["type"] = "physical_request"
	} else if pf == 0xDA && ps == 0xF1 {
		meta["type"] = "physical_response"
	} else if pf == 0xDB {
		meta["type"] = "functional_request"
	}

	return &kernel.Response{Data: f.Data, Metadata: meta}, nil
}

// buildPhysicalReqID builds a 29-bit CAN ID for physical request (0x18DAxxF1).
func (c *j1850Codec) buildPhysicalReqID() uint32 {
	pri := uint32(c.priority)
	if pri == 0 {
		pri = 6
	}
	return (pri << 26) | (1 << 25) | (0xDA << 16) | (uint32(c.targetAddr) << 8) | uint32(c.sourceAddr)
}

// buildPhysicalRespID builds a 29-bit CAN ID for physical response (0x18DAF1xx).
func (c *j1850Codec) buildPhysicalRespID() uint32 {
	pri := uint32(c.priority)
	if pri == 0 {
		pri = 6
	}
	return (pri << 26) | (1 << 25) | (0xDA << 16) | (uint32(c.sourceAddr) << 8) | uint32(c.targetAddr)
}

// buildFunctionalID builds a 29-bit CAN ID for functional broadcast (0x18DB33F1).
func (c *j1850Codec) buildFunctionalID() uint32 {
	pri := uint32(c.priority)
	if pri == 0 {
		pri = 6
	}
	return (pri << 26) | (1 << 25) | (0xDB << 16) | (0x33 << 8) | uint32(c.sourceAddr)
}

// buildDiagPayload builds the ISO 15765-2 diagnostic payload.
func (c *j1850Codec) buildDiagPayload(req *kernel.Request) []byte {
	mode := byte(0x01)
	if v, ok := req.Metadata["mode"].(float64); ok {
		mode = byte(v)
	}
	pid := byte(0x00)
	if v, ok := req.Metadata["pid"].(float64); ok {
		pid = byte(v)
	}
	payload := make([]byte, 8)
	// Single frame: length in high nibble of byte 0
	payload[0] = 0x02 // 2 bytes of diagnostic data
	payload[1] = mode
	payload[2] = pid
	return payload
}

// marshalFrame serializes a CANFrame into an 8-byte wire format.
func (c *j1850Codec) marshalFrame(id uint32, data []byte, ext bool) []byte {
	buf := make([]byte, 8)
	buf[0] = byte(id)
	buf[1] = byte(id >> 8)
	buf[2] = byte(id >> 16)
	buf[3] = byte(id >> 24)
	if ext {
		buf[3] |= 0x80
	}
	copy(buf[4:], data)
	return buf
}

// unmarshalFrame deserializes an 8-byte wire format into a CANFrame.
func (c *j1850Codec) unmarshalFrame(data []byte) bridge.CANFrame {
	id := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	ext := data[3]&0x80 != 0
	id &= 0x7FFFFFFF
	return bridge.CANFrame{ID: id, Data: data[4:8], Ext: ext}
}
