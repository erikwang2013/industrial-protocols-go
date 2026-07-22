// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package profinet

import (
	"encoding/binary"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "profinet" {
		t.Errorf("expected profinet, got %s", p.Name())
	}
	if p.DefaultPort() != 34964 {
		t.Errorf("expected 34964, got %d", p.DefaultPort())
	}
}

func TestEncodeDCPIdentify(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	data, err := c.Encode(&kernel.Request{Function: "dcp_identify"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 10 {
		t.Fatalf("expected 10 bytes, got %d", len(data))
	}
	fid := int(data[0])<<8 | int(data[1])
	if fid != FrameIDDCPIdentify {
		t.Errorf("expected 0xFEFD, got 0x%04X", fid)
	}
}

func TestDecodeDCPResponse(t *testing.T) {
	c, _ := New().NewCodec("nrt")
	resp := make([]byte, 12)
	binary.BigEndian.PutUint16(resp[0:2], FrameIDDCPIdentify)
	resp[2] = ServiceIDIdentify
	resp[3] = ServiceTypeResponse
	binary.BigEndian.PutUint32(resp[4:8], 1)
	binary.BigEndian.PutUint16(resp[8:10], 2)
	resp[10] = 0x01
	resp[11] = 0x02

	result, err := c.Decode(resp)
	if err != nil {
		t.Fatal(err)
	}
	if result.Metadata["service_type"].(int) != ServiceTypeResponse {
		t.Error("expected response type")
	}
}
