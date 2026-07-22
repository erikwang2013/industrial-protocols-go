// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kline

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	DefaultBaud    = 10400
	ReqFormat      = 0x68
	RespFormat     = 0x48
	SIDCurrentData = 0x01
	SIDRequestDTCs = 0x03
	SIDVehicleInfo = 0x09
	SIDFastInit    = 0x33
	FastInitSync   = 0x55
	KeyByte1       = 0x08
	KeyByte2       = 0x08
)

type KLineProtocol struct{}

func New() *KLineProtocol { return &KLineProtocol{} }

func (p *KLineProtocol) Name() string       { return "kline" }
func (p *KLineProtocol) Variants() []string  { return []string{"serial"} }
func (p *KLineProtocol) DefaultPort() int    { return DefaultBaud }

func (p *KLineProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "serial" {
		return nil, fmt.Errorf("kline: unsupported variant %q", v)
	}
	return &kLineCodec{}, nil
}

type kLineCodec struct{}

func (c *kLineCodec) Encode(req *kernel.Request) ([]byte, error) {
	target := byte(0x33)
	source := byte(0xF1)
	if v, ok := req.Metadata["target"].(byte); ok {
		target = v
	}
	if v, ok := req.Metadata["source"].(byte); ok {
		source = v
	}

	var sid byte
	switch req.Function {
	case "current_data":
		sid = SIDCurrentData
	case "request_dtcs":
		sid = SIDRequestDTCs
	case "vehicle_info":
		sid = SIDVehicleInfo
	default:
		sid = SIDCurrentData
	}

	// Format + Target + Source + SID + Data + CS
	frame := []byte{ReqFormat, target, source, sid}
	frame = append(frame, req.Data...)
	cs := checksum(frame)
	frame = append(frame, cs)

	return frame, nil
}

func (c *kLineCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("kline: frame too short (%d bytes)", len(data))
	}

	if data[0] == FastInitSync {
		return &kernel.Response{Metadata: map[string]any{"fast_init_sync": true}}, nil
	}

	// Response: Format(0x48) + Source + Target + SID + Data[] + CS
	if len(data) < 2 {
		return nil, fmt.Errorf("kline: response too short")
	}

	// Verify checksum
	expectedCS := checksum(data[:len(data)-1])
	if expectedCS != data[len(data)-1] {
		return nil, fmt.Errorf("kline: checksum mismatch (got %02X, expected %02X)", data[len(data)-1], expectedCS)
	}

	sid := data[3]
	dataBytes := data[4 : len(data)-1]

	return &kernel.Response{
		Data: dataBytes,
		Metadata: map[string]any{
			"sid":    sid & 0x3F, // response SID = request SID + 0x40
			"source": data[1],
			"target": data[2],
		},
	}, nil
}

// FastInitSequence returns the bytes to write for a 5-baud fast init wakeup.
// Callers should send this byte at 5 baud (200ms per bit), then switch to 10400 baud.
func FastInitSequence() []byte {
	return []byte{SIDFastInit} // 0x33
}

// KeyBytes returns the key bytes sent after receiving the sync byte (0x55).
func KeyBytes() []byte {
	return []byte{KeyByte1, KeyByte2}
}

// VerifyKeyResponse verifies the ECU's response to the key byte exchange.
// The ECU should respond with the bitwise complement of each key byte.
func VerifyKeyResponse(resp []byte) bool {
	if len(resp) < 2 {
		return false
	}
	return resp[0] == byte(0xFF ^ KeyByte1) && resp[1] == byte(0xFF ^ KeyByte2)
}

func checksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return sum
}
