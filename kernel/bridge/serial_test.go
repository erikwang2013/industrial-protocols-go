// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package bridge

import (
	"bytes"
	"net"
	"testing"
	"time"
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

func TestSerialBridge_TransportNilWithoutPort(t *testing.T) {
	b := NewSerialBridge("/dev/ttyS0", 9600)
	if tr := b.Transport(); tr != nil {
		t.Error("expected nil transport without a port")
	}
}

func TestSerialBridge_StopWithoutPort(t *testing.T) {
	b := NewSerialBridge("/dev/ttyS0", 9600)
	if err := b.Stop(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestSerialBridge_TransportRoundTrip(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c2.Close()

	b := NewSerialBridge("/dev/ttyUSB1", 115200)
	b.SetPort(c1)

	// Echo: whatever the bridge writes comes back through its transport.
	go func() {
		buf := make([]byte, 64)
		for {
			n, err := c2.Read(buf)
			if err != nil {
				return
			}
			c2.Write(buf[:n])
		}
	}()

	if err := b.Start(); err != nil {
		t.Fatal(err)
	}

	tr := b.Transport()
	if tr == nil {
		t.Fatal("transport is nil after SetPort")
	}
	if _, err := tr.Write([]byte("AT\r")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "AT\r" {
		t.Errorf("expected 'AT\\r', got %q", string(buf[:n]))
	}

	if err := b.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestSerialBridge_StopClosesPort(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c2.Close()

	closed := make(chan struct{})
	go func() {
		buf := make([]byte, 1)
		for {
			if _, err := c2.Read(buf); err != nil {
				close(closed)
				return
			}
		}
	}()

	b := NewSerialBridge("/dev/ttyUSB2", 9600)
	b.SetPort(c1)
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	if err := b.Stop(); err != nil {
		t.Fatal(err)
	}

	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Error("port should have been closed by Stop")
	}
}
