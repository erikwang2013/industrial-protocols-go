package bridge

import (
	"bytes"
	"testing"
)

type mockSerial struct{ bytes.Buffer }

func (m *mockSerial) Close() error { return nil }

func TestSerialBridgeMock(t *testing.T) {
	b := NewSerialBridge("/dev/ttyUSB0", 115200)
	if b.running {
		t.Error("not running before Start()")
	}
	rw := &mockSerial{}
	b.SetPort(rw)
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	tr := b.Transport()
	if tr == nil {
		t.Fatal("transport is nil")
	}
	tr.Write([]byte("AT\r"))
	b.Stop()
}

func TestSerialBridgeInterface(t *testing.T) {
	var b Bridge = NewSerialBridge("/dev/ttyS0", 9600)
	if b == nil {
		t.Fatal("Bridge interface not satisfied")
	}
}
