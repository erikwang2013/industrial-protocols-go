// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

//go:build linux

package bridge

import "testing"

func TestNewCANBridge(t *testing.T) {
	b := NewCANBridge("vcan0")
	if b.iface != "vcan0" {
		t.Errorf("expected vcan0, got %s", b.iface)
	}
	if b.running {
		t.Error("not running before Start()")
	}
}

func TestCANFrame(t *testing.T) {
	f := CANFrame{ID: 0x123, Data: []byte{1, 2, 3}, Ext: true}
	if f.ID != 0x123 {
		t.Errorf("expected 0x123, got 0x%X", f.ID)
	}
}

func TestCANBridge_StopBeforeStart(t *testing.T) {
	b := NewCANBridge("vcan0")
	if err := b.Stop(); err != nil {
		t.Errorf("expected nil error before Start, got %v", err)
	}
}

func TestCANBridge_TransportBeforeStart(t *testing.T) {
	b := NewCANBridge("vcan0")
	tr := b.Transport()
	if tr == nil {
		t.Fatal("Transport() returned nil")
	}
	if tr.Addr() != "vcan0" {
		t.Errorf("expected addr vcan0, got %s", tr.Addr())
	}
	if tr.Alive() {
		t.Error("transport should not be alive before Start")
	}
}
