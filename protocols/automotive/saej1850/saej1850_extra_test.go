// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package saej1850

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// newCodec returns the concrete j1850Codec so tests can exercise
// unexported fields and ID-building helpers.
func newCodec(t *testing.T) *j1850Codec {
	t.Helper()
	c, err := New().NewCodec("can")
	if err != nil {
		t.Fatal(err)
	}
	return c.(*j1850Codec)
}

func TestBuildIDs(t *testing.T) {
	c := newCodec(t)
	// Priority 6 << 26 | ext 1 << 25 | PF | PS | SA.
	if got := c.buildPhysicalReqID(); got != 0x1ADA01F1 {
		t.Errorf("physical request ID = 0x%08X, want 0x1ADA01F1", got)
	}
	if got := c.buildPhysicalRespID(); got != 0x1ADAF101 {
		t.Errorf("physical response ID = 0x%08X, want 0x1ADAF101", got)
	}
	if got := c.buildFunctionalID(); got != 0x1ADB33F1 {
		t.Errorf("functional ID = 0x%08X, want 0x1ADB33F1", got)
	}
}

func TestBuildIDsCustomAddresses(t *testing.T) {
	c := newCodec(t)
	c.targetAddr = 0x10
	c.sourceAddr = 0x02
	if got := c.buildPhysicalReqID(); got != 0x1ADA1002 {
		t.Errorf("physical request ID = 0x%08X, want 0x1ADA1002", got)
	}
	if got := c.buildPhysicalRespID(); got != 0x1ADA0210 {
		t.Errorf("physical response ID = 0x%08X, want 0x1ADA0210", got)
	}
}

func TestBuildIDsCustomPriority(t *testing.T) {
	c := newCodec(t)
	c.priority = 2
	if got := c.buildFunctionalID(); got != 0x0ADB33F1 {
		t.Errorf("functional ID with priority 2 = 0x%08X, want 0x0ADB33F1", got)
	}
}

func TestBuildDiagPayload(t *testing.T) {
	c := newCodec(t)
	payload := c.buildDiagPayload(&kernel.Request{
		Metadata: map[string]any{"mode": float64(0x22), "pid": float64(0x0C)},
	})
	want := []byte{0x02, 0x22, 0x0C, 0, 0, 0, 0, 0}
	for i := range want {
		if payload[i] != want[i] {
			t.Fatalf("payload = %v, want %v", payload, want)
		}
	}
}

func TestEncodeMode01PID(t *testing.T) {
	c := newCodec(t)
	raw, err := c.Encode(&kernel.Request{
		Function: "mode01",
		Metadata: map[string]any{"pid": float64(0x0C)},
	})
	if err != nil {
		t.Fatal(err)
	}
	// CAN frame: ID(4) + data(8); data[0]=0x02 len, data[1]=mode, data[2]=PID.
	// The wire format carries the first 4 data bytes at raw[4:8].
	if raw[4] != 0x02 || raw[5] != 0x01 || raw[6] != 0x0C {
		t.Errorf("mode01 frame data = % X, want 02 01 0C 00", raw[4:])
	}
}

func TestDecodeMetadataTypes(t *testing.T) {
	c := newCodec(t)
	// Functional request frame: 0x1ADB33F1.
	raw, _ := c.Encode(&kernel.Request{Function: "broadcast", Data: []byte{1}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "functional_request" {
		t.Errorf("type = %v, want functional_request", resp.Metadata["type"])
	}
	if resp.Metadata["pf"] != byte(0xDB) {
		t.Errorf("pf = %v, want 0xDB", resp.Metadata["pf"])
	}
	if resp.Metadata["priority"] != byte(6) {
		t.Errorf("priority = %v, want 6", resp.Metadata["priority"])
	}
}

func TestDecodePhysicalResponseMetadata(t *testing.T) {
	c := newCodec(t)
	raw, _ := c.Encode(&kernel.Request{Function: "diag_response", Data: []byte{0x41}})
	resp, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Metadata["type"] != "physical_response" {
		t.Errorf("type = %v, want physical_response", resp.Metadata["type"])
	}
}
