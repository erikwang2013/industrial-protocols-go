// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package dnp3

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

const (
	TPDUStart1 = 0x05
	TPDUStart2 = 0x64

	FuncConfirm  = 0x00
	FuncRead     = 0x01
	FuncResponse = 0x81

	FIR = 0x80
	FIN = 0x40
)

type DNP3Protocol struct{}

func New() *DNP3Protocol                        { return &DNP3Protocol{} }
func (p *DNP3Protocol) Name() string            { return "dnp3" }
func (p *DNP3Protocol) Variants() []string      { return []string{"serial", "tcp"} }
func (p *DNP3Protocol) DefaultPort() int        { return 20000 }

func (p *DNP3Protocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "tcp" && v != "serial" {
		return nil, fmt.Errorf("dnp3: unsupported variant %q", v)
	}
	return &dnp3Codec{seq: 0}, nil
}

type dnp3Codec struct{ seq byte }

func (c *dnp3Codec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "read":
		return c.encodeAppLayer(FuncRead)
	case "class0_poll":
		return c.encodeClass0Poll()
	case "confirm":
		return c.encodeAppLayer(FuncConfirm)
	default:
		return c.encodeClass0Poll()
	}
}

func (c *dnp3Codec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("dnp3: frame too short (%d bytes)", len(data))
	}
	// Detect TCP vs serial (serial has CRC between blocks)
	if data[0] == TPDUStart1 && data[1] == TPDUStart2 {
		return c.decodeTPDU(data)
	}
	return c.decodeAPDU(data)
}

func (c *dnp3Codec) encodeAppLayer(funcCode byte) ([]byte, error) {
	c.seq++
	// APDU: Control(1) + Function(1)
	apdu := []byte{(c.seq & 0x3F) | FIR | FIN, funcCode}
	return c.wrapTPDU(apdu, 1, 2), nil
}

func (c *dnp3Codec) encodeClass0Poll() ([]byte, error) {
	c.seq++
	// Class 0 poll: no object header (reads all static objects)
	apdu := []byte{(c.seq & 0x3F) | FIR | FIN, FuncRead}
	return c.wrapTPDU(apdu, 1, 2), nil
}

func (c *dnp3Codec) wrapTPDU(apdu []byte, src, dst uint16) []byte {
	// TPDU: start(2) + length(1) + control(1) + dest(2) + src(2) + CRC(2) + ... + CRC(2)
	length := byte(len(apdu) + 5) // control + dest + src + apdu
	control := byte(c.seq & 0x3F) | FIR | FIN

	tpdu := []byte{TPDUStart1, TPDUStart2, length, control}
	addrBytes := make([]byte, 4)
	binary.LittleEndian.PutUint16(addrBytes[0:2], dst)
	binary.LittleEndian.PutUint16(addrBytes[2:4], src)
	tpdu = append(tpdu, addrBytes...)

	// CRC over header before APDU (control+dest+src = 5 bytes)
	hdrCRC := crc16DNP(tpdu[3:])
	tpdu = append(tpdu, byte(hdrCRC), byte(hdrCRC>>8))
	tpdu = append(tpdu, apdu...)

	// CRC over APDU data
	if len(apdu) > 0 {
		apduCRC := crc16DNP(apdu)
		tpdu = append(tpdu, byte(apduCRC), byte(apduCRC>>8))
	}

	return tpdu
}

func (c *dnp3Codec) decodeTPDU(data []byte) (*kernel.Response, error) {
	if data[0] != TPDUStart1 || data[1] != TPDUStart2 {
		return nil, fmt.Errorf("dnp3: invalid TPDU start bytes")
	}
	length := int(data[2])
	hdrCRC := binary.LittleEndian.Uint16(data[8:10])
	expectedCRC := crc16DNP(data[3:8])
	if hdrCRC != expectedCRC {
		return nil, fmt.Errorf("dnp3: TPDU header CRC mismatch")
	}

	apduStart := 10
	if apduStart+2 > len(data) {
		return nil, fmt.Errorf("dnp3: no APDU data")
	}

	apduEnd := apduStart + (length - 5) // subtract control+dest+src
	if apduEnd+2 <= len(data) {
		apduCRC := binary.LittleEndian.Uint16(data[apduEnd : apduEnd+2])
		expectedAPDUCRC := crc16DNP(data[apduStart:apduEnd])
		if apduCRC != expectedAPDUCRC {
			return nil, fmt.Errorf("dnp3: APDU CRC mismatch")
		}
	}

	return c.decodeAPDU(data[apduStart:apduEnd])
}

func (c *dnp3Codec) decodeAPDU(apdu []byte) (*kernel.Response, error) {
	if len(apdu) < 2 {
		return nil, fmt.Errorf("dnp3: APDU too short")
	}
	control := apdu[0]
	function := apdu[1]
	isFIR := control&FIR != 0
	isFIN := control&FIN != 0

	return &kernel.Response{
		Data: apdu,
		Metadata: map[string]any{
			"function": int(function),
			"fir":      isFIR,
			"fin":      isFIN,
			"sequence": int(control & 0x3F),
		},
	}, nil
}

func crc16DNP(data []byte) uint16 {
	var crc uint16 = 0x0000
	for _, b := range data {
		temp := crc ^ uint16(b)
		for i := 0; i < 8; i++ {
			if temp&1 != 0 {
				temp = (temp >> 1) ^ 0xA6BC
			} else {
				temp >>= 1
			}
		}
		crc = (crc >> 8) ^ temp
	}
	return ^crc
}
