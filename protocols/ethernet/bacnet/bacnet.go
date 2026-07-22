// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package bacnet

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	BVLCTypeOrigUnicast   = 0x81
	BVLCTypeOrigBroadcast = 0x82
	BVLCFunctionLength    = 0x0B // 11 bytes fixed (function 0x0B = original message length)

	APDUTypeConfirmedReq   = 0x00
	APDUTypeUnconfirmedReq = 0x10 // actually type 8 (0x10 << 0 = 0x10)

	ServiceWhoIs        = 3
	ServiceIAm          = 0
	ServiceReadProperty = 12
)

type BacnetProtocol struct{}

func New() *BacnetProtocol                        { return &BacnetProtocol{} }
func (p *BacnetProtocol) Name() string            { return "bacnet" }
func (p *BacnetProtocol) Variants() []string      { return []string{"ip"} }
func (p *BacnetProtocol) DefaultPort() int        { return 47808 }

func (p *BacnetProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "ip" {
		return nil, fmt.Errorf("bacnet: unsupported variant %q", v)
	}
	return &bacnetCodec{}, nil
}

type bacnetCodec struct{}

func (c *bacnetCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "who_is":
		return c.encodeWhoIs()
	case "read_property":
		return c.encodeReadProperty(req)
	default:
		return c.encodeWhoIs()
	}
}

func (c *bacnetCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("bacnet: frame too short (%d bytes)", len(data))
	}

	bvlcType := data[0]
	npduStart := 4
	if bvlcType == BVLCTypeOrigUnicast || bvlcType == BVLCTypeOrigBroadcast {
		npduStart = 4
	} else {
		return nil, fmt.Errorf("bacnet: unsupported BVLC type 0x%02X", bvlcType)
	}

	npdu := data[npduStart:]
	if len(npdu) < 2 {
		return nil, fmt.Errorf("bacnet: NPDU too short")
	}

	apduStart := npduStart + 2
	apdu := data[apduStart:]

	metadata := map[string]any{"bvlc_type": bvlcType}

	if len(apdu) >= 2 {
		apduType := apdu[0] >> 4
		service := apdu[1]
		metadata["apdu_type"] = int(apduType)
		metadata["service"] = int(service)
	}

	return &kernel.Response{Data: data, Metadata: metadata}, nil
}

func (c *bacnetCodec) encodeWhoIs() ([]byte, error) {
	// BVLC: type(1) + function(1) + length(2)
	// NPDU: version(1) + control(1)
	// APDU: type+flags(1) = 0x10 (unconfirmed), service(1) = 3 (Who-Is)
	buf := make([]byte, 8)
	buf[0] = BVLCTypeOrigBroadcast
	buf[1] = BVLCFunctionLength
	binary.BigEndian.PutUint16(buf[2:4], 4) // length of NPDU+APDU
	buf[4] = 0x01                            // NPDU version
	buf[5] = 0x00                            // NPDU control
	buf[6] = 0x10                            // APDU: unconfirmed request
	buf[7] = ServiceWhoIs
	return buf, nil
}

func (c *bacnetCodec) encodeReadProperty(req *kernel.Request) ([]byte, error) {
	objType := byte(0)
	objInstance := uint32(0)
	propID := uint32(8) // Present Value

	if v, ok := req.Metadata["object_type"].(byte); ok {
		objType = v
	}
	if v, ok := req.Metadata["object_instance"].(float64); ok {
		objInstance = uint32(v)
	}
	if v, ok := req.Metadata["property_id"].(float64); ok {
		propID = uint32(v)
	}

	// BVLC(4) + NPDU(2) + APDU: type(1) + maxSegments(1) + invokeID(1) + service(1)
	// + objectID(4) + propID(4)
	buf := make([]byte, 18)
	buf[0] = BVLCTypeOrigUnicast
	buf[1] = BVLCFunctionLength
	binary.BigEndian.PutUint16(buf[2:4], 14)
	buf[4] = 0x01 // NPDU version
	buf[5] = 0x00
	buf[6] = APDUTypeConfirmedReq // confirmed request
	buf[7] = 0x00                 // max segments + max ADPU
	buf[8] = 1                     // invoke ID
	buf[9] = ServiceReadProperty
	// Object ID: 10 bits type + 22 bits instance
	objID := (uint32(objType)&0x3FF)<<22 | (objInstance & 0x3FFFFF)
	binary.BigEndian.PutUint32(buf[10:14], objID)
	binary.BigEndian.PutUint32(buf[14:18], propID)
	return buf, nil
}
