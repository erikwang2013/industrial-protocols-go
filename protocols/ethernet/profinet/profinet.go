// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package profinet

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	FrameIDDCPIdentify = 0xFEFD
	FrameIDDCPSet      = 0xFEFE
	FrameIDRecordData  = 0xFEFC

	ServiceIDIdentify = 5
	ServiceIDSet      = 4

	ServiceTypeRequest  = 0
	ServiceTypeResponse = 1
)

type ProfinetProtocol struct{}

func New() *ProfinetProtocol { return &ProfinetProtocol{} }

func (p *ProfinetProtocol) Name() string       { return "profinet" }
func (p *ProfinetProtocol) Variants() []string  { return []string{"nrt"} }
func (p *ProfinetProtocol) DefaultPort() int    { return 34964 }

func (p *ProfinetProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "nrt" {
		return nil, fmt.Errorf("profinet: unsupported variant %q", v)
	}
	return &pnCodec{xid: 1}, nil
}

type pnCodec struct{ xid uint32 }

func (c *pnCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "dcp_identify":
		return c.encodeDCP(FrameIDDCPIdentify, ServiceIDIdentify, nil)
	case "dcp_set":
		name := ""
		ip := ""
		if v, ok := req.Metadata["device_name"].(string); ok {
			name = v
		}
		if v, ok := req.Metadata["ip_address"].(string); ok {
			ip = v
		}
		return c.encodeDCPSet(name, ip)
	case "record_read":
		return c.encodeRecordRead(req)
	default:
		return nil, fmt.Errorf("profinet: unknown function %q", req.Function)
	}
}

func (c *pnCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("profinet: frame too short (%d bytes)", len(data))
	}
	frameID := binary.BigEndian.Uint16(data[0:2])
	serviceID := data[2]
	serviceType := data[3]
	xid := binary.BigEndian.Uint32(data[4:8])
	blockLen := binary.BigEndian.Uint16(data[8:10])
	if len(data) < 10+int(blockLen) {
		return nil, fmt.Errorf("profinet: block length %d exceeds frame length %d", blockLen, len(data))
	}
	blocks := data[10 : 10+blockLen]

	return &kernel.Response{
		Data: blocks,
		Metadata: map[string]any{
			"frame_id":     int(frameID),
			"service_id":   int(serviceID),
			"service_type": int(serviceType),
			"xid":          xid,
		},
	}, nil
}

func (c *pnCodec) encodeDCP(frameID uint16, serviceID byte, blocks []byte) ([]byte, error) {
	c.xid++
	hdr := make([]byte, 10)
	binary.BigEndian.PutUint16(hdr[0:2], frameID)
	hdr[2] = serviceID
	hdr[3] = ServiceTypeRequest
	binary.BigEndian.PutUint32(hdr[4:8], c.xid)
	binary.BigEndian.PutUint16(hdr[8:10], uint16(len(blocks)))
	return append(hdr, blocks...), nil
}

func (c *pnCodec) encodeDCPSet(name, ip string) ([]byte, error) {
	// DCP Set with NameOfStation and IPParameter blocks
	blocks := append(buildNameBlock(name), buildIPBlock(ip)...)
	return c.encodeDCP(FrameIDDCPSet, ServiceIDSet, blocks)
}

func (c *pnCodec) encodeRecordRead(req *kernel.Request) ([]byte, error) {
	// Simplified RecordDataRead
	index := uint32(0)
	if v, ok := req.Metadata["index"].(float64); ok {
		index = uint32(v)
	}
	hdr := make([]byte, 10)
	binary.BigEndian.PutUint16(hdr[0:2], FrameIDRecordData)
	hdr[2] = 0x01 // RecordDataRead
	hdr[3] = ServiceTypeRequest
	c.xid++
	binary.BigEndian.PutUint32(hdr[4:8], c.xid)
	payload := make([]byte, 6)
	binary.BigEndian.PutUint32(payload[2:6], index)
	binary.BigEndian.PutUint16(hdr[8:10], uint16(len(payload)))
	return append(hdr, payload...), nil
}

func buildNameBlock(name string) []byte {
	if name == "" {
		name = "device"
	}
	opts := []byte{0x02, 0x00, 0x02, byte(len(name))}
	return append(opts, []byte(name)...)
}

func buildIPBlock(ip string) []byte {
	block := []byte{0x01, 0x00, 0x0E, 0x0C}
	// Simple IP parsing: "192.168.1.10"
	if ip != "" {
		var a, b, c, d byte
		fmt.Sscanf(ip, "%d.%d.%d.%d", &a, &b, &c, &d)
		block = append(block, a, b, c, d)
		block = append(block, 0xFF, 0xFF, 0xFF, 0x00) // subnet 255.255.255.0
		block = append(block, a, b, c, 1)               // gateway
	} else {
		block = append(block, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	}
	return block
}
