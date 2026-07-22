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
